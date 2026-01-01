package service

import (
	"api/config"
	"api/pkg/services/api"
	"api/pkg/storage"
	"context"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"testing"

	"github.com/google/uuid"
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
	mediaitemType                = api.MediaItemType_PHOTO.String()
	mediaitemCategory            = api.MediaItemCategory_DEFAULT.String()
	previewUrl                   = "preview_url"
	sourceUrl                    = "source_url"
	latitude                     = "latitude"
	longitude                    = "longitude"
	components                   = "METADATA,PLACES"
	badcreationtime              = "bad-creation-time"
	creationtime                 = "2022-09-22 11:22:33"
	width                  int32 = 1080
	height                 int32 = 720
	placeholder                  = "placeholder"
	sampleId, _                  = uuid.Parse("019b7796-6072-76ee-8be3-485ff2b32fd7")
	sampleEmbedding              = pgvector.NewVector([]float32{0.42})
	mediaItemResultRequest       = api.MediaItemMetadataRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		MimeType: &mimetype, Type: api.MediaItemType_PHOTO, Category: api.MediaItemCategory_DEFAULT, Width: &width, Height: &height, CreationTime: &creationtime,
	}
	mediaItemPreviewThumbnailRequest = api.MediaItemPreviewThumbnailRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		Status: api.MediaItemStatus_READY, Placeholder: &placeholder,
	}
	country            = "country"
	postcode           = "postcode"
	locality           = "locality"
	area               = "area"
	mediaItemEmbedding = api.MediaItemEmbedding{
		Embedding: []float32{0.0, 0.42, 0.111},
	}
	mediaItemPlaceRequest = api.MediaItemPlaceRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		Country: &country, Postcode: &postcode, Locality: &locality, Area: &area,
	}
	mediaItemPlaceLocalityRequest = api.MediaItemPlaceRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		Locality: &locality,
	}
	mediaItemPlaceAreaRequest = api.MediaItemPlaceRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		Area: &area,
	}
	mediaItemFacesRequest = api.MediaItemFacesRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		Embeddings: []*api.MediaItemEmbedding{&mediaItemEmbedding}, Thumbnails: []string{"thumbnail"},
	}
	mediaItemPeopleRequest = api.MediaItemPeopleRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
			"019b7796-6072-76ee-8be3-485ff2b32fd7": {
				FacePeople: map[string]string{
					"019b7796-6072-76ee-8be3-485ff2b32fd7": "019b7796-6072-76ee-8be3-485ff2b32fd7",
				},
			}, "019b7796-6072-76ee-8be3-485ff2b33fd7": {
				FacePeople: map[string]string{
					"019b7796-6072-76ee-8be3-485ff2b33fd7": "019b7796-6072-76ee-8be3-485ff2b32fd7",
				},
			}, "019b7796-6072-76ee-8be3-485ff2b34fd7": {
				FacePeople: map[string]string{
					"019b7796-6072-76ee-8be3-485ff2b34fd7": "1",
				},
			}, "019b7796-6072-76ee-8be3-485ff2b35fd7": {
				FacePeople: map[string]string{
					"019b7796-6072-76ee-8be3-485ff2b35fd7": "1",
				},
			},
		},
	}
	mediaItemFaceEmbeddingsRequest = api.MediaItemFaceEmbeddingsRequest{
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
	}
	mediaItemFinalResultRequest = api.MediaItemFinalResultRequest{
		Id:     "019b7796-6072-76ee-8be3-485ff2b32fd7",
		UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7",
		DetectedText: "some detected text", Caption: "some caption",
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
			"get worker config with success", &config.Config{
				ML: config.ML{Places: true, PlacesProvider: "openstreetmap"},
			}, []byte(`[{"name":"METADATA"},{"name":"PLACES","source":"openstreetmap"}]`), nil,
		},
		{
			"get worker config with success with all config", &config.Config{ML: config.ML{
				Places: true, PlacesProvider: "openstreetmap", OCR: true, OCRProvider: "paddlepaddle",
				OCRParams: `{"det_model_dir":"/det_infer"}`, Search: true, SearchProvider: "pytorch",
				SearchParams: `{"tokenizer_dir":"/tokenizer"}`, Faces: true, FacesParams: `{"face_threshold":"0.9"}`,
				PreviewThumbnailParams: `{"thumbnail_size":"256"}`,
			}}, []byte(`[{"name":"METADATA"},{"name":"PREVIEW_THUMBNAIL","params":"{\"thumbnail_size\":\"256\"}"},` +
				`{"name":"PLACES","source":"openstreetmap"},{"name":"OCR","source":"paddlepaddle","params":` +
				`"{\"det_model_dir\":\"/det_infer\"}"},{"name":"SEARCH","source":"pytorch","params":` +
				`"{\"tokenizer_dir\":\"/tokenizer\"}"},{"name":"FACES","params":"{\"face_threshold\":\"0.9\"}"}]`), nil,
		},
		{
			"get worker config with no error", &config.Config{}, []byte(`[{"name":"METADATA"}]`), nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// service
			service := Init(test.Config, nil, nil)
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

func TestGetMediaItemProcess(t *testing.T) {
	tests := []struct {
		Name           string
		MockDB         func(mock pgxmock.PgxPoolIface)
		ExpectedResult *api.MediaItemProcessResponse
		ExpectedErr    error
	}{
		{
			"get mediaitem to process with error", func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE queue`)).
					WillReturnError(errors.New("some db error"))
			}, nil, status.Error(codes.Internal, "error getting mediaitem to process: some db error"),
		},
		{
			"get mediaitem to process with error due to empty rows", func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE queue`)).
					WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "mediaitem_id", "components", "mime_type", "source_url",
						"preview_url", "mediaitem_type", "mediaitem_category", "latitude", "longitude"}))
			}, &api.MediaItemProcessResponse{}, nil,
		},
		{
			"get mediaitem to process with success", func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE queue`)).
					WillReturnRows(getMockedMediaItemToProcessRow())
			}, &api.MediaItemProcessResponse{
				Id:          "019b7796-6072-76ee-8be3-485ff2b32fd7",
				UserId:      "019b7796-6072-76ee-8be3-485ff2b33fd7",
				MediaItemId: "019b7796-6072-76ee-8be3-485ff2b34fd7",
				Components:  []api.MediaItemComponent{api.MediaItemComponent_METADATA, api.MediaItemComponent_PLACES},
				Payload: map[string]string{"category": mediaitemCategory, "latitude": "latitude",
					"longitude": "longitude", "mime_type": "mimetype",
					"preview_url": "preview_url", "source_url": "source_url", "type": mediaitemType},
			}, nil,
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
			service := Init(&config.Config{ML: config.ML{Places: true}}, mockDB, nil)
			// server
			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer(service)))
			assert.Nil(t, err)
			defer conn.Close()
			client := api.NewAPIClient(conn)
			res, err := client.GetMediaItemProcess(ctx, &emptypb.Empty{})
			// assert
			assert.Equal(t, test.ExpectedErr, err)
			if test.ExpectedErr == nil {
				assert.Equal(t, test.ExpectedResult.Id, res.Id)
				assert.Equal(t, test.ExpectedResult.UserId, res.UserId)
				assert.Equal(t, test.ExpectedResult.MediaItemId, res.MediaItemId)
				assert.Equal(t, test.ExpectedResult.Components, res.Components)
				assert.Equal(t, test.ExpectedResult.Payload, res.Payload)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		Name           string
		MockDB         func(mock pgxmock.PgxPoolIface)
		ExpectedResult *api.UsersResponse
		ExpectedErr    error
	}{
		{
			"get users with error", func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users`)).
					WillReturnError(errors.New("some db error"))
			}, nil, status.Error(codes.Internal, "error getting users: some db error"),
		},
		{
			"get users with error due to scanning", func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users`)).
					WillReturnRows(getMockedUserIDRows(true))
			}, nil, status.Error(codes.Internal, "error scanning user: Scanning value error for column 'id': Scan: invalid UUID length: 7"),
		},
		{
			"get users with success", func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users`)).
					WillReturnRows(getMockedUserIDRows(false))
			}, &api.UsersResponse{
				Users: []string{
					"019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b33fd7",
				},
			}, nil,
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
			service := Init(&config.Config{}, mockDB, nil)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
		ExpectedErr error
	}{
		{
			"save mediaitem result with invalid mediaitem user id", &api.MediaItemMetadataRequest{UserId: "bad-mediaitem-user-id"}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem result with invalid mediaitem id", &api.MediaItemMetadataRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "bad-mediaitem-id",
			}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem result with incorrect creation time", &api.MediaItemMetadataRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7", CreationTime: &badcreationtime,
			}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem creation time"),
		},
		{
			"save mediaitem result with error", &mediaItemResultRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error updating mediaitem result: some db error"),
		},
		{
			"save mediaitem result with success", &mediaItemResultRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}, nil,
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
			service := Init(&config.Config{}, mockDB, &storage.Disk{Root: tmpRoot})
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

