package service

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"testing"

	"api/config"
	"api/internal/models"
	"api/pkg/services/api"
	"api/pkg/storage"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	mimetype                     = "mimetype"
	mediaitemType                = "photo"
	mediaitemCategory            = "default"
	badcreationtime              = "bad-creation-time"
	creationtime                 = "2022-09-22 11:22:33"
	width                  int32 = 1080
	height                 int32 = 720
	existingPlaceKeywords        = "placecity placepostcode"
	placeholder                  = "placeholder"
	mediaItemResultRequest       = api.MediaItemMetadataRequest{
		UserId:       "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:           "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		MimeType:     &mimetype,
		Type:         mediaitemType,
		Category:     mediaitemCategory,
		Width:        &width,
		Height:       &height,
		CreationTime: &creationtime,
	}
	mediaItemPreviewThumbnailRequest = api.MediaItemPreviewThumbnailRequest{
		UserId:      "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:          "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Status:      string(models.StatusReady),
		Placeholder: &placeholder,
	}
	country            = "country"
	postcode           = "postcode"
	locality           = "locality"
	area               = "area"
	embedding          = pgvector.NewVector([]float32{0.0, 0.42, 0.111})
	mediaItemEmbedding = api.MediaItemEmbedding{
		Embedding: []float32{0.0, 0.42, 0.111},
	}
	mediaItemPlaceRequest = api.MediaItemPlaceRequest{
		UserId:   "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:       "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Country:  &country,
		Postcode: &postcode,
		Locality: &locality,
		Area:     &area,
	}
	mediaItemPlaceLocalityRequest = api.MediaItemPlaceRequest{
		UserId:   "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:       "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Locality: &locality,
	}
	mediaItemPlaceAreaRequest = api.MediaItemPlaceRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Id:     "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		Area:   &area,
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
		Thumbnails: []string{"thumbnail"},
	}
	mediaItemPeopleRequest = api.MediaItemPeopleRequest{
		UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
		MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179": {
				FacePeople: map[string]string{
					"4d05b5f6-17c2-475e-87fe-3fc8b9567179": "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				},
			},
			"4d05b5f6-17c2-475e-87fe-3fc8b9567180": {
				FacePeople: map[string]string{
					"4d05b5f6-17c2-475e-87fe-3fc8b9567180": "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				},
			},
			"4d05b5f6-17c2-475e-87fe-3fc8b9567181": {
				FacePeople: map[string]string{
					"4d05b5f6-17c2-475e-87fe-3fc8b9567181": "1",
				},
			},
			"4d05b5f6-17c2-475e-87fe-3fc8b9567182": {
				FacePeople: map[string]string{
					"4d05b5f6-17c2-475e-87fe-3fc8b9567182": "1",
				},
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
			&config.Config{
				ML: config.ML{Places: true, PlacesProvider: "openstreetmap"},
			},
			[]byte(
				`[{"name":"METADATA"},{"name":"PLACES","source":"openstreetmap"}]`,
			),
			nil,
		},
		{
			"get worker config with success with all config",
			&config.Config{ML: config.ML{
				Places: true, PlacesProvider: "openstreetmap",
				Classification: true, ClassificationProvider: "pytorch", ClassificationParams: `{"file":"model-file-name.pt"}`,
				OCR: true, OCRProvider: "paddlepaddle", OCRParams: `{"det_model_dir":"/det_infer"}`,
				Search: true, SearchProvider: "pytorch", SearchParams: `{"tokenizer_dir":"/tokenizer"}`,
				Faces: true, FacesParams: `{"face_threshold":"0.9"}`,
				PreviewThumbnailParams: `{"thumbnail_size":"256"}`,
			}},
			[]byte(
				`[{"name":"METADATA"},{"name":"PREVIEW_THUMBNAIL","params":"{\"thumbnail_size\":\"256\"}"},{"name":"PLACES","source":"openstreetmap"},` +
					`{"name":"CLASSIFICATION","source":"pytorch","params":"{\"file\":\"model-file-name.pt\"}"},{"name":"OCR","source":"paddlepaddle",` +
					`"params":"{\"det_model_dir\":\"/det_infer\"}"},{"name":"SEARCH","source":"pytorch","params":"{\"tokenizer_dir\":\"/tokenizer\"}"},` +
					`{"name":"FACES","params":"{\"face_threshold\":\"0.9\"}"}]`,
			),
			nil,
		},
		{
			"get worker config with no error",
			&config.Config{},
			[]byte(`[{"name":"METADATA"}]`),
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
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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
		MockDB         func(mock pgxmock.PgxPoolIface)
		ExpectedResult *api.GetUsersResponse
		ExpectedErr    error
	}{
		{
			"get users with error",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			status.Error(codes.Internal, "error getting users: some db error"),
		},
		{
			"get users with success",
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users`)).
					WillReturnRows(getMockedUserIDRows())
			},
			&api.GetUsersResponse{
				Users: []string{
					"4d05b5f6-17c2-475e-87fe-3fc8b9567179",
					"4d05b5f6-17c2-475e-87fe-3fc8b9567180",
				},
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
		ExpectedErr error
	}{
		{
			"save mediaitem result with invalid mediaitem user id",
			&api.MediaItemMetadataRequest{UserId: "bad-mediaitem-user-id"},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem result with invalid mediaitem id",
			&api.MediaItemMetadataRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "bad-mediaitem-id",
			},
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
			status.Errorf(
				codes.InvalidArgument,
				"invalid mediaitem creation time",
			),
		},
		{
			"save mediaitem result with error",
			&mediaItemResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(
				codes.Internal,
				"error updating mediaitem result: some db error",
			),
		},
		{
			"save mediaitem result with success",
			&mediaItemResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			tmpRoot := os.TempDir()
			service := &Service{
				Config:  &config.Config{},
				DB:      mockDB,
				Storage: &storage.Disk{Root: tmpRoot},
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemMetadata(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestSaveMediaItemPreviewThumbnail(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemPreviewThumbnailRequest
		MockDB      func(mock pgxmock.PgxPoolIface)
		MockParams  func(string) (string, string, string, func(), error)
		ExpectedErr error
	}{
		{
			"save mediaitem preview and thumbnail with invalid mediaitem user id",
			&api.MediaItemPreviewThumbnailRequest{
				UserId: "bad-mediaitem-user-id",
			},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem preview and thumbnail with invalid mediaitem id",
			&api.MediaItemPreviewThumbnailRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "bad-mediaitem-id",
			},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem preview and thumbnail with error uploading original file",
			&mediaItemPreviewThumbnailRequest,
			nil,
			func(tmpRoot string) (string, string, string, func(), error) {
				return "", "", "", func() {}, nil
			},
			status.Errorf(codes.Internal, "error uploading original file"),
		},
		{
			"save mediaitem preview and thumbnail with error uploading preview file",
			&mediaItemPreviewThumbnailRequest,
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
			"save mediaitem preview and thumbnail with error uploading thumbnail file",
			&mediaItemPreviewThumbnailRequest,
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
			"save mediaitem preview and thumbnail with error",
			&mediaItemPreviewThumbnailRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("some db error"))
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
			status.Error(
				codes.Internal,
				"error updating mediaitem result: some db error",
			),
		},
		{
			"save mediaitem preview and thumbnail with success",
			&mediaItemPreviewThumbnailRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
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
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			tmpRoot := os.TempDir()
			service := &Service{
				Config:  &config.Config{},
				DB:      mockDB,
				Storage: &storage.Disk{Root: tmpRoot},
			}
			// mock tmp params
			if test.MockParams != nil {
				originalPath, previewPath, thumbnailPath, clear, err := test.MockParams(
					tmpRoot,
				)
				assert.NoError(t, err)
				test.Request.SourcePath = &originalPath
				test.Request.PreviewPath = &previewPath
				test.Request.ThumbnailPath = &thumbnailPath
				defer clear()
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIClient(conn)
			_, err = client.SaveMediaItemPreviewThumbnail(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestSaveMediaItemPlace(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemPlaceRequest
		MockDB      func(mock pgxmock.PgxPoolIface)
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
			&api.MediaItemPlaceRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "bad-mediaitem-id",
			},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem place with error starting transaction",
			&mediaItemPlaceRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(
					pgx.TxOptions{},
				).WillReturnError(
					errors.New("some db error"),
				)
			},
			status.Error(
				codes.Internal,
				"error starting transaction for saving mediaitem place: some db error",
			),
		},
		{
			"save mediaitem place with error saving place",
			&mediaItemPlaceRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("some db error"))

			},
			status.Error(codes.Internal, "error saving place: some db error"),
		},
		{
			"save mediaitem place with error saving mediaitem place",
			&mediaItemPlaceRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO place_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(
						errors.New("some db error"),
					)

			},
			status.Error(
				codes.Internal,
				"error saving mediaitem place: some db error",
			),
		},
		{
			"save mediaitem place with error committing transaction",
			&mediaItemPlaceRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO place_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))

			},
			status.Error(
				codes.Internal,
				"error committing transaction for saving mediaitem place: some db error",
			),
		},
		{
			"save mediaitem place with all details success",
			&mediaItemPlaceRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO place_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem place with locality success",
			&mediaItemPlaceLocalityRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO place_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectCommit()
			},
			nil,
		},
		{
			"save mediaitem place with area success",
			&mediaItemPlaceAreaRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO place_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectCommit()
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
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
			&api.MediaItemThingRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "bad-mediaitem-id",
			},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem thing with error starting transaction",
			&mediaItemThingRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(
					pgx.TxOptions{},
				).WillReturnError(
					errors.New("some db error"),
				)
			},
			status.Error(
				codes.Internal,
				"error starting transaction for saving mediaitem thing: some db error",
			),
		},
		{
			"save mediaitem thing with error saving thing",
			&mediaItemThingRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO things`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(codes.Internal, "error saving thing: some db error"),
		},
		{
			"save mediaitem thing with error saving mediaitem thing",
			&mediaItemThingRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO things`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO thing_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(
						errors.New("some db error"),
					)
			},
			status.Error(
				codes.Internal,
				"error saving mediaitem thing: some db error",
			),
		},
		{
			"save mediaitem thing with error committing transaction",
			&mediaItemThingRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO things`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO thing_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			},
			status.Error(
				codes.Internal,
				"error committing transaction for saving mediaitem thing: some db error",
			),
		},
		{
			"save mediaitem thing with success",
			&mediaItemThingRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO things`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO thing_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectCommit()
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
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
			&api.MediaItemFacesRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "bad-mediaitem-id",
			},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem faces with error",
			&mediaItemFacesRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO mediaitem_faces`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(
						errors.New("some db error"),
					)
			},
			status.Error(
				codes.Internal,
				"error saving mediaitem faces: some db error",
			),
		},
		{
			"save mediaitem faces with success",
			&mediaItemFacesRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO mediaitem_faces`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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
		MockDB         func(mock pgxmock.PgxPoolIface)
		ExpectedResult *api.MediaItemFaceEmbeddingsResponse
		ExpectedErr    error
	}{
		{
			"get mediaitem face embeddings with invalid mediaitem user id",
			&api.MediaItemFaceEmbeddingsRequest{
				UserId: "bad-mediaitem-user-id",
			},
			nil,
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"get mediaitem face embeddings with error",
			&mediaItemFaceEmbeddingsRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT id, mediaitem_id, people_id, embedding FROM mediaitem_faces`,
					),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(
						errors.New("some db error"),
					)
			},
			nil,
			status.Error(
				codes.Internal,
				"error getting mediaitem face embeddings: some db error",
			),
		},
		{
			"get mediaitem face embeddings with success",
			&mediaItemFaceEmbeddingsRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT id, mediaitem_id, people_id, embedding FROM mediaitem_faces`,
					),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnRows(
						getMockedMediaItemFaceEmbeddingRows(),
					)
			},
			&api.MediaItemFaceEmbeddingsResponse{
				MediaItemFaceEmbeddings: []*api.MediaItemFaceEmbedding{
					{
						MediaItemId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
						Embedding:   nil,
					},
				},
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIClient(conn)
			res, err := client.GetMediaItemFaceEmbeddings(ctx, test.Request)
			// assert
			assert.Equal(t, test.ExpectedErr, err)
			if test.ExpectedResult != nil {
				for idx, mediaItemFaceEmbedding := range test.ExpectedResult.MediaItemFaceEmbeddings {
					assert.Equal(
						t,
						mediaItemFaceEmbedding.MediaItemId,
						res.MediaItemFaceEmbeddings[idx].MediaItemId,
					)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
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
			&api.MediaItemPeopleRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
					"bad-mediaitem-id": nil,
				},
			},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem people with invalid face id",
			&api.MediaItemPeopleRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
					"4d05b5f6-17c2-475e-87fe-3fc8b9567179": {
						FacePeople: map[string]string{
							"bad-face-id": "bad-people-id",
						},
					},
				},
			},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid face id"),
		},
		{
			"save mediaitem people with error starting transaction",
			&mediaItemPeopleRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(
					pgx.TxOptions{},
				).WillReturnError(
					errors.New("some db error"),
				)
			},
			status.Error(
				codes.Internal,
				"error starting transaction for saving mediaitem people: some db error",
			),
		},
		{
			"save mediaitem people with error saving people",
			&mediaItemPeopleRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(codes.Internal, "error saving people: some db error"),
		},
		{
			"save mediaitem people with error saving people mediaitems",
			&mediaItemPeopleRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(
						errors.New("some db error"),
					)
			},
			status.Error(
				codes.Internal,
				"error saving people mediaitems: some db error",
			),
		},
		{
			"save mediaitem people with error saving mediaitem faces",
			&mediaItemPeopleRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(
				codes.Internal,
				"error saving mediaitem faces: some db error",
			),
		},
		{
			"save mediaitem people with error committing transaction",
			&mediaItemPeopleRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			},
			status.Error(
				codes.Internal,
				"error committing transaction for saving mediaitem people: some db error",
			),
		},
		{
			"save mediaitem people with success",
			&mediaItemPeopleRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO people_mediaitems`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
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
			&api.MediaItemFinalResultRequest{
				UserId: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				Id:     "bad-mediaitem-id",
			},
			nil,
			status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem final result with error saving keywords",
			&mediaItemFinalResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("some db error"))
			},
			status.Error(
				codes.Internal,
				"error saving mediaitem final result keywords: some db error",
			),
		},
		{
			"save mediaitem final result with error saving embeddings",
			&mediaItemFinalResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO mediaitem_embeddings`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnError(
						errors.New("some db error"),
					)
			},
			status.Error(
				codes.Internal,
				"error saving mediaitem final result embedding: some db error",
			),
		},
		{
			"save mediaitem final result with success",
			&mediaItemFinalResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO mediaitem_embeddings`),
				).
					WithArgs(
						pgxmock.AnyArg(),
						pgxmock.AnyArg(),
					).
					WillReturnResult(
						pgxmock.NewResult("INSERT", 1),
					)
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			// service
			service := &Service{
				Config: &config.Config{},
				DB:     mockDB,
			}
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(
				ctx,
				"",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)),
			)
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

func getMockedMediaItemFaceEmbeddingRows() *pgxmock.Rows {
	return pgxmock.NewRows(mediaitemFaceCols).
		AddRow(
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			nil,
		).
		AddRow(
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			nil,
		)
}

func getMockedUserIDRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id"}).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180")
}
