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
	api.UnimplementedAPIServiceServer
	Config *config.Config
	DB     *gorm.DB
}

func (s *Service) SaveMediaItemResult(ctx context.Context, req *api.MediaItemResultRequest) (*empty.Empty, error) {
	uid, err := uuid.FromString(req.Id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem id")
	}
	creationTime := time.Now()
	if req.CreationTime != nil {
		creationTime, err = time.Parse("2006-01-02 15:04:05 -0700", *req.CreationTime)
		if err != nil {
			log.Printf("error getting mediaitem creation time: %+v", err)
			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time")
		}
	}
	mediaItem := models.MediaItem{
		ID:              uid,
		Status:          models.MediaItemStatus(req.Status),
		Filename:        *req.Filename,
		MimeType:        *req.MimeType,
		SourceURL:       *req.SourceUrl,
		PreviewURL:      *req.PreviewUrl,
		ThumbnailURL:    *req.ThumbnailUrl,
		MediaItemType:   models.MediaItemType(*req.Type),
		Width:           int(*req.Width),
		Height:          int(*req.Height),
		CreationTime:    creationTime,
		CameraMake:      req.CameraMake,
		CameraModel:     req.CameraModel,
		FocalLength:     req.FocalLength,
		ApertureFnumber: req.ApertureFNumber,
		IsoEquivalent:   req.IsoEquivalent,
		ExposureTime:    req.ExposureTime,
		FPS:             req.Fps,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
	}
	result := s.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error updating mediaitem result: %+v", result.Error)
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "error updating mediaitem result: %s", result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) SaveMediaItemPlace(ctx context.Context, req *api.MediaItemPlaceRequest) (*empty.Empty, error) {
	return &emptypb.Empty{}, nil
}
