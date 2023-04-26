package service

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/api"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// Service ...
type Service struct {
	api.UnimplementedAPIServer
	Config *config.Config
	DB     *gorm.DB
}

func (s *Service) GetWorkerConfig(context.Context, *empty.Empty) (*api.ConfigResponse, error) {
	type WorkerTask struct {
		Name     string   `json:"name"`
		Source   string   `json:"source,omitempty"`
		Download []string `json:"download,omitempty"`
	}
	var workerTasks []WorkerTask
	if s.Config.ML.Places {
		workerTasks = append(workerTasks, WorkerTask{Name: "places", Source: s.Config.ML.PlacesSource})
	}
	if s.Config.ML.Classification {
		workerTasks = append(workerTasks, WorkerTask{Name: "classification", Download: s.Config.ClassificationDownload})
	}
	if s.Config.ML.Detection {
		workerTasks = append(workerTasks, WorkerTask{Name: "detection", Download: s.Config.DetectionDownload})
	}
	if s.Config.ML.Faces {
		workerTasks = append(workerTasks, WorkerTask{Name: "faces", Download: s.Config.FacesDownload})
	}
	if s.Config.ML.OCR {
		workerTasks = append(workerTasks, WorkerTask{Name: "ocr", Download: s.Config.OCRDownload})
	}
	if s.Config.ML.Speech {
		workerTasks = append(workerTasks, WorkerTask{Name: "speech", Download: s.Config.SpeechDownload})
	}
	configBytes, err := json.Marshal(&workerTasks)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error parsing worker config: %s", err.Error())
	}
	return &api.ConfigResponse{
		Config: configBytes,
	}, nil
}

func (s *Service) SaveMediaItemMetadata(_ context.Context, req *api.MediaItemMetadataRequest) (*empty.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		log.Printf("error getting mediaitem user id: %+v", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	creationTime := time.Now()
	if req.CreationTime != nil {
		creationTime, err = time.Parse("2006-01-02 15:04:05", *req.CreationTime)
		if err != nil {
			log.Printf("error getting mediaitem creation time: %+v", err)
			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time")
		}
	}

	mediaItem := models.MediaItem{
		UserID: userID, ID: uid, CreationTime: creationTime,
	}
	parseMediaItem(&mediaItem, req)
	result := s.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error updating mediaitem result: %+v", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", result.Error.Error())
	}
	log.Printf("saved metadata for mediaitem: %s", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPlace(_ context.Context, req *api.MediaItemPlaceRequest) (*empty.Empty, error) {
	userID, err := uuid.FromString(req.UserId)
	if err != nil {
		log.Printf("error getting mediaitem user id: %+v", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id")
	}
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
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
		Assign(models.Place{CoverMediaItemID: uid}).
		FirstOrCreate(&place)
	if result.Error != nil {
		log.Printf("error getting or creating place: %+v", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal,
			"error getting or creating place: %s", result.Error.Error())
	}
	mediaItem := models.MediaItem{ID: uid}
	err = s.DB.Omit("MediaItems.*").Model(&place).Association("MediaItems").Append(&mediaItem)
	if err != nil {
		log.Printf("error saving mediaitem place: %+v", err)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error saving mediaitem place: %s", err.Error())
	}
	log.Printf("saved place for mediaitem: %s", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func getNameForPlace(place models.Place) string {
	if place.City != nil {
		return *place.City
	}
	return *place.Town
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
	mediaItem.SourceURL = req.SourceUrl
	if req.PreviewUrl != nil {
		mediaItem.PreviewURL = *req.PreviewUrl
	}
	if req.ThumbnailUrl != nil {
		mediaItem.ThumbnailURL = *req.ThumbnailUrl
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
