package service

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/api"
	"context"
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

// nolint: cyclop
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
		UserID: userID, ID: uid, Status: models.MediaItemStatus(req.Status),
		CreationTime: creationTime, CameraMake: req.CameraMake, CameraModel: req.CameraModel,
		FocalLength: req.FocalLength, ApertureFnumber: req.ApertureFNumber, IsoEquivalent: req.IsoEquivalent,
		ExposureTime: req.ExposureTime, FPS: req.Fps, Latitude: req.Latitude, Longitude: req.Longitude,
	}
	if req.MimeType != nil {
		mediaItem.MimeType = *req.MimeType
	}
	if req.SourceUrl != nil {
		mediaItem.SourceURL = *req.SourceUrl
	}
	if req.PreviewUrl != nil {
		mediaItem.PreviewURL = *req.PreviewUrl
	}
	if req.ThumbnailUrl != nil {
		mediaItem.ThumbnailURL = *req.ThumbnailUrl
	}
	if req.Type != nil {
		mediaItem.MediaItemType = models.MediaItemType(*req.Type)
	}
	if req.Width != nil {
		mediaItem.Width = int(*req.Width)
	}
	if req.Height != nil {
		mediaItem.Height = int(*req.Height)
	}
	result := s.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error updating mediaitem result: %+v", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", result.Error.Error())
	}
	log.Printf("saved metadata for mediaitem: %s", mediaItem.ID.String())
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPlace(_ context.Context, req *api.MediaItemPlaceRequest) (*empty.Empty, error) {
	// nolint: golint
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
