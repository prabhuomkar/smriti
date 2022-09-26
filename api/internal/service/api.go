package service

import (
	"api/config"
	"api/pkg/services/api"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// Service ...
type Service struct {
	api.UnimplementedAPIServiceServer
	Config *config.Config
	DB     *gorm.DB
}

func (s *Service) SaveMediaItemResult(context.Context, *api.MediaItemResultRequest) (*empty.Empty, error) {
	return nil, grpc.ErrServerStopped
}

func (s *Service) SaveMediaItemPlace(context.Context, *api.MediaItemPlaceRequest) (*empty.Empty, error) {
	return nil, grpc.ErrServerStopped
}
