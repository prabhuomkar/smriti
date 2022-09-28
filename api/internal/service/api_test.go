package service

import (
	"api/config"
	"api/pkg/services/api"
	"context"
	"errors"
	"log"
	"net"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	filename                    = "filename"
	mimetype                    = "mimetype"
	sourceurl                   = "sourceurl"
	previewurl                  = "previewurl"
	thumbnailurl                = "thumbnailurl"
	mediaitemtype               = "photo"
	badcreationtime             = "bad-creation-time"
	creationtime                = "2022-09-22 11:22:33 +0530"
	width                 int32 = 1080
	height                int32 = 720
	mediaItemReultRequest       = api.MediaItemResultRequest{
		Id:           "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Filename:     &filename,
		MimeType:     &mimetype,
		SourceUrl:    &sourceurl,
		PreviewUrl:   &previewurl,
		ThumbnailUrl: &thumbnailurl,
		Type:         &mediaitemtype,
		Width:        &width,
		Height:       &height,
		CreationTime: &creationtime,
	}
	country               = "country"
	state                 = "state"
	town                  = "town"
	city                  = "city"
	postcode              = "postcode"
	mediaItemPlaceRequest = api.MediaItemPlaceRequest{
		Id:       "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Country:  &country,
		State:    &state,
		Town:     &town,
		City:     &city,
		Postcode: &postcode,
	}
	mediaItemPlaceCountryRequest = api.MediaItemPlaceRequest{
		Id:      "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Country: &country,
	}
	mediaItemPlaceStateRequest = api.MediaItemPlaceRequest{
		Id:    "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		State: &state,
	}
	mediaItemPlaceTownRequest = api.MediaItemPlaceRequest{
		Id:   "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Town: &town,
	}
	sampleTime, _ = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	placeCols     = []string{"id", "name", "postcode", "town", "city", "state",
		"country", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at"}
)

func TestSaveMediaItemResult(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemResultRequest
		MockDB      func(mock sqlmock.Sqlmock)
		ExpectedErr error
	}{
		{
			"save mediaitem result with invalid mediaitem id",
			&api.MediaItemResultRequest{Id: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem result with incorrect creation time",
			&api.MediaItemResultRequest{Id: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", CreationTime: &badcreationtime},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time"),
		},
		{
			"save mediaitem result with success",
			&mediaItemReultRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", "mimetype", "sourceurl", "previewurl",
						"thumbnailurl", "photo", 1080, 720, sqlmock.AnyArg(), sqlmock.AnyArg(),
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem result with error",
			&mediaItemReultRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", "mimetype", "sourceurl", "previewurl",
						"thumbnailurl", "photo", 1080, 720, sqlmock.AnyArg(), sqlmock.AnyArg(),
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error updating mediaitem result: some db error"),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()
			mockGDB, err := gorm.Open(postgres.New(postgres.Config{
				DSN:                  "sqlmock",
				DriverName:           "postgres",
				Conn:                 mockDB,
				PreferSimpleProtocol: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Error),
			})
			assert.NoError(t, err)
			if test.MockDB != nil {
				test.MockDB(mock)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockGDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)))
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIServiceClient(conn)
			_, err = client.SaveMediaItemResult(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestSaveMediaItemPlace(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemPlaceRequest
		MockDB      func(mock sqlmock.Sqlmock)
		ExpectedErr error
	}{
		{
			"save mediaitem place with invalid mediaitem id",
			&api.MediaItemPlaceRequest{Id: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem place with city success",
			&mediaItemPlaceRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem place with country success",
			&mediaItemPlaceCountryRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem place with state success",
			&mediaItemPlaceStateRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem place with town success",
			&mediaItemPlaceTownRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem place with place find or create error",
			&mediaItemPlaceRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(codes.Internal, "error getting or creating place: some db error"),
		},
		{
			"save mediaitem place with error",
			&mediaItemPlaceRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving mediaitem place: some db error"),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()
			mockGDB, err := gorm.Open(postgres.New(postgres.Config{
				DSN:                  "sqlmock",
				DriverName:           "postgres",
				Conn:                 mockDB,
				PreferSimpleProtocol: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Error),
			})
			assert.NoError(t, err)
			if test.MockDB != nil {
				test.MockDB(mock)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockGDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)))
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIServiceClient(conn)
			_, err = client.SaveMediaItemPlace(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func dialer(service *Service) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	api.RegisterAPIServiceServer(server, service)
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func getMockedPlaceRow() *sqlmock.Rows {
	return sqlmock.NewRows(placeCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "postcode", "town", "city",
			"state", "country", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime)
}
