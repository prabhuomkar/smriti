package service

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/api"
	"api/pkg/storage"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
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

func (s *Service) GetWorkerConfig(context.Context, *empty.Empty) (*api.ConfigResponse, error) {
	type WorkerTask struct {
		Name   string   `json:"name"`
		Source string   `json:"source,omitempty"`
		Params []string `json:"params,omitempty"`
	}
	var workerTasks []WorkerTask
	if s.Config.ML.Places {
		workerTasks = append(workerTasks, WorkerTask{Name: "places", Source: s.Config.ML.PlacesProvider})
	}
	if s.Config.ML.Classification {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   "classification",
			Source: s.Config.ClassificationProvider,
			Params: s.Config.ClassificationParams,
		})
	}
	if s.Config.ML.OCR {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   "ocr",
			Source: s.Config.OCRProvider,
			Params: s.Config.OCRParams,
		})
	}
	if s.Config.ML.Search {
		workerTasks = append(workerTasks, WorkerTask{
			Name:   "search",
			Source: s.Config.SearchProvider,
			Params: s.Config.SearchParams,
		})
	}
	if s.Config.ML.Faces {
		workerTasks = append(workerTasks, WorkerTask{Name: "faces", Params: s.Config.FacesParams})
	}
	if s.Config.ML.Speech {
		workerTasks = append(workerTasks, WorkerTask{Name: "speech", Params: s.Config.SpeechParams})
	}
	configBytes, err := json.Marshal(&workerTasks)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error parsing worker config: %s", err.Error())
	}
	return &api.ConfigResponse{
		Config: configBytes,
	}, nil
}

func (s *Service) SaveMediaItemMetadata(_ context.Context, req *api.MediaItemMetadataRequest) (*empty.Empty, error) { //nolint: cyclop
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem metadata", slog.Any("userId", req.UserId), slog.Any("mediaitem", req.Id), slog.Any("body", req.String()))
	creationTime := time.Now()
	if req.CreationTime != nil {
		creationTime, err = time.Parse("2006-01-02 15:04:05", *req.CreationTime)
		if err != nil {
			slog.Error("error getting mediaitem creation time", slog.Any("error", err))
			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time")
		}
	}

	mediaItem := models.MediaItem{
		UserID: userID, ID: uid, CreationTime: creationTime,
	}
	parseMediaItem(&mediaItem, req)
	mediaItem.SourceURL, err = uploadFile(s.Storage, req.SourcePath, "originals", req.Id)
	if err != nil {
		slog.Error("error uploading original file for mediaitem %s: %+v", req.Id, err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading original file")
	}
	if req.PreviewPath != nil {
		mediaItem.PreviewURL, err = uploadFile(s.Storage, *req.PreviewPath, "previews", req.Id)
		if err != nil {
			slog.Error("error uploading preview file for mediaitem %s: %+v", req.Id, err)
			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading preview file")
		}
	}
	if req.ThumbnailPath != nil {
		mediaItem.ThumbnailURL, err = uploadFile(s.Storage, *req.ThumbnailPath, "thumbnails", req.Id)
		if err != nil {
			slog.Error("error uploading thumbnail file for mediaitem %s: %+v", req.Id, err)
			return &emptypb.Empty{}, status.Error(codes.Internal, "error uploading thumbnail file")
		}
	}
	result := s.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil {
		slog.Error("error updating mediaitem result", slog.Any("error", result.Error))
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", result.Error.Error())
	}
	slog.Info("saved metadata for mediaitem", slog.Any("mediaitem", mediaItem.ID.String()))
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPlace(_ context.Context, req *api.MediaItemPlaceRequest) (*empty.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem place", slog.Any("userId", req.UserId), slog.Any("mediaitem", req.Id), slog.Any("body", req.String()))
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
		slog.Error("error getting or creating place", slog.Any("error", result.Error))
		return &emptypb.Empty{}, status.Errorf(codes.Internal,
			"error getting or creating place: %s", result.Error.Error())
	}
	mediaItem := models.MediaItem{ID: uid}
	err = s.DB.Omit("MediaItems.*").Model(&place).Association("MediaItems").Append(&mediaItem)
	if err != nil {
		slog.Error("error saving mediaitem place", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem place: %s", err.Error())
	}
	slog.Info("saved place for mediaitem", slog.Any("mediaitem", mediaItem.ID.String()))
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemThing(_ context.Context, req *api.MediaItemThingRequest) (*empty.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem thing", slog.Any("userId", req.UserId), slog.Any("mediaitem", req.Id), slog.Any("body", req.String()))
	thing := models.Thing{
		UserID: userID,
		Name:   req.Name,
	}
	result := s.DB.Where(models.Thing{UserID: userID, Name: thing.Name}).
		Attrs(models.Thing{ID: uuid.NewV4()}).
		Assign(models.Thing{CoverMediaItemID: &uid}).
		FirstOrCreate(&thing)
	if result.Error != nil {
		slog.Error("error getting or creating thing", slog.Any("error", result.Error))
		return &emptypb.Empty{}, status.Errorf(codes.Internal,
			"error getting or creating thing: %s", result.Error.Error())
	}
	mediaItem := models.MediaItem{ID: uid}
	err = s.DB.Omit("MediaItems.*").Model(&thing).Association("MediaItems").Append(&mediaItem)
	if err != nil {
		slog.Error("error saving mediaitem thing", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem thing: %s", err.Error())
	}
	slog.Info("saved thing for mediaitem", slog.Any("mediaitem", mediaItem.ID.String()))
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemFaces(_ context.Context, req *api.MediaItemFacesRequest) (*empty.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving mediaitem faces", slog.Any("userId", req.UserId), slog.Any("mediaitem", req.Id))

	mediaItemFaces := make([]models.MediaitemFace, len(req.GetEmbeddings()))
	for idx, reqEmbedding := range req.GetEmbeddings() {
		faceEmbedding := pgvector.NewVector(reqEmbedding.Embedding)
		mediaItemFaces[idx] = models.MediaitemFace{MediaitemID: uid, Embedding: &faceEmbedding}
	}
	result := s.DB.Create(mediaItemFaces)
	if result.Error != nil {
		slog.Error("error saving mediaitem faces", slog.Any("error", result.Error))
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem faces: %s", result.Error.Error())
	}

	slog.Info("saved faces for mediaitem", slog.Any("userId", userID.String()), slog.Any("mediaitem", uid.String()))
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemFinalResult(_ context.Context, req *api.MediaItemFinalResultRequest) (*empty.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		slog.Error("error getting mediaitem user id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	slog.Info("saving final mediaitem result", slog.Any("userId", req.UserId), slog.Any("mediaitem", req.Id))

	mediaItemEmbeddings := make([]models.MediaitemEmbedding, len(req.GetEmbeddings()))
	for idx, reqEmbedding := range req.GetEmbeddings() {
		mediaItemEmbedding := pgvector.NewVector(reqEmbedding.Embedding)
		mediaItemEmbeddings[idx] = models.MediaitemEmbedding{MediaitemID: uid, Embedding: &mediaItemEmbedding}
	}
	result := s.DB.Create(mediaItemEmbeddings)
	if result.Error != nil {
		slog.Error("error saving mediaitem embeddings", slog.Any("error", result.Error))
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem final result: %s", result.Error.Error())
	}

	slog.Info("saved final mediaitem result", slog.Any("userId", userID.String()), slog.Any("mediaitem", uid.String()))
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
