package service

import (
	"api/config"
	"api/pkg/services/api"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
)

// Service ...
type Service struct {
	Config *config.Config
	DB     *sqlx.DB
	api.UnimplementedAPIServiceServer
}

func (s *Service) SaveMediaItemResult(context.Context, *api.MediaItemResultRequest) (*empty.Empty, error) {
	return nil, nil
}

func (s *Service) SaveMediaItemPlace(context.Context, *api.MediaItemPlaceRequest) (*empty.Empty, error) {
	return nil, nil
}
