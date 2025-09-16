package service

import (
	"api/config"
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/services/api"
	"api/pkg/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service ...
type Service struct {
	api.UnimplementedAPIServer

	Config  *config.Config
	DB      database.DBInterface
	Storage storage.Provider

	enabledComponents *[]api.MediaItemComponent
}

var errUploadingInvalidFilePath = errors.New("error uploading due to invalid file path")

const (
	queryGetUsers              = `SELECT id FROM users`
	querySaveMediaItemMetadata = `UPDATE mediaitems SET creation_time=$3, camera_make=$4, camera_model=$5,` +
		` focal_length=$6, aperture_fnumber=$7, iso_equivalent=$8, exposure_time=$9, megapixels=$10, fps=$11,` +
		` latitude=$12, longitude=$13, exif_data=$14, mime_type=$15, mediaitem_type=$16, mediaitem_category=$17,` +
		` width=$18, height=$19, status=$20 WHERE user_id=$1 AND id=$2`
	querySaveMediaItemPreviewThumbnail = `UPDATE mediaitems SET status=$3, source_url=$4, placeholder=$5,` +
		` preview_url=$6, thumbnail_url=$7 WHERE user_id=$1 AND id=$2`
	querySavePlace = `INSERT INTO places (id, user_id, name, postcode, country,` +
		` locality, area, is_hidden, cover_mediaitem_id, created_at, updated_at) VALUES` +
		` ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (user_id, name, postcode)` +
		` DO UPDATE SET cover_mediaitem_id=$9, updated_at=$11`
	querySaveMediaItemPlace = `INSERT INTO place_mediaitems (mediaitem_id, place_id)` +
		` VALUES ($1, $2) ON CONFLICT (mediaitem_id, place_id) DO NOTHING`
	querySaveThing = `INSERT INTO things (id, user_id, name, is_hidden,` +
		` cover_mediaitem_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)` +
		` ON CONFLICT (user_id, name) DO UPDATE SET cover_mediaitem_id=$5, updated_at=$7`
	querySaveMediaItemThing = `INSERT INTO thing_mediaitems (mediaitem_id, thing_id)` +
		` VALUES ($1, $2) ON CONFLICT (mediaitem_id, thing_id) DO NOTHING`
	querySaveMediaItemFaces = `INSERT INTO mediaitem_faces (id, mediaitem_id,` +
		` people_id, embedding, thumbnail) VALUES ($1, $2, $3, $4, $5) ON CONFLICT` +
		` (mediaitem_id, people_id) DO UPDATE SET embedding=$4, thumbnail=$5`
	queryGetMediaItemFaces = `SELECT id, mediaitem_id, people_id,` +
		` embedding FROM mediaitem_faces WHERE mediaitem_id IN` +
		` (SELECT id FROM mediaitems WHERE user_id=$1 AND status=$2)`
	querySaveMediaItemFinalResultKeywords = `UPDATE mediaitems SET keywords=$3 WHERE` +
		` user_id=$1 AND id=$2`
	querySaveMediaItemFinalResultEmbeddings = `INSERT INTO mediaitem_embeddings VALUES($1, $2)`
	querySavePerson                         = `INSERT INTO people (id, user_id, name,` +
		` is_hidden, cover_mediaitem_id, cover_mediaitem_face_id, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT(id) DO UPDATE cover_mediaitem_id=$5,` +
		` cover_mediaitem_face_id=$6, updated_at=$8`
	querySaveMediaItemPeople = `INSERT INTO people_mediaitems (mediaitem_id, people_id)` +
		` VALUES ($1, $2) ON CONFLICT (mediaitem_id, people_id) DO NOTHING`
	querySaveMediaItemFacePeople = `UPDATE mediaitem_faces SET people_id=$2 WHERE id=$1`
	queryUnqueueMediaItem        = `DELETE FROM queue WHERE id=$1`
	queryGetMediaItemProcess     = `UPDATE queue q SET status='PROCESSING' FROM mediaitems m WHERE` +
		` q.id=(SELECT id FROM queue WHERE status='UNSPECIFIED' ORDER BY id FOR UPDATE SKIP LOCKED LIMIT 1)` +
		` AND m.id = q.mediaitem_id RETURNING q.id, q.user_id, q.mediaitem_id, q.components,` +
		` m.mime_type, m.source_url, m.preview_url, m.mediaitem_type, m.mediaitem_category, m.latitude, m.longitude`
)

func Init(cfg *config.Config, dbi database.DBInterface, storage storage.Provider) *Service {
	enabledComponents := []api.MediaItemComponent{api.MediaItemComponent_METADATA}
	if len(cfg.PreviewThumbnailParams) > 0 {
		enabledComponents = append(enabledComponents, api.MediaItemComponent_PREVIEW_THUMBNAIL)
	}
	if cfg.ML.Places {
		enabledComponents = append(enabledComponents, api.MediaItemComponent_PLACES)
	}
	if cfg.Classification {
		enabledComponents = append(enabledComponents, api.MediaItemComponent_CLASSIFICATION)
	}
	if cfg.OCR {
		enabledComponents = append(enabledComponents, api.MediaItemComponent_OCR)
	}
	if cfg.Search {
		enabledComponents = append(enabledComponents, api.MediaItemComponent_SEARCH)
	}
	if cfg.Faces {
		enabledComponents = append(enabledComponents, api.MediaItemComponent_FACES)
	}

	return &Service{Config: cfg, DB: dbi, Storage: storage, enabledComponents: &enabledComponents}
}

func (s *Service) GetWorkerConfig(_ context.Context, _ *emptypb.Empty) (*api.ConfigResponse, error) {
	type Component struct {
		Name   string `json:"name"`
		Source string `json:"source,omitempty"`
		Params string `json:"params,omitempty"`
	}
	var components []Component
	components = append(components, Component{Name: api.MediaItemComponent_METADATA.String()})
	if len(s.Config.PreviewThumbnailParams) > 0 {
		components = append(components, Component{
			Name: api.MediaItemComponent_PREVIEW_THUMBNAIL.String(), Params: s.Config.PreviewThumbnailParams,
		})
	}
	if s.Config.ML.Places {
		components = append(components, Component{
			Name: api.MediaItemComponent_PLACES.String(), Source: s.Config.PlacesProvider,
		})
	}
	if s.Config.Classification {
		components = append(components, Component{
			Name: api.MediaItemComponent_CLASSIFICATION.String(), Source: s.Config.ClassificationProvider, Params: s.Config.ClassificationParams,
		})
	}
	if s.Config.OCR {
		components = append(components, Component{
			Name: api.MediaItemComponent_OCR.String(), Source: s.Config.OCRProvider, Params: s.Config.OCRParams,
		})
	}
	if s.Config.Search {
		components = append(components, Component{
			Name: api.MediaItemComponent_SEARCH.String(), Source: s.Config.SearchProvider, Params: s.Config.SearchParams,
		})
	}
	if s.Config.Faces {
		components = append(components, Component{
			Name: api.MediaItemComponent_FACES.String(), Source: s.Config.FacesProvider, Params: s.Config.FacesParams,
		})
	}
	configBytes, err := json.Marshal(&components)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error parsing worker config: %s", err.Error())
	}

	return &api.ConfigResponse{Config: configBytes}, nil
}