func TestSaveMediaItemPreviewThumbnail(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemPreviewThumbnailRequest
		MockDB      func(mock pgxmock.PgxPoolIface)
		MockParams  func(string) (string, string, string, func(), error)
		ExpectedErr error
	}{
		{
			"save mediaitem preview and thumbnail with invalid mediaitem user id", &api.MediaItemPreviewThumbnailRequest{
				UserId: "bad-mediaitem-user-id",
			}, nil, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem preview and thumbnail with invalid mediaitem id", &api.MediaItemPreviewThumbnailRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "bad-mediaitem-id",
			}, nil, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem preview and thumbnail with error uploading original file",
			&mediaItemPreviewThumbnailRequest, nil, func(tmpRoot string) (string, string, string, func(), error) {
				return "", "", "", func() {}, nil
			}, status.Errorf(codes.Internal, "error uploading original file"),
		},
		{
			"save mediaitem preview and thumbnail with error uploading preview file",
			&mediaItemPreviewThumbnailRequest, nil, func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0o777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0o777)
				return originalFile.Name(), "", "", func() {
					defer os.Remove(tmpRoot + "/originals/")
					defer os.Remove(originalFile.Name())
				}, nil
			}, status.Errorf(codes.Internal, "error uploading preview file"),
		},
		{
			"save mediaitem preview and thumbnail with error uploading thumbnail file",
			&mediaItemPreviewThumbnailRequest, nil, func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0o777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0o777)
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
			}, status.Errorf(codes.Internal, "error uploading thumbnail file"),
		},
		{
			"save mediaitem preview and thumbnail with error", &mediaItemPreviewThumbnailRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0o777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0o777)
				previewFile, err := os.CreateTemp(tmpRoot, "preview")
				assert.NoError(t, err)
				os.Mkdir(tmpRoot+"/thumbnails/", 0o777)
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
			}, status.Error(codes.Internal, "error updating mediaitem result: some db error"),
		},
		{
			"save mediaitem preview and thumbnail with success", &mediaItemPreviewThumbnailRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}, func(tmpRoot string) (string, string, string, func(), error) {
				os.Mkdir(tmpRoot+"/originals/", 0o777)
				originalFile, err := os.CreateTemp(tmpRoot, "original")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/previews/", 0o777)
				previewFile, err := os.CreateTemp(tmpRoot, "preview")
				if err != nil {
					return "", "", "", nil, err
				}
				os.Mkdir(tmpRoot+"/thumbnails/", 0o777)
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
			}, nil,
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
			service := Init(&config.Config{}, mockDB, &storage.Disk{Root: tmpRoot})
			// mock tmp params
			if test.MockParams != nil {
				originalPath, previewPath, thumbnailPath, clear, err := test.MockParams(tmpRoot)
				assert.NoError(t, err)
				test.Request.SourcePath = &originalPath
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
			"save mediaitem place with invalid mediaitem user id", &api.MediaItemPlaceRequest{UserId: "bad-mediaitem-id"}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem place with invalid mediaitem id", &api.MediaItemPlaceRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "bad-mediaitem-id",
			}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem place with error starting transaction", &mediaItemPlaceRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{}).WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error starting transaction for saving mediaitem place: some db error"),
		},
		{
			"save mediaitem place with error saving place", &mediaItemPlaceRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving place: some db error"),
		},
		{
			"save mediaitem place with error saving mediaitem place", &mediaItemPlaceRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving mediaitem place: some db error"),
		},
		{
			"save mediaitem place with error committing transaction", &mediaItemPlaceRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error committing transaction for saving mediaitem place: some db error"),
		},
		{
			"save mediaitem place with all details success", &mediaItemPlaceRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			}, nil,
		},
		{
			"save mediaitem place with locality success", &mediaItemPlaceLocalityRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			}, nil,
		},
		{
			"save mediaitem place with area success", &mediaItemPlaceAreaRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			}, nil,
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
			service := Init(&config.Config{}, mockDB, nil)
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

