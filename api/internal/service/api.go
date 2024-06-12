package service

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/api"
	"api/pkg/services/worker"
	"api/pkg/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pgvector/pgvector-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// Service ...
type Service struct {
	api.UnimplementedAPIServer
	Config  *config.Config
	DB      *gorm.DB
	Storage storage.Provider
}

func (s *Service) GetWorkerConfig(_ context.Context, _ *emptypb.Empty) (*api.ConfigResponse, error) {
	type WorkerTask struct {
		Name   string `json:"name"`
		Source string `json:"source,omitempty"`
		Params string `json:"params,omitempty"`
	}
	var workerTasks []WorkerTask
	workerTasks = append(workerTasks, WorkerTask{Name: worker.MediaItemComponent_METADATA.String()})
	if len(s.Config.ML.PreviewThumbnailParams) > 0 {
		workerTasks = append(workerTasks, WorkerTask{Name: worker.MediaItemComponent_PREVIEW_THUMBNAIL.String(), Params: s.Config.PreviewThumbnailParams})
	}
	if s.Config.ML.Places {
		workerTasks = append(workerTasks, WorkerTask{Name: worker.MediaItemComponent_PLACES.String(), Source: s.Config.ML.PlacesProvider})
	}
	if s.Config.ML.Classification {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   worker.MediaItemComponent_CLASSIFICATION.String(),
			Source: s.Config.ClassificationProvider,
			Params: s.Config.ClassificationParams,
		})
	}
	if s.Config.ML.OCR {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   worker.MediaItemComponent_OCR.String(),
			Source: s.Config.OCRProvider,
			Params: s.Config.OCRParams,
		})
	}
	if s.Config.ML.Search {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   worker.MediaItemComponent_SEARCH.String(),
			Source: s.Config.SearchProvider,
			Params: s.Config.SearchParams,
		})
	}
	if s.Config.ML.Faces {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   worker.MediaItemComponent_FACES.String(),
			Source: s.Config.FacesProvider,
			Params: s.Config.FacesParams,
		})
	}
	configBytes, err := json.Marshal(&workerTasks)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error parsing worker config: %s", err.Error())
	}
	return &api.ConfigResponse{
		Config: configBytes,
	}, nil
}

func (s *Service) GetUsers(_ context.Context, _ *emptypb.Empty) (*api.GetUsersResponse, error) {
	var userUUIDs []uuid.UUID
	result := s.DB.Model(&models.User{}).Pluck("id", &userUUIDs)
	if result.Error != nil {
		slog.Error("error getting users", "error", result.Error)
		return nil, status.Errorf(codes.Internal, "error getting users: %s", result.Error.Error())
	}

	users := []string{}
	for _, userUUID := range userUUIDs {
		users = append(users, userUUID.String())
	}

	return &api.GetUsersResponse{
		Users: users,
	}, nil
}

