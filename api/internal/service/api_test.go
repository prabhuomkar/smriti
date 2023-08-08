package service

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"testing"
	"time"

	"api/config"
	"api/pkg/services/api"
	"api/pkg/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mimetype                    = "mimetype"
	mediaitemType               = "photo"
	mediaitemCategory           = "default"
	badcreationtime             = "bad-creation-time"
	creationtime                = "2022-09-22 11:22:33"
	width                 int32 = 1080
	height                int32 = 720
	existingPlaceKeywords       = "placecity placepostcode"
	mediaItemReultRequest       = api.MediaItemMetadataRequest{
		UserId:       "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:           "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		MimeType:     &mimetype,
		Type:         mediaitemType,
		Category:     mediaitemCategory,
		Width:        &width,
		Height:       &height,
		CreationTime: &creationtime,
	}
	country               = "country"
	state                 = "state"
	town                  = "town"
	city                  = "city"
	postcode              = "postcode"
	embedding             = pgvector.NewVector([]float32{0.0, 0.42, 0.111})
	mediaItemEmbedding    = api.MediaItemEmbedding{Embedding: []float32{0.0, 0.42, 0.111}}
	mediaItemPlaceRequest = api.MediaItemPlaceRequest{
		UserId:   "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:       "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Country:  &country,
		State:    &state,
		Town:     &town,
		City:     &city,
		Postcode: &postcode,
	}
	mediaItemPlaceTownRequest = api.MediaItemPlaceRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Town:   &town,
	}
	mediaItemPlaceStateRequest = api.MediaItemPlaceRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		State:  &state,
	}
	mediaItemThingRequest = api.MediaItemThingRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Name:   "Pizza",
	}
	mediaItemFacesRequest = api.MediaItemFacesRequest{
		UserId:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:         "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Embeddings: []*api.MediaItemEmbedding{&mediaItemEmbedding},
	}
	mediaItemPeopleRequest = api.MediaItemPeopleRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179": {
				FacePeople: map[string]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179": "4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			},
			"4d05b5f6-17c2-475e-87fe-3fc8b9567180": {
				FacePeople: map[string]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567180": "4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			},
			"4d05b5f6-17c2-475e-87fe-3fc8b9567181": {
				FacePeople: map[string]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567181": "1"},
			},
			"4d05b5f6-17c2-475e-87fe-3fc8b9567182": {
				FacePeople: map[string]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567182": "1"},
			},
		},
	}
	mediaItemFaceEmbeddingsRequest = api.MediaItemFaceEmbeddingsRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
	}
	mediaItemFinalResultRequest = api.MediaItemFinalResultRequest{
		UserId:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:         "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Keywords:   "some keywords",
		Embeddings: []*api.MediaItemEmbedding{&mediaItemEmbedding},
	}
	sampleTime, _ = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	placeCols     = []string{
		"id", "name", "postcode", "town", "city", "state",
		"country", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at",
	}
	thingCols = []string{
		"id", "name", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at",
	}
	peopleCols = []string{
		"id", "name", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at",
	}
	mediaitemCols = []string{
		"id", "user_id", "filename", "description", "mime_type", "keywords", "source_url", "preview_url",
		"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "mediaitem_category",
		"width", "height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "latitude", "longitude", "fps", "created_at", "updated_at",
	}
	mediaitemFaceCols = []string{
		"id", "mediaitem_id", "people_id", "embedding",
	}
)

