package handlers

import (
	"api/config"
	"api/internal/models"
	"api/pkg/cache"
	"api/pkg/services/worker"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bluele/gcache"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
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
	MockDB           func(mock pgxmock.PgxPoolIface)
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
			ctx.Set("userID", "019b7796-6072-76ee-8be3-485ff2b32fd7")
			if _, ok := test.Header[echo.HeaderAuthorization]; ok {
				var features models.Features
				_ = json.Unmarshal([]byte(`{"albums":true,"explore":true,"places":true,"people":true}`), &features)
				ctx.Set("features", features)
			}
			// database
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()
			if test.MockDB != nil {
				test.MockDB(mockDB)
			}
			mockCache := &cache.InMemoryCache{
				Connection: gcache.New(1024).LRU().Build(),
			}
			if test.mockCache != nil {
				mockCache = &cache.InMemoryCache{Connection: gcache.New(1024).
					LRU().
					SerializeFunc(test.mockCache[0]).
					DeserializeFunc(test.mockCache[1]).
					Build()}
			}
			for key, val := range test.Header {
				if key == echo.HeaderAuthorization {
					mockCache.SetWithExpire(val, true, 1*time.Minute)
				}
			}
			// handler
			handler := &Handler{
				Config: &config.Config{
					Storage: config.Storage{DiskRoot: os.TempDir()}, Auth: config.Auth{RefreshTTL: 60}, Feature: config.Feature{
						Albums: true, Explore: true, Places: true, People: true,
					}, ML: config.ML{
						Places: true, Faces: true, Search: true,
					},
				}, DB: mockDB, Cache: mockCache, Worker: test.mockWorkerClient,
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
		wantErr bool
		wantOk  bool
	}
)

func (mwc *mockWorkerGRPCClient) GenerateEmbedding(ctx context.Context, request *worker.GenerateEmbeddingRequest, opts ...grpc.CallOption) (*worker.GenerateEmbeddingResponse, error) {
	if mwc.wantErr {
		return nil, errors.New("some grpc error")
	}
	return &worker.GenerateEmbeddingResponse{Embedding: make([]float32, 0)}, nil
}

type AnyID struct{}

func (a AnyID) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
