package handlers

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/worker"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluele/gcache"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Test struct {
	Name             string
	Method           string
	Route            string
	Path             string
	ParamNames       []string
	ParamValues      []string
	Header           map[string]string
	Body             io.Reader
	MockDB           func(mock sqlmock.Sqlmock)
	mockCache        []func(interface{}, interface{}) (interface{}, error)
	mockWorkerClient *mockWorkerGRPCClient
	Handler          func(handler *Handler) func(ctx echo.Context) error
	ExpectedResCode  int
	ExpectedResBody  string
}

func executeTests(t *testing.T, tests []Test) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// server
			server := echo.New()
			req := httptest.NewRequest(test.Method, test.Path, test.Body)
			for key, val := range test.Header {
				if key == echo.HeaderAuthorization {
					req.Header.Set(key, fmt.Sprintf("Bearer %s", val))
				} else {
					req.Header.Set(key, val)
				}
			}
			rec := httptest.NewRecorder()
			// context
			ctx := server.NewContext(req, rec)
			ctx.SetPath(test.Route)
			ctx.SetParamNames(test.ParamNames...)
			ctx.SetParamValues(test.ParamValues...)
			ctx.Set("userID", "4d05b5f6-17c2-475e-87fe-3fc8b9567179")
			if _, ok := test.Header[echo.HeaderAuthorization]; ok {
				var features models.Features
				_ = json.Unmarshal([]byte(`{"albums":true,"explore":true,"places":true,"ml":{"places":true,"faces":true}}`), &features)
				ctx.Set("features", features)
			}
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
			mockCache := gcache.New(1024).LRU().Build()
			if test.mockCache != nil {
				mockCache = gcache.New(1024).
					LRU().
					SerializeFunc(test.mockCache[0]).
					DeserializeFunc(test.mockCache[1]).
					Build()
			}
			for key, val := range test.Header {
				if key == echo.HeaderAuthorization {
					mockCache.Set(val, true)
				}
			}
			// handler
			handler := &Handler{
				Config: &config.Config{
					Auth: config.Auth{RefreshTTL: 60},
					Feature: config.Feature{Albums: true, Explore: true, Places: true,
						ML: config.MLFeatures{Places: true, Faces: true},
					},
				},
				DB:     mockGDB,
				Cache:  mockCache,
				Worker: test.mockWorkerClient,
			}
			err = test.Handler(handler)(ctx)
			if test.ExpectedResCode >= http.StatusBadRequest {
				assert.Equal(t, test.ExpectedResCode, err.(*echo.HTTPError).Code)
				assert.Contains(t, err.(*echo.HTTPError).Message.(string), test.ExpectedResBody)
			} else {
				assert.Equal(t, test.ExpectedResCode, rec.Code)
				assert.Contains(t, strings.TrimSpace(rec.Body.String()), test.ExpectedResBody)
			}
		})
	}
}

type (
	mockWorkerGRPCClient struct {
		wantErr             bool
		wantSendErr         bool
		wantCloseAndRecvErr bool
		wantOk              bool
	}

	mockWorkerMediaItemProcessClient struct {
		wantSendErr         bool
		wantCloseAndRecvErr bool
		wantOk              bool
	}
)

func (mwmpc *mockWorkerMediaItemProcessClient) Send(*worker.MediaItemProcessRequest) error {
	if mwmpc.wantSendErr {
		return errors.New("some error in send")
	}
	return nil
}
func (mwmpc *mockWorkerMediaItemProcessClient) CloseAndRecv() (*worker.MediaItemProcessResponse, error) {
	if mwmpc.wantCloseAndRecvErr {
		return nil, errors.New("some error in close and recv")
	}
	return &worker.MediaItemProcessResponse{Ok: mwmpc.wantOk}, nil
}

func (mwmpc *mockWorkerMediaItemProcessClient) Header() (metadata.MD, error) { return nil, nil }

func (mwmpc *mockWorkerMediaItemProcessClient) Trailer() metadata.MD { return nil }

func (mwmpc *mockWorkerMediaItemProcessClient) CloseSend() error { return nil }

func (mwmpc *mockWorkerMediaItemProcessClient) Context() context.Context { return nil }

func (mwmpc *mockWorkerMediaItemProcessClient) SendMsg(m interface{}) error { return nil }

func (mwmpc *mockWorkerMediaItemProcessClient) RecvMsg(m interface{}) error { return nil }

func (mwc *mockWorkerGRPCClient) MediaItemProcess(ctx context.Context, opts ...grpc.CallOption) (worker.Worker_MediaItemProcessClient, error) {
	if mwc.wantErr {
		return nil, errors.New("some grpc error")
	}
	return &mockWorkerMediaItemProcessClient{wantSendErr: mwc.wantSendErr, wantCloseAndRecvErr: mwc.wantCloseAndRecvErr, wantOk: mwc.wantOk}, nil
}