func TestGetWorkerConfig(t *testing.T) {
	tests := []struct {
		Name           string
		Config         *config.Config
		ExpectedConfig []byte
		ExpectedErr    error
	}{
		{
			"get worker config with success",
			&config.Config{ML: config.ML{Places: true, PlacesProvider: "openstreetmap"}},
			[]byte(`[{"name":"places","source":"openstreetmap"}]`),
			nil,
		},
		{
			"get worker config with success with all config",
			&config.Config{ML: config.ML{
				Places: true, PlacesProvider: "openstreetmap",
				Classification: true, ClassificationProvider: "pytorch", ClassificationParams: []string{"model-file-name.pt"},
				OCR: true, OCRProvider: "paddlepaddle", OCRParams: []string{"ocr-v1-model.pt"},
				Search: true, SearchProvider: "pytorch", SearchParams: []string{"search-model.pt"},
				Faces: true, FacesParams: []string{"http://faces/model/link"},
				MetadataParams: []string{"512"},
			}},
			[]byte(`[{"name":"metadata","params":["512"]},` +
				`{"name":"places","source":"openstreetmap"},{"name":"classification","source":"pytorch",` +
				`"params":["model-file-name.pt"]},{"name":"ocr","source":"paddlepaddle","params":["ocr-v1-model.pt"]}` +
				`,{"name":"search","source":"pytorch","params":["search-model.pt"]}` +
				`,{"name":"faces","params":["http://faces/model/link"]}]`),
			nil,
		},
		{
			"get worker config with no error",
			&config.Config{},
			[]byte(`null`),
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// service
			service := &Service{
				Config: test.Config,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)))
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIClient(conn)
			res, err := client.GetWorkerConfig(ctx, &emptypb.Empty{})
			// assert
			assert.Equal(t, test.ExpectedConfig, res.Config)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		Name           string
		MockDB         func(mock sqlmock.Sqlmock)
		ExpectedResult *api.GetUsersResponse
		ExpectedErr    error
	}{
		{
			"get users with success",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users"`)).
					WillReturnRows(getMockedUserIDRows())
			},
			&api.GetUsersResponse{
				Users: []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567180"},
			},
			nil,
		},
		{
			"get users with error",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			status.Error(codes.Internal, "error getting users: some db error"),
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
			client := api.NewAPIClient(conn)
			res, err := client.GetUsers(ctx, &emptypb.Empty{})
			// assert
			assert.Equal(t, test.ExpectedErr, err)
			if test.ExpectedResult != nil {
				assert.Equal(t, test.ExpectedResult.Users, res.Users)
			} else {
				assert.Equal(t, test.ExpectedResult, res)
			}
		})
	}
}
func TestSaveMediaItemMetadata(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemMetadataRequest
		MockDB      func(mock sqlmock.Sqlmock)
		MockParams  func(string) (string, string, string, func(), error)
		ExpectedErr error
	}{
		{
			"save mediaitem result with invalid mediaitem user id",
			&api.MediaItemMetadataRequest{UserId: "bad-mediaitem-user-id"},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem result with invalid mediaitem id",
			&api.MediaItemMetadataRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", Id: "bad-mediaitem-id"},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem result with incorrect creation time",
			&api.MediaItemMetadataRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179", CreationTime: &badcreationtime,
			},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time"),
		},
		{
			"save mediaitem result with error uploading original file",
			&mediaItemReultRequest,
			nil,
			func(tmpRoot string) (string, string, string, func(), error) {
				return "", "", "", func() {}, nil
			},
			status.Errorf(codes.Internal, "error uploading original file"),
		},
		{
			"save mediaitem result with error uploading preview file",
			&mediaItemReultRequest,
			nil,
			func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0777)
				return originalFile.Name(), "", "", func() {
					defer os.Remove(tmpRoot + "/originals/")
					defer os.Remove(originalFile.Name())
				}, nil
			},
			status.Errorf(codes.Internal, "error uploading preview file"),
		},
		{
			"save mediaitem result with error uploading thumbnail file",
			&mediaItemReultRequest,
			nil,
			func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0777)
				previewFile, err := os.CreateTemp(tmpRoot, "preview")
				if err != nil {
					return "", "", "", nil, err
				}
				return originalFile.Name(), previewFile.Name(), "", func() {
					defer os.Remove(tmpRoot + "/originals/")
					defer os.Remove(originalFile.Name())
					defer os.Remove(tmpRoot + "/previews/")
					defer os.Remove(previewFile.Name())
				}, nil
			},
			status.Errorf(codes.Internal, "error uploading thumbnail file"),
		},
		{
			"save mediaitem result with success",
			&mediaItemReultRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
						"mimetype", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "photo", "default",
						1080, 720, sqlmock.AnyArg(), sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0777)
				previewFile, err := os.CreateTemp(tmpRoot, "preview")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/thumbnails/", 0777)
				thumbnailFile, err := os.CreateTemp(tmpRoot, "thumbnail")
				if err != nil {
					return "", "", "", nil, err
				}
				return originalFile.Name(), previewFile.Name(), thumbnailFile.Name(), func() {
					defer os.Remove(tmpRoot + "/originals/")
					defer os.Remove(originalFile.Name())
					defer os.Remove(tmpRoot + "/previews/")
					defer os.Remove(previewFile.Name())
					defer os.Remove(tmpRoot + "/thumbnails/")
					defer os.Remove(thumbnailFile.Name())
				}, nil
			},
			nil,
		},
		{
			"save mediaitem result with error",
			&mediaItemReultRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
						"mimetype", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "photo", "default",
						1080, 720, sqlmock.AnyArg(), sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0777)
				previewFile, err := os.CreateTemp(tmpRoot, "preview")
				assert.NoError(t, err)
				os.Mkdir(tmpRoot+"/thumbnails/", 0777)
				thumbnailFile, err := os.CreateTemp(tmpRoot, "thumbnail")
				assert.NoError(t, err)
				return originalFile.Name(), previewFile.Name(), thumbnailFile.Name(), func() {
					defer os.Remove(tmpRoot + "/originals/")
					defer os.Remove(originalFile.Name())
					defer os.Remove(tmpRoot + "/previews/")
					defer os.Remove(previewFile.Name())
					defer os.Remove(tmpRoot + "/thumbnails/")
					defer os.Remove(thumbnailFile.Name())
				}, nil
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
			tmpRoot := os.TempDir()
			service := &Service{
				Config:  &config.Config{},
				DB:      mockGDB,
				Storage: &storage.Disk{Root: tmpRoot},
			}
			// mock tmp params
			if test.MockParams != nil {
				originalPath, previewPath, thumbnailPath, clear, err := test.MockParams(tmpRoot)
				assert.NoError(t, err)
				test.Request.SourcePath = originalPath
				test.Request.PreviewPath = &previewPath
				test.Request.ThumbnailPath = &thumbnailPath
				defer clear()
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)))
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemMetadata(ctx, test.Request)
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
			"save mediaitem place with invalid mediaitem user id",
			&api.MediaItemPlaceRequest{UserId: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem place with invalid mediaitem id",
			&api.MediaItemPlaceRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", Id: "bad-mediaitem-id"},
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
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow(nil))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
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
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow(nil))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
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
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "places"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "place_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow(nil))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
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
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
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
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemPlace(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestSaveMediaItemThing(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemThingRequest
		MockDB      func(mock sqlmock.Sqlmock)
		ExpectedErr error
	}{
		{
			"save mediaitem thing with invalid mediaitem user id",
			&api.MediaItemThingRequest{UserId: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem thing with invalid mediaitem id",
			&api.MediaItemThingRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", Id: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem thing with success",
			&mediaItemThingRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnRows(getMockedThingRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "things"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "things"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "thing_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow(&existingPlaceKeywords))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem thing with thing find or create error",
			&mediaItemThingRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(codes.Internal, "error getting or creating thing: some db error"),
		},
		{
			"save mediaitem thing with error",
			&mediaItemThingRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnRows(getMockedThingRow())
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "things"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "things"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving mediaitem thing: some db error"),
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
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemThing(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestSaveMediaItemFaces(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemFacesRequest
		MockDB      func(mock sqlmock.Sqlmock)
		ExpectedErr error
	}{
		{
			"save mediaitem faces with invalid mediaitem user id",
			&api.MediaItemFacesRequest{UserId: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem faces with invalid mediaitem id",
			&api.MediaItemFacesRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", Id: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem faces with success",
			&mediaItemFacesRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitem_faces"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem faces with error",
			&mediaItemFacesRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitem_faces"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving mediaitem faces: some db error"),
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
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemFaces(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestGetMediaItemFaceEmbeddings(t *testing.T) {
	tests := []struct {
		Name           string
		Request        *api.MediaItemFaceEmbeddingsRequest
		MockDB         func(mock sqlmock.Sqlmock)
		ExpectedResult *api.MediaItemFaceEmbeddingsResponse
		ExpectedErr    error
	}{
		{
			"get mediaitem face embeddings with invalid mediaitem user id",
			&api.MediaItemFaceEmbeddingsRequest{UserId: "bad-mediaitem-user-id"},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"get mediaitem face embeddings with success",
			&mediaItemFaceEmbeddingsRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow(&existingPlaceKeywords))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitem_faces"`)).
					WillReturnRows(getMockedMediaItemFaceEmbeddingRows())
			},
			&api.MediaItemFaceEmbeddingsResponse{
				MediaItemFaceEmbeddings: []*api.MediaItemFaceEmbedding{
					{
						MediaItemId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
						Embedding:   &mediaItemEmbedding,
					},
				},
			},
			nil,
		},
		{
			"get mediaitem face embeddings with error",
			&mediaItemFaceEmbeddingsRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			status.Error(codes.Internal, "error getting mediaitem face embeddings: some db error"),
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
			client := api.NewAPIClient(conn)
			res, err := client.GetMediaItemFaceEmbeddings(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
			if test.ExpectedResult != nil {
				for idx, mediaItemFaceEmbedding := range test.ExpectedResult.MediaItemFaceEmbeddings {
					assert.Equal(t, mediaItemFaceEmbedding.Embedding.Embedding, res.MediaItemFaceEmbeddings[idx].Embedding.Embedding)
					assert.Equal(t, mediaItemFaceEmbedding.MediaItemId, res.MediaItemFaceEmbeddings[idx].MediaItemId)
				}
			} else {
				assert.Equal(t, test.ExpectedResult, res)
			}
		})
	}
}

func TestSaveMediaItemPeople(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemPeopleRequest
		MockDB      func(mock sqlmock.Sqlmock)
		ExpectedErr error
	}{
		{
			"save mediaitem people with invalid mediaitem user id",
			&api.MediaItemPeopleRequest{UserId: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem people with invalid mediaitem id",
			&api.MediaItemPeopleRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", MediaItemFacePeople: map[string]*api.MediaItemFacePeople{"bad-mediaitem-id": nil}},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem people with invalid face id",
			&api.MediaItemPeopleRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", MediaItemFacePeople: map[string]*api.MediaItemFacePeople{"4d05b5f6-17c2-475e-87fe-3fc8b9567179": &api.MediaItemFacePeople{
				FacePeople: map[string]string{"bad-face-id": "bad-people-id"},
			}}},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid face id"),
		},
		{
			"save mediaitem people with success",
			&mediaItemPeopleRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "people"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitem_faces"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitem_faces"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem people with error saving people",
			&mediaItemPeopleRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving people: some db error"),
		},
		{
			"save mediaitem people with error saving people mediaitems",
			&mediaItemPeopleRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "people"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving people mediaitems: some db error"),
		},
		{
			"save mediaitem people with error saving mediaitem faces people",
			&mediaItemPeopleRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "people"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "people_mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitem_faces"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving mediaitem faces people: some db error"),
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
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemPeople(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestSaveMediaItemFinalResult(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemFinalResultRequest
		MockDB      func(mock sqlmock.Sqlmock)
		ExpectedErr error
	}{
		{
			"save mediaitem ml result with invalid mediaitem user id",
			&api.MediaItemFinalResultRequest{UserId: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem ml result with invalid mediaitem id",
			&api.MediaItemFinalResultRequest{UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", Id: "bad-mediaitem-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem final result with success",
			&mediaItemFinalResultRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitem_embeddings"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem final result with error",
			&mediaItemFinalResultRequest,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitem_embeddings"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			status.Error(codes.Internal, "error saving mediaitem final result: some db error"),
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
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemFinalResult(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func dialer(service *Service) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	api.RegisterAPIServer(server, service)
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

func getMockedThingRow() *sqlmock.Rows {
	return sqlmock.NewRows(thingCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime)
}

func getMockedMediaItemRow(existingKeyword *string) *sqlmock.Rows {
	return sqlmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"filename", "description", "mime_type", existingKeyword, "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", "mediaitem_category", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime)
}

func getMockedMediaItemFaceEmbeddingRows() *sqlmock.Rows {
	return sqlmock.NewRows(mediaitemFaceCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil, embedding).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179", embedding)
}

func getMockedUserIDRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id"}).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180")
}
