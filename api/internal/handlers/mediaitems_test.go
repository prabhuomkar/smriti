package handlers

import (
	"api/config"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	sampleTime, _  = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	mediaitem_cols = []string{"id", "filename", "description", "mime_type", "source_url", "preview_url",
		"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "width",
		"height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "location", "fps", "created_at", "updated_at"}
	mediaitem_rows = sqlmock.NewRows(mediaitem_cols).
			AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "location", "fps", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "false", "true", "true", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "location", "fps", sampleTime, sampleTime)
)

func TestGetMediaItemPlaces(t *testing.T) {

}

func TestGetMediaItemThings(t *testing.T) {

}

func TestGetMediaItemPeople(t *testing.T) {

}

func TestGetMediaItem(t *testing.T) {

}

func TestUpdateMediaItem(t *testing.T) {

}

func TestDeleteMediaItem(t *testing.T) {

}

func TestGetMediaItems(t *testing.T) {
	tests := []struct {
		Name            string
		Path            string
		MockDB          func(mock sqlmock.Sqlmock)
		ExpectedErr     bool
		ExpectedResCode int
		ExpectedResBody string
	}{
		{
			"get mediaitems with empty table",
			"/v1/mediaItems",
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			false,
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitems with 2 rows",
			"/v1/mediaItems",
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`))
				expectedQuery.WillReturnRows(mediaitem_rows)
			},
			false,
			http.StatusOK,
			`[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename","description":"description",` +
				`"mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url","thumbnailUrl":"thumbnail_url",` +
				`"status":"status","mediaItemType":"mediaitem_type","width":720,"height":480,` +
				`"creationTime":"2022-09-22T11:22:33+05:30","cameraMake":"camera_make","cameraModel":"camera_model",` +
				`"focalLength":"focal_length","apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent",` +
				`"exposureTime":"exposure_time","location":"bG9jYXRpb24=","fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
				`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","filename":"filename",` +
				`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
				`"thumbnailUrl":"thumbnail_url","status":"status","mediaItemType":"mediaitem_type","width":720,"height":480,` +
				`"creationTime":"2022-09-22T11:22:33+05:30","cameraMake":"camera_make","cameraModel":"camera_model",` +
				`"focalLength":"focal_length","apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent",` +
				`"exposureTime":"exposure_time","location":"bG9jYXRpb24=","fps":"fps",` +
				`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}]`,
		},
		{
			"get mediaitems with error",
			"/v1/mediaItems",
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			true,
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// server
			server := echo.New()
			req := httptest.NewRequest(http.MethodGet, test.Path, nil)
			rec := httptest.NewRecorder()

			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()
			mockDBx := sqlx.NewDb(mockDB, "sqlmock")
			test.MockDB(mock)

			// handler
			handler := &Handler{
				Config: &config.Config{},
				DB:     mockDBx,
			}
			server.GET(test.Path, handler.GetMediaItems)
			server.ServeHTTP(rec, req)

			assert.Equal(t, test.ExpectedResCode, rec.Code)
			assert.Equal(t, test.ExpectedResBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestUploadMediaItems(t *testing.T) {

}