func TestSaveMediaItemFaces(t *testing.T) {
	tests := []struct {
		Name        string
		Request     *api.MediaItemFacesRequest
		MockDB      func(mock pgxmock.PgxPoolIface)
		MockParams  func(string) (string, func(), error)
		ExpectedErr error
	}{
		{
			"save mediaitem faces with invalid mediaitem user id", &api.MediaItemFacesRequest{UserId: "bad-mediaitem-id"}, nil, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem faces with invalid mediaitem id", &api.MediaItemFacesRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemId: "bad-mediaitem-id",
			}, nil, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem faces with error uploading face thumbnail file", &mediaItemFacesRequest, nil, nil,
			status.Error(codes.Internal, "error uploading mediaitem face thumbnail file"),
		},
		{
			"save mediaitem faces with error", &mediaItemFacesRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, func(tmpRoot string) (string, func(), error) {
				os.Mkdir(tmpRoot+"/faces/", 0o777)
				faceThumbnailFile, err := os.CreateTemp(tmpRoot, "thumbnail")
				if err != nil {
					return "", nil, err
				}
				return faceThumbnailFile.Name(), func() {
					defer os.Remove(tmpRoot + "/faces/")
					defer os.Remove(faceThumbnailFile.Name())
				}, nil
			}, status.Error(codes.Internal, "error saving mediaitem faces: some db error"),
		},
		{
			"save mediaitem faces with success", &mediaItemFacesRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}, func(tmpRoot string) (string, func(), error) {
				os.Mkdir(tmpRoot+"/faces/", 0o777)
				faceThumbnailFile, err := os.CreateTemp(tmpRoot, "thumbnail")
				if err != nil {
					return "", nil, err
				}
				return faceThumbnailFile.Name(), func() {
					defer os.Remove(tmpRoot + "/faces/")
					defer os.Remove(faceThumbnailFile.Name())
				}, nil
			}, nil,
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
			service := Init(&config.Config{}, mockDB, &storage.Disk{Root: tmpRoot})
			// mock tmp params
			if test.MockParams != nil {
				facesPath, clear, err := test.MockParams(tmpRoot)
				assert.NoError(t, err)
				thumbnails := make([]string, len(mediaItemFacesRequest.Thumbnails))
				for idx := range test.Request.Thumbnails {
					thumbnails[idx] = facesPath
				}
				test.Request.Thumbnails = thumbnails
				defer clear()
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
		MockDB         func(mock pgxmock.PgxPoolIface)
		ExpectedResult *api.MediaItemFaceEmbeddingsResponse
		ExpectedErr    error
	}{
		{
			"get mediaitem face embeddings with invalid mediaitem user id", &api.MediaItemFaceEmbeddingsRequest{
				UserId: "bad-mediaitem-user-id",
			}, nil, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"get mediaitem face embeddings with error", &mediaItemFaceEmbeddingsRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mediaitem_id, people_id, embedding FROM mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, status.Error(codes.Internal, "error getting mediaitem face embeddings: some db error"),
		},
		{
			"get mediaitem face embeddings with error due to scanning", &mediaItemFaceEmbeddingsRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mediaitem_id, people_id, embedding FROM mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemFaceEmbeddingRows(true))
			}, nil, status.Error(codes.Internal, "error scanning mediaitem face embedding: Scanning value error for column 'id': Scan: invalid UUID length: 7"),
		},
		{
			"get mediaitem face embeddings with success", &mediaItemFaceEmbeddingsRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, mediaitem_id, people_id, embedding FROM mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemFaceEmbeddingRows(false))
			}, &api.MediaItemFaceEmbeddingsResponse{
				MediaItemFaceEmbeddings: []*api.MediaItemFaceEmbedding{
					{
						MediaItemId: "019b7796-6072-76ee-8be3-485ff2b32fd7", Embedding: nil,
					},
				},
			}, nil,
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
			service := Init(&config.Config{}, mockDB, nil)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
		ExpectedErr error
	}{
		{
			"save mediaitem people with invalid mediaitem user id", &api.MediaItemPeopleRequest{UserId: "bad-mediaitem-id"}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem people with invalid mediaitem id", &api.MediaItemPeopleRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
					"bad-mediaitem-id": nil,
				},
			}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem people with invalid face id", &api.MediaItemPeopleRequest{
				UserId: "019b7796-6072-76ee-8be3-485ff2b32fd7", MediaItemFacePeople: map[string]*api.MediaItemFacePeople{
					"019b7796-6072-76ee-8be3-485ff2b32fd7": {
						FacePeople: map[string]string{
							"bad-face-id": "bad-people-id",
						},
					},
				},
			}, nil, status.Errorf(codes.InvalidArgument, "invalid face id"),
		},
		{
			"save mediaitem people with error starting transaction", &mediaItemPeopleRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{}).WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error starting transaction for saving mediaitem people: some db error"),
		},
		{
			"save mediaitem people with error saving people", &mediaItemPeopleRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving people: some db error"),
		},
		{
			"save mediaitem people with error saving people mediaitems", &mediaItemPeopleRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving people mediaitems: some db error"),
		},
		{
			"save mediaitem people with error saving mediaitem faces", &mediaItemPeopleRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving mediaitem faces: some db error"),
		},
		{
			"save mediaitem people with error committing transaction", &mediaItemPeopleRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error committing transaction for saving mediaitem people: some db error"),
		},
		{
			"save mediaitem people with success", &mediaItemPeopleRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			}, nil,
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
			service := Init(&config.Config{}, mockDB, nil)
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
		MockDB      func(mock pgxmock.PgxPoolIface)
		ExpectedErr error
	}{
		{
			"save mediaitem ml result with invalid queue id", &api.MediaItemFinalResultRequest{Id: "bad-queue-id"}, nil, status.Errorf(codes.InvalidArgument, "invalid queue id"),
		},
		{
			"save mediaitem ml result with invalid user id", &api.MediaItemFinalResultRequest{
				Id:     "019b7796-6072-76ee-8be3-485ff2b32fd7",
				UserId: "bad-mediaitem-user-id",
			}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem user id"),
		},
		{
			"save mediaitem ml result with invalid mediaitem id", &api.MediaItemFinalResultRequest{
				Id:          "019b7796-6072-76ee-8be3-485ff2b32fd7",
				UserId:      "019b7796-6072-76ee-8be3-485ff2b32fd7",
				MediaItemId: "bad-mediaitem-id",
			}, nil, status.Errorf(codes.InvalidArgument, "invalid mediaitem id"),
		},
		{
			"save mediaitem final result with error saving final result", &mediaItemFinalResultRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving mediaitem final result: some db error"),
		},
		{
			"save mediaitem final result with error saving embeddings", &mediaItemFinalResultRequest, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitem_embeddings`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error saving mediaitem final result embedding: some db error"),
		},
		{
			"save mediaitem final result with error unqueuing mediaitem from processing", &mediaItemFinalResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitem_embeddings`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM queue`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, status.Error(codes.Internal, "error unqueuing mediaitem from processing: some db error"),
		},
		{
			"save mediaitem final result with success", &mediaItemFinalResultRequest,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitem_embeddings`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM queue`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			}, nil,
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
			service := Init(&config.Config{}, mockDB, nil)
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

func getMockedMediaItemFaceEmbeddingRows(bad bool) *pgxmock.Rows {
	if bad {
		return pgxmock.NewRows(mediaitemFaceCols).
			AddRow("invalid", "invalid", nil, nil)
	}
	return pgxmock.NewRows(mediaitemFaceCols).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleId, &sampleEmbedding).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleId, &sampleEmbedding)
}

func getMockedMediaItemToProcessRow() *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id", "user_id", "mediaitem_id", "components", "mime_type", "source_url",
		"preview_url", "mediaitem_type", "mediaitem_category", "latitude", "longitude"}).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b33fd7",
			"019b7796-6072-76ee-8be3-485ff2b34fd7", components, &mimetype, sourceUrl,
			&previewUrl, &mediaitemType, &mediaitemCategory, &latitude, &longitude)

}

func getMockedUserIDRows(bad bool) *pgxmock.Rows {
	if bad {
		return pgxmock.NewRows([]string{"id"}).
			AddRow("invalid")
	}
	return pgxmock.NewRows([]string{"id"}).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7").
		AddRow("019b7796-6072-76ee-8be3-485ff2b33fd7")
}