func (s *Service) GetMediaItemProcess(ctx context.Context, _ *emptypb.Empty) (*api.MediaItemProcessResponse, error) { //nolint:cyclop
	var (
		queueID           uuid.UUID
		userID            uuid.UUID
		mediaItemID       uuid.UUID
		components        string
		mimeType          *string
		sourceURL         string
		previewURL        *string
		mediaItemType     *string
		mediaItemCategory *string
		latitude          *string
		longitude         *string
	)
	err := s.DB.QueryRow(ctx, queryGetMediaItemProcess).
		Scan(&queueID, &userID, &mediaItemID, &components, &mimeType, &sourceURL, &previewURL, &mediaItemType,
			&mediaItemCategory, &latitude, &longitude)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.MediaItemProcessResponse{}, nil
		}

		slog.Error("error getting mediaitem to process", "error", err)

		return nil, status.Errorf(codes.Internal, "error getting mediaitem to process: %s", err.Error())
	}

	filteredComponents := []api.MediaItemComponent{}
	queueComponents := strings.Split(components, ",")
	for _, queueComponent := range queueComponents {
		component := api.MediaItemComponent(api.MediaItemComponent_value[queueComponent])
		if slices.Contains(*s.enabledComponents, component) {
			filteredComponents = append(filteredComponents, component)
		}
	}

	payload := map[string]string{"source_url": sourceURL}
	if mimeType != nil {
		payload["mime_type"] = *mimeType
	}
	if previewURL != nil {
		payload["preview_url"] = *previewURL
	}
	if mediaItemType != nil {
		payload["type"] = *mediaItemType
	}
	if mediaItemCategory != nil {
		payload["category"] = *mediaItemCategory
	}
	if latitude != nil {
		payload["latitude"] = *latitude
	}
	if longitude != nil {
		payload["longitude"] = *longitude
	}

	return &api.MediaItemProcessResponse{
		Id: queueID.String(), UserId: userID.String(),
		MediaItemId: mediaItemID.String(), Components: filteredComponents,
		Payload: payload,
	}, nil
}