func (s *Service) SaveMediaItemMetadata(_ context.Context, req *api.MediaItemMetadataRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem metadata", "userId", req.UserId, "mediaitem", req.Id, "body", req.String())
	creationTime := time.Now()
	if req.CreationTime != nil {
		creationTime, err = time.Parse("2006-01-02 15:04:05", *req.CreationTime)
		if err != nil {
			slog.Error("error getting mediaitem creation time", "error", err)
			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time")
		}
	}

	mediaItem := models.MediaItem{
		UserID: userID, ID: uid, CreationTime: creationTime,
	}
	parseMediaItem(&mediaItem, req)
	result := s.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil {
		slog.Error("error updating mediaitem result", "error", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", result.Error.Error())
	}
	slog.Info("saved metadata for mediaitem", "mediaitem", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPreviewThumbnail(_ context.Context, req *api.MediaItemPreviewThumbnailRequest) (*emptypb.Empty, error) { //nolint: cyclop
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving preview and thumbnail for mediaitem", "userId", req.UserId, "mediaitem", req.Id, "body", req.String())
	mediaItem := models.MediaItem{
		UserID: userID, ID: uid,
	}
	mediaItemUpdates := map[string]interface{}{
		"status": req.Status,
	}
	if req.SourcePath != nil {
		mediaItemUpdates["source_url"], err = uploadFile(s.Storage, *req.SourcePath, "originals", req.Id)
		if err != nil {
			slog.Error("error uploading original file for mediaitem", "id", req.Id, "error", err)
			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading original file")
		}
	}
	if req.Placeholder != nil {
		mediaItemUpdates["placeholder"] = *req.Placeholder
	}
	if req.PreviewPath != nil {
		mediaItemUpdates["preview_url"], err = uploadFile(s.Storage, *req.PreviewPath, "previews", req.Id)
		if err != nil {
			slog.Error("error uploading preview file for mediaitem", "id", req.Id, "error", err)
			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading preview file")
		}
	}
	if req.ThumbnailPath != nil {
		mediaItemUpdates["thumbnail_url"], err = uploadFile(s.Storage, *req.ThumbnailPath, "thumbnails", req.Id)
		if err != nil {
			slog.Error("error uploading thumbnail file for mediaitem", "id", req.Id, "error", err)
			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading thumbnail file")
		}
	}
	result := s.DB.Model(&mediaItem).Updates(mediaItemUpdates)
	if result.Error != nil {
		slog.Error("error updating mediaitem result", "error", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", result.Error.Error())
	}
	slog.Info("saved preview and thumbnail for mediaitem", "mediaitem", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPlace(_ context.Context, req *api.MediaItemPlaceRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem place", "userId", req.UserId, "mediaitem", req.Id, "body", req.String())
	place := models.Place{
		UserID:   userID,
		Postcode: req.Postcode,
		Town:     req.Town,
		City:     req.City,
		State:    req.State,
		Country:  req.Country,
	}
	place.Name = getNameForPlace(place)
	result := s.DB.Where(models.Place{UserID: userID, Name: place.Name, Postcode: place.Postcode}).
		Attrs(models.Place{ID: uuid.NewV4()}).
		Assign(models.Place{CoverMediaItemID: &uid}).
		FirstOrCreate(&place)
	if result.Error != nil {
		slog.Error("error getting or creating place", "error", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal,
			"error getting or creating place: %s", result.Error.Error())
	}
	mediaItem := models.MediaItem{ID: uid}
	err = s.DB.Omit("MediaItems.*").Model(&place).Association("MediaItems").Append(&mediaItem)
	if err != nil {
		slog.Error("error saving mediaitem place", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem place: %s", err.Error())
	}
	slog.Info("saved place for mediaitem", "mediaitem", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemThing(_ context.Context, req *api.MediaItemThingRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem thing", "userId", req.UserId, "mediaitem", req.Id, "body", req.String())
	thing := models.Thing{
		UserID: userID,
		Name:   req.Name,
	}
	result := s.DB.Where(models.Thing{UserID: userID, Name: thing.Name}).
		Attrs(models.Thing{ID: uuid.NewV4()}).
		Assign(models.Thing{CoverMediaItemID: &uid}).
		FirstOrCreate(&thing)
	if result.Error != nil {
		slog.Error("error getting or creating thing", "error", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal,
			"error getting or creating thing: %s", result.Error.Error())
	}
	mediaItem := models.MediaItem{ID: uid}
	err = s.DB.Omit("MediaItems.*").Model(&thing).Association("MediaItems").Append(&mediaItem)
	if err != nil {
		slog.Error("error saving mediaitem thing", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem thing: %s", err.Error())
	}
	slog.Info("saved thing for mediaitem", "mediaitem", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemFaces(_ context.Context, req *api.MediaItemFacesRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem faces", "userId", req.UserId, "mediaitem", req.Id)

	mediaItemFaces := make([]models.MediaitemFace, len(req.GetEmbeddings()))
	faceThumbnails := req.GetThumbnails()
	for idx, reqEmbedding := range req.GetEmbeddings() {
		faceEmbedding := pgvector.NewVector(reqEmbedding.Embedding)
		mediaItemFaces[idx] = models.MediaitemFace{
			MediaitemID: uid, ID: uuid.NewV4(),
			Embedding: &faceEmbedding,
			Thumbnail: faceThumbnails[idx],
		}
	}
	result := s.DB.Create(mediaItemFaces)
	if result.Error != nil {
		slog.Error("error saving mediaitem faces", "error", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem faces: %s", result.Error.Error())
	}

	slog.Info("saved faces for mediaitem", "userId", userID.String(), "mediaitem", uid.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) GetMediaItemFaceEmbeddings(_ context.Context, req *api.MediaItemFaceEmbeddingsRequest) (*api.MediaItemFaceEmbeddingsResponse, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	slog.Info("getting mediaitem face embeddings", "userId", req.UserId)

	mediaItems := []models.MediaItem{}
	result := s.DB.Model(&models.MediaItem{}).
		Where("status=? AND user_id=?", models.Ready, userID).
		Preload("Faces").
		Find(&mediaItems)
	if result.Error != nil {
		slog.Error("error getting mediaitem face embeddings", "error", result.Error)
		return nil, status.Errorf(codes.Internal, "error getting mediaitem face embeddings: %s", result.Error.Error())
	}

	mediaItemFaceEmbeddings := []*api.MediaItemFaceEmbedding{}
	for _, mediaItem := range mediaItems {
		for _, mediaItemFace := range mediaItem.Faces {
			mediaItemFaceEmbedding := &api.MediaItemFaceEmbedding{
				Id:          mediaItemFace.ID.String(),
				MediaItemId: mediaItemFace.MediaitemID.String(),
				Embedding:   &api.MediaItemEmbedding{Embedding: mediaItemFace.Embedding.Slice()},
			}
			if mediaItemFace.PeopleID != nil {
				mediaItemFaceEmbedding.PeopleId = mediaItemFace.PeopleID.String()
			}
			mediaItemFaceEmbeddings = append(mediaItemFaceEmbeddings, mediaItemFaceEmbedding)
		}
	}

	return &api.MediaItemFaceEmbeddingsResponse{
		MediaItemFaceEmbeddings: mediaItemFaceEmbeddings,
	}, nil
}

func (s *Service) SaveMediaItemPeople(_ context.Context, req *api.MediaItemPeopleRequest) (*emptypb.Empty, error) { //nolint:gocognit,funlen,cyclop
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	slog.Info("saving mediaitem people", "userId", req.UserId)

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
				slog.Warn("error getting people id", "faceId", reqPeople, "error", err)
				createdPeopleID, ok := peopleIdxUUIDs[reqPeople]
				if !ok {
					peopleID = uuid.NewV4()
					peopleIdxUUIDs[reqPeople] = peopleID
				} else {
					peopleID = createdPeopleID
				}
			}
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

	for _, peopleID := range peopleIdxUUIDs {
		mediaItems := []*models.MediaItem{}
		for _, mediaItemID := range peopleWithMediaItems[peopleID] {
			mediaItems = append(mediaItems, &models.MediaItem{ID: mediaItemID})
		}
		defaultPeopleVisibility := false
		defaultCoverMediaItemID := mediaItems[0].ID
		defaultCoverMediaItemFaceID := peopleWithFaces[peopleID][0]
		people := models.People{
			IsHidden:             &defaultPeopleVisibility,
			Name:                 "",
			CoverMediaItemID:     &defaultCoverMediaItemID,
			CoverMediaItemFaceID: &defaultCoverMediaItemFaceID,
		}
		result := s.DB.Where(models.People{UserID: userID, ID: peopleID}).
			Attrs(models.People{ID: peopleID}).
			Assign(models.People{CoverMediaItemID: &defaultCoverMediaItemID}).
			FirstOrCreate(&people)
		if result.Error != nil {
			slog.Error("error saving people", "error", result.Error)
			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving people: %s", result.Error.Error())
		}
		err = s.DB.Omit("MediaItems.*").Model(&people).Association("MediaItems").Append(mediaItems)
		if err != nil {
			slog.Error("error saving people mediaitems", "error", err)
			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving people mediaitems: %s", err.Error())
		}
	}

	for peopleID, faces := range peopleWithFaces {
		result := s.DB.Model(&models.MediaitemFace{}).Where("id IN ?", faces).Updates(map[string]interface{}{
			"PeopleID": peopleID,
		})
		if result.Error != nil {
			slog.Error("error saving mediaitem faces people", "error", result.Error, "peopleId", peopleID, "faces", faces)
			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem faces people: %s", result.Error.Error())
		}
	}

	slog.Info("saved people for mediaitem", "userId", userID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemFinalResult(_ context.Context, req *api.MediaItemFinalResultRequest) (*emptypb.Empty, error) { //nolint:cyclop,funlen,gocognit
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving final mediaitem result", "userId", req.UserId, "mediaitem", req.Id)

	if len(req.GetKeywords()) > 0 {
		mediaItem := models.MediaItem{}
		mediaItem.ID = uid
		mediaItem.UserID = userID
		mediaItem.Keywords = &req.Keywords
		result := s.DB.Model(&mediaItem).Updates(mediaItem)
		if result.Error != nil {
			slog.Error("error saving mediaitem keywords", "error", result.Error)
			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem final result: %s", result.Error.Error())
		}
	}

	if len(req.GetEmbeddings()) > 0 {
		mediaItemEmbeddings := make([]models.MediaitemEmbedding, len(req.GetEmbeddings()))
		for idx, reqEmbedding := range req.GetEmbeddings() {
			mediaItemEmbedding := pgvector.NewVector(reqEmbedding.Embedding)
			mediaItemEmbeddings[idx] = models.MediaitemEmbedding{MediaitemID: uid, Embedding: &mediaItemEmbedding}
		}
		result := s.DB.Create(mediaItemEmbeddings)
		if result.Error != nil {
			slog.Error("error saving mediaitem embeddings", "error", result.Error)
			return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem final result: %s", result.Error.Error())
		}
	}

	defer func() {
		err := filepath.WalkDir(s.Config.DiskRoot, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				slog.Error("error iterating over directory for mediaitem", "mediaitem", req.Id, "error", err)
				return err
			}
			if !d.IsDir() && strings.Contains(d.Name(), req.Id) && filepath.Dir(path) == s.Config.DiskRoot {
				// acquire lock to check if not copied
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
						return fmt.Errorf("error removing file for mediaitem %s: %w", req.Id, err)
					}
					break
				}
			}
			return nil
		})
		if err != nil {
			slog.Error("error clearing the files for mediaitem", "mediaitem", req.Id, "error", err)
		} else {
			slog.Debug("cleared the files for mediaitem", "mediaitem", req.Id)
		}
	}()

	slog.Info("saved final mediaitem result", "userId", userID.String(), "mediaitem", uid.String())
	return &emptypb.Empty{}, nil
}

func getNameForPlace(place models.Place) string {
	if place.City != nil {
		return *place.City
	}
	if place.Town != nil {
		return *place.Town
	}
	return *place.State
}

func parseMediaItem(mediaItem *models.MediaItem, req *api.MediaItemMetadataRequest) {
	mediaItem.Status = models.MediaItemStatus(req.Status)
	mediaItem.CameraMake = req.CameraMake
	mediaItem.CameraModel = req.CameraModel
	mediaItem.FocalLength = req.FocalLength
	mediaItem.ApertureFnumber = req.ApertureFNumber
	mediaItem.IsoEquivalent = req.IsoEquivalent
	mediaItem.ExposureTime = req.ExposureTime
	mediaItem.FPS = req.Fps
	mediaItem.Latitude = req.Latitude
	mediaItem.Longitude = req.Longitude
	mediaItem.EXIFData = req.ExifData
	if req.MimeType != nil {
		mediaItem.MimeType = *req.MimeType
	}
	mediaItem.MediaItemType = models.MediaItemType(req.Type)
	mediaItem.MediaItemCategory = models.MediaItemCategory(req.Category)
	if req.Width != nil {
		mediaItem.Width = int(*req.Width)
	}
	if req.Height != nil {
		mediaItem.Height = int(*req.Height)
	}
}

func uploadFile(provider storage.Provider, filePath, fileType, fileID string) (string, error) {
	if len(filePath) == 0 {
		return "", errors.New("error uploading due to invalid file path")
	}
	return provider.Upload(filePath, fileType, fileID)
}