func (s *Service) GetUsers(ctx context.Context, _ *emptypb.Empty) (*api.UsersResponse, error) {
	var userUUIDs []uuid.UUID
	rows, err := s.DB.Query(ctx, queryGetUsers)
	if err != nil {
		slog.Error("error getting users", "error", err)

		return nil, status.Errorf(codes.Internal, "error getting users: %s", err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var userUUID uuid.UUID
		if err := rows.Scan(&userUUID); err != nil {
			slog.Error("error scanning user", "error", err)

			return nil, status.Errorf(codes.Internal, "error scanning user: %s", err.Error())
		}
		userUUIDs = append(userUUIDs, userUUID)
	}

	users := []string{}
	for _, userUUID := range userUUIDs {
		users = append(users, userUUID.String())
	}

	return &api.UsersResponse{Users: users}, nil
}

func (s *Service) SaveMediaItemMetadata(ctx context.Context, req *api.MediaItemMetadataRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	mediaItemID, err := uuid.FromString(req.MediaItemId)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem metadata", "user", req.UserId, "mediaitem", req.MediaItemId, "body", req.String())
	creationTime := time.Now()
	if req.CreationTime != nil {
		creationTime, err = time.Parse("2006-01-02 15:04:05", *req.CreationTime)
		if err != nil {
			slog.Error("error getting mediaitem creation time", "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time")
		}
	}

	mediaItem := models.MediaItem{UserID: userID, ID: mediaItemID, CreationTime: creationTime}
	parseMediaItem(&mediaItem, req)
	_, err = s.DB.Exec(ctx, querySaveMediaItemMetadata, userID, mediaItemID, mediaItem.CreationTime, mediaItem.CameraMake,
		mediaItem.CameraModel, mediaItem.FocalLength, mediaItem.ApertureFnumber, mediaItem.IsoEquivalent,
		mediaItem.ExposureTime, mediaItem.Megapixels, mediaItem.FPS, mediaItem.Latitude, mediaItem.Longitude,
		mediaItem.EXIFData, mediaItem.MimeType, mediaItem.MediaItemType, mediaItem.MediaItemCategory,
		mediaItem.Width, mediaItem.Height, mediaItem.Status)
	if err != nil {
		slog.Error("error saving mediaitem metadata", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", err.Error())
	}
	slog.Info("saved metadata for mediaitem", "user", req.UserId, "mediaitem", mediaItem.ID.String())

	return &emptypb.Empty{}, nil
}

//nolint:cyclop
func (s *Service) SaveMediaItemPreviewThumbnail(ctx context.Context, req *api.MediaItemPreviewThumbnailRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	mediaItemID, err := uuid.FromString(req.MediaItemId)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving preview and thumbnail for mediaitem", "user", req.UserId, "mediaitem", req.MediaItemId, "body", req.String())
	mediaItemUpdates := map[string]interface{}{"status": req.Status}
	if req.SourcePath != nil {
		mediaItemUpdates["source_url"], err = uploadFile(s.Storage, *req.SourcePath, "originals", req.MediaItemId)
		if err != nil {
			slog.Error("error uploading original file for mediaitem", "id", req.MediaItemId, "error", err)

			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading original file")
		}
	}
	if req.Placeholder != nil {
		mediaItemUpdates["placeholder"] = *req.Placeholder
	}
	if req.PreviewPath != nil {
		mediaItemUpdates["preview_url"], err = uploadFile(s.Storage, *req.PreviewPath, "previews", req.MediaItemId)
		if err != nil {
			slog.Error("error uploading preview file for mediaitem", "id", req.MediaItemId, "error", err)

			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading preview file")
		}
	}
	if req.ThumbnailPath != nil {
		mediaItemUpdates["thumbnail_url"], err = uploadFile(s.Storage, *req.ThumbnailPath, "thumbnails", req.MediaItemId)
		if err != nil {
			slog.Error("error uploading thumbnail file for mediaitem", "id", req.MediaItemId, "error", err)

			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading thumbnail file")
		}
	}
	_, err = s.DB.Exec(ctx, querySaveMediaItemPreviewThumbnail, userID, mediaItemID, req.Status,
		mediaItemUpdates["source_url"], mediaItemUpdates["placeholder"], mediaItemUpdates["preview_url"],
		mediaItemUpdates["thumbnail_url"])
	if err != nil {
		slog.Error("error saving mediaitem preview and thumbnail", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", err.Error())
	}
	slog.Info("saved preview and thumbnail for mediaitem", "user", req.UserId, "mediaitem", req.MediaItemId)

	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPlace(ctx context.Context, req *api.MediaItemPlaceRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	mediaItemID, err := uuid.FromString(req.MediaItemId)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem place", "user", req.UserId, "mediaitem", req.MediaItemId, "body", req.String())
	place := models.Place{
		ID:     uuid.NewV4(),
		UserID: userID, Postcode: req.Postcode, Country: req.Country, Locality: req.Locality, Area: req.Area,
	}
	place.Name = getNameForPlace(place)
	place.CreatedAt = time.Now()
	place.UpdatedAt = place.CreatedAt
	ptx, err := s.DB.Begin(ctx)
	if err != nil {
		slog.Error("error starting transaction for saving mediaitem place", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error starting transaction for saving mediaitem place: %s", err.Error())
	}
	defer func() {
		if err != nil {
			_ = ptx.Rollback(ctx)
		}
	}()
	_, err = ptx.Exec(ctx, querySavePlace, place.ID, userID, place.Name, place.Postcode,
		place.Country, place.Locality, place.Area, false, mediaItemID, place.CreatedAt, place.UpdatedAt)
	if err != nil {
		slog.Error("error saving place", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving place: %s", err.Error())
	}
	_, err = ptx.Exec(ctx, querySaveMediaItemPlace, mediaItemID, place.ID)
	if err != nil {
		slog.Error("error saving mediaitem place", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem place: %s", err.Error())
	}
	if err = ptx.Commit(ctx); err != nil {
		slog.Error("error committing transaction for saving mediaitem place", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error committing transaction for saving mediaitem place: %s", err.Error())
	}
	slog.Info("saved place for mediaitem", "user", req.UserId, "mediaitem", req.MediaItemId)

	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemThing(ctx context.Context, req *api.MediaItemThingRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	mediaItemID, err := uuid.FromString(req.MediaItemId)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem thing", "user", req.UserId, "mediaitem", req.MediaItemId, "body", req.String())
	thing := models.Thing{ID: uuid.NewV4(), UserID: userID, Name: req.Name}
	thing.CreatedAt = time.Now()
	thing.UpdatedAt = thing.CreatedAt
	ttx, err := s.DB.Begin(ctx)
	if err != nil {
		slog.Error("error starting transaction for saving mediaitem thing", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error starting transaction for saving mediaitem thing: %s", err.Error())
	}
	defer func() {
		if err != nil {
			_ = ttx.Rollback(ctx)
		}
	}()
	_, err = ttx.Exec(ctx, querySaveThing, thing.ID, thing.UserID, thing.Name, false, mediaItemID, thing.CreatedAt, thing.UpdatedAt)
	if err != nil {
		slog.Error("error saving thing", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving thing: %s", err.Error())
	}
	_, err = ttx.Exec(ctx, querySaveMediaItemThing, mediaItemID, thing.ID)
	if err != nil {
		slog.Error("error saving mediaitem thing", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem thing: %s", err.Error())
	}
	if err = ttx.Commit(ctx); err != nil {
		slog.Error("error committing transaction for saving mediaitem thing", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error committing transaction for saving mediaitem thing: %s", err.Error())
	}
	slog.Info("saved thing for mediaitem", "user", req.UserId, "mediaitem", req.MediaItemId)

	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemFaces(ctx context.Context, req *api.MediaItemFacesRequest) (*emptypb.Empty, error) {
	_, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	mediaItemID, err := uuid.FromString(req.MediaItemId)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem faces", "user", req.UserId, "mediaitem", req.MediaItemId)

	mediaItemFaces := make([]models.MediaitemFace, len(req.GetEmbeddings()))
	faceThumbnails := req.GetThumbnails()
	for idx, reqEmbedding := range req.GetEmbeddings() {
		faceEmbedding := pgvector.NewVector(reqEmbedding.Embedding)
		mediaItemFaces[idx] = models.MediaitemFace{
			MediaitemID: mediaItemID, ID: uuid.NewV4(), Embedding: &faceEmbedding, Thumbnail: faceThumbnails[idx],
		}
	}
	for idx, mediaItemFace := range mediaItemFaces {
		_, err = s.DB.Exec(ctx, querySaveMediaItemFaces, mediaItemFace.ID, mediaItemFace.MediaitemID, nil, mediaItemFace.Embedding, mediaItemFace.Thumbnail)
		if err != nil {
			slog.Error("error saving mediaitem faces", "idx", idx, "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem faces: %s", err.Error())
		}
	}

	slog.Info("saved faces for mediaitem", "user", req.UserId, "mediaitem", req.MediaItemId)

	return &emptypb.Empty{}, nil
}

func (s *Service) GetMediaItemFaceEmbeddings(ctx context.Context, req *api.MediaItemFaceEmbeddingsRequest) (*api.MediaItemFaceEmbeddingsResponse, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	slog.Info("getting mediaitem face embeddings", "user", req.UserId)

	mediaItemFaces := []models.MediaitemFace{}
	rows, err := s.DB.Query(ctx, queryGetMediaItemFaces, userID, api.MediaItemStatus_READY.String())
	if err != nil {
		slog.Error("error getting mediaitem face embeddings", "error", err)

		return nil, status.Errorf(codes.Internal, "error getting mediaitem face embeddings: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var mediaItemFace models.MediaitemFace
		if err := rows.Scan(&mediaItemFace.ID, &mediaItemFace.MediaitemID, &mediaItemFace.PeopleID, &mediaItemFace.Embedding); err != nil {
			slog.Error("error scanning mediaitem face embedding", "error", err)

			return nil, status.Errorf(codes.Internal, "error scanning mediaitem face embedding: %s", err.Error())
		}
		mediaItemFaces = append(mediaItemFaces, mediaItemFace)
	}

	mediaItemFaceEmbeddings := []*api.MediaItemFaceEmbedding{}
	for _, mediaItemFace := range mediaItemFaces {
		mediaItemFaceEmbedding := &api.MediaItemFaceEmbedding{
			Id: mediaItemFace.ID.String(), MediaItemId: mediaItemFace.MediaitemID.String(),
		}
		if mediaItemFace.Embedding != nil {
			mediaItemFaceEmbedding.Embedding = &api.MediaItemEmbedding{Embedding: mediaItemFace.Embedding.Slice()}
		}
		if mediaItemFace.PeopleID != nil {
			mediaItemFaceEmbedding.PeopleId = mediaItemFace.PeopleID.String()
		}
		mediaItemFaceEmbeddings = append(mediaItemFaceEmbeddings, mediaItemFaceEmbedding)
	}

	return &api.MediaItemFaceEmbeddingsResponse{MediaItemFaceEmbeddings: mediaItemFaceEmbeddings}, nil
}

//nolint:gocognit,cyclop
func (s *Service) SaveMediaItemPeople(ctx context.Context, req *api.MediaItemPeopleRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	slog.Info("saving mediaitem people", "user", req.UserId)

	peopleWithFaces := map[uuid.UUID][]uuid.UUID{}
	peopleWithMediaItems := map[uuid.UUID][]uuid.UUID{}
	peopleIdxUUIDs := map[string]uuid.UUID{}
	for reqMediaItem, reqFacePeople := range req.GetMediaItemFacePeople() {
		mediaItemID, err := uuid.FromString(reqMediaItem)
		if err != nil {
			slog.Error("error getting mediaitem id", "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
		}
		for reqFace, reqPeople := range reqFacePeople.GetFacePeople() {
			faceID, err := uuid.FromString(reqFace)
			if err != nil {
				slog.Error("error getting face id", "error", err)

				return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid face id")
			}
			peopleID, err := uuid.FromString(reqPeople)
			if err != nil {
				slog.Warn("error getting people id", "faceId", faceID, "peopleId", reqPeople, "error", err)
				createdPeopleID, ok := peopleIdxUUIDs[reqPeople]
				if !ok {
					peopleID = uuid.NewV4()
				} else {
					peopleID = createdPeopleID
				}
			}
			peopleIdxUUIDs[reqPeople] = peopleID
			_, idxOk := peopleWithFaces[peopleID]
			if idxOk {
				peopleWithFaces[peopleID] = append(peopleWithFaces[peopleID], faceID)
			} else {
				peopleWithFaces[peopleID] = []uuid.UUID{faceID}
			}
			_, idxOk = peopleWithMediaItems[peopleID]
			if idxOk {
				peopleWithMediaItems[peopleID] = append(peopleWithMediaItems[peopleID], mediaItemID)
			} else {
				peopleWithMediaItems[peopleID] = []uuid.UUID{mediaItemID}
			}
		}
	}

	for idx, peopleID := range peopleIdxUUIDs {
		defaultPeopleVisibility := false
		defaultCoverMediaItemID := peopleWithMediaItems[peopleID][0]
		defaultCoverMediaItemFaceID := peopleWithFaces[peopleID][0]
		people := models.People{
			IsHidden: &defaultPeopleVisibility, Name: "", CoverMediaItemID: &defaultCoverMediaItemID, CoverMediaItemFaceID: &defaultCoverMediaItemFaceID,
		}
		people.CreatedAt = time.Now()
		people.UpdatedAt = people.CreatedAt
		ptx, err := s.DB.Begin(ctx)
		if err != nil {
			slog.Error("error starting transaction for saving mediaitem people", "idx", idx, "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error starting transaction for saving mediaitem people: %s", err.Error())
		}
		defer func() {
			if err != nil {
				_ = ptx.Rollback(ctx)
			}
		}()
		_, err = ptx.Exec(ctx, querySavePerson, peopleID, userID, people.Name, people.IsHidden, people.CoverMediaItemID, people.CoverMediaItemFaceID, people.CreatedAt, people.UpdatedAt)
		if err != nil {
			slog.Error("error saving people", "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving people: %s", err.Error())
		}
		for midx, mediaItemID := range peopleWithMediaItems[peopleID] {
			_, err = ptx.Exec(ctx, querySaveMediaItemPeople, mediaItemID, peopleID)
			if err != nil {
				slog.Error("error saving people mediaitems", "idx", idx, "mediaItemIdx", midx, "error", err)

				return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving people mediaitems: %s", err.Error())
			}
		}
		for fidx, faceID := range peopleWithFaces[peopleID] {
			_, err = ptx.Exec(ctx, querySaveMediaItemFacePeople, faceID, peopleID)
			if err != nil {
				slog.Error("error saving mediaitem faces", "idx", idx, "faceIdx", fidx, "error", err)

				return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem faces: %s", err.Error())
			}
		}
		if err = ptx.Commit(ctx); err != nil {
			slog.Error("error committing transaction for saving mediaitem people", "idx", idx, "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error committing transaction for saving mediaitem people: %s", err.Error())
		}
	}

	slog.Info("saved people for mediaitem", "user", userID.String())

	return &emptypb.Empty{}, nil
}

//nolint:gocognit,cyclop
func (s *Service) SaveMediaItemFinalResult(ctx context.Context, req *api.MediaItemFinalResultRequest) (*emptypb.Empty, error) {
	queueID, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting queue id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid queue id")
	}
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	mediaItemID, err := uuid.FromString(req.MediaItemId)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving final mediaitem result", "user", req.UserId, "mediaitem", req.MediaItemId)

	if len(req.GetKeywords()) > 0 {
		_, err = s.DB.Exec(ctx, querySaveMediaItemFinalResultKeywords, userID, mediaItemID, req.GetKeywords())
		if err != nil {
			slog.Error("error saving mediaitem keywords", "error", err)

			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem final result keywords: %s", err.Error())
		}
	}

	if len(req.GetEmbeddings()) > 0 {
		for idx, reqEmbedding := range req.GetEmbeddings() {
			mediaItemEmbedding := pgvector.NewVector(reqEmbedding.Embedding)
			_, err = s.DB.Exec(ctx, querySaveMediaItemFinalResultEmbeddings, mediaItemID, mediaItemEmbedding)
			if err != nil {
				slog.Error("error saving mediaitem embedding", "idx", idx, "error", err)

				return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem final result embedding: %s", err.Error())
			}
		}
	}

	_, err = s.DB.Exec(ctx, queryUnqueueMediaItem, queueID)
	if err != nil {
		slog.Error("error unqueuing mediaitem from processing", "error", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error unqueuing mediaitem from processing: %s", err.Error())
	}

	defer func() {
		err := filepath.WalkDir(s.Config.DiskRoot, func(path string, dir os.DirEntry, err error) error {
			if err != nil {
				slog.Error("error iterating over directory for mediaitem", "mediaitem", req.MediaItemId, "error", err)

				return err
			}
			if !dir.IsDir() && strings.Contains(dir.Name(), req.MediaItemId) &&
				filepath.Dir(path) == s.Config.DiskRoot { // acquire lock to check if not copied
				for {
					slog.Debug("deleting file", "path", path)
					file, err := os.Open(path)
					if err != nil {
						continue
					}
					if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
						continue
					}
					if err = os.Remove(path); err != nil {
						return fmt.Errorf("error removing file for mediaitem %s: %w", req.MediaItemId, err)
					}

					break
				}
			}

			return nil
		})
		if err != nil {
			slog.Error("error clearing the files for mediaitem", "mediaitem", req.MediaItemId, "error", err)
		} else {
			slog.Debug("cleared the files for mediaitem", "mediaitem", req.MediaItemId)
		}
	}()

	slog.Info("saved final mediaitem result", "user", req.UserId, "mediaitem", req.MediaItemId)

	return &emptypb.Empty{}, nil
}

func getNameForPlace(place models.Place) string {
	if place.Locality != nil {
		return *place.Locality
	}
	if place.Area != nil {
		return *place.Area
	}

	return *place.Country
}

func parseMediaItem(mediaItem *models.MediaItem, req *api.MediaItemMetadataRequest) {
	mediaItem.Status = req.Status.String()
	mediaItem.CameraMake = req.CameraMake
	mediaItem.CameraModel = req.CameraModel
	mediaItem.FocalLength = req.FocalLength
	mediaItem.ApertureFnumber = req.ApertureFNumber
	mediaItem.IsoEquivalent = req.IsoEquivalent
	mediaItem.ExposureTime = req.ExposureTime
	mediaItem.Megapixels = req.Megapixels
	mediaItem.FPS = req.Fps
	mediaItem.Latitude = req.Latitude
	mediaItem.Longitude = req.Longitude
	mediaItem.EXIFData = req.ExifData
	if req.MimeType != nil {
		mediaItem.MimeType = *req.MimeType
	}
	mediaItem.MediaItemType = req.Type.String()
	mediaItem.MediaItemCategory = req.Category.String()
	if req.Width != nil {
		mediaItem.Width = int(*req.Width)
	}
	if req.Height != nil {
		mediaItem.Height = int(*req.Height)
	}
}

func uploadFile(provider storage.Provider, filePath, fileType, fileID string) (string, error) {
	if len(filePath) == 0 {
		return "", errUploadingInvalidFilePath
	}

	return provider.Upload(filePath, fileType, fileID)
}
