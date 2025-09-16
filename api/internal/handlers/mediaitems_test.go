package handlers

import (
	"api/pkg/services/api"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	sampleTime, _           = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	sampleDescription       = "description"
	sampleBoolTrue          = true
	sampleBoolFalse         = false
	sampleCameraMake        = "camera_make"
	sampleCameraModel       = "camera_model"
	sampleFocalLength       = "focal_length"
	sampleApertureFnumber   = "aperture_fnumber"
	sampleIsoEquivalent     = "iso_equivalent"
	sampleExposureTime      = "exposure_time"
	sampleMegapixels        = "18.4"
	sampleLatitude          = 17.580249
	sampleLongitude         = -70.278493
	sampleFPS               = "fps"
	sampleSourceURL         = "source_url"
	samplePreviewURL        = "preview_url"
	sampleThumbnailURL      = "thumbnail_url"
	samplePlaceholder       = "placeholder"
	sampleMediaItemType     = "mediaitem_type"
	sampleMediaItemCategory = "mediaitem_category"
	sampleWidth             = 720
	sampleHeight            = 480

	coverMediaItemCols = []string{
		"id", "user_id", "source_url", "preview_url", "thumbnail_url", "placeholder", "mediaitem_type",
		"mediaitem_category", "width", "height",
	}
	mediaitemCols = []string{
		"id", "user_id", "filename", "hash", "description", "mime_type", "source_url", "preview_url", "thumbnail_url",
		"placeholder", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "mediaitem_category",
		"width", "height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "megapixels", "latitude", "longitude", "fps", "exif_data",
		"keywords", "created_at", "updated_at",
	}
	mediaitemFaceCols = []string{
		"id", "mediaitem_id", "people_id", "embeddings", "thumbnail",
	}
	mediaitemResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","favourite":true,"hidden":false,"deleted":false,"status":"READY",` +
		`"mediaItemType":"PHOTO","mediaItemCategory":"DEFAULT","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"megapixels":"18.4","latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	mediaitemsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","favourite":true,"hidden":false,"deleted":false,"status":"READY",` +
		`"mediaItemType":"PHOTO","mediaItemCategory":"DEFAULT","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"megapixels":"18.4","latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","favourite":false,"hidden":true,"deleted":true,"status":"READY",` +
		`"mediaItemType":"VIDEO","mediaItemCategory":"DEFAULT","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"megapixels":"18.4","latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetMediaItemPlaces(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem places bad request", http.MethodGet, "/v1/mediaItems/:id/places", "/v1/mediaItems/bad-uuid/places", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"get mediaitem places with empty table", http.MethodGet, "/v1/mediaItems/:id/places", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(placeCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			}, http.StatusOK, "[]",
		},
		{
			"get mediaitem places with error", http.MethodGet, "/v1/mediaItems/:id/places", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get mediaitem places with error in scanning", http.MethodGet, "/v1/mediaItems/:id/places", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(placeCols, coverMediaItemCols...)).
						AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &samplePostCode, &sampleCountry, &sampleLocality, &sampleArea, &sampleBoolTrue, &sampleCoverMediaItemID, sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get mediaitem places with success", http.MethodGet, "/v1/mediaItems/:id/places", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedPlaceRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			}, http.StatusOK, placesResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemThings(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem things bad request", http.MethodGet, "/v1/mediaItems/:id/things", "/v1/mediaItems/bad-uuid/things", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"get mediaitem things with empty table", http.MethodGet, "/v1/mediaItems/:id/things", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM things t LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(thingCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			}, http.StatusOK, "[]",
		},
		{
			"get mediaitem things with error", http.MethodGet, "/v1/mediaItems/:id/things", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM things t LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get mediaitem things with error in scanning", http.MethodGet, "/v1/mediaItems/:id/things", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM things t LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(thingCols, coverMediaItemCols...)).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleCoverMediaItemID, "true", sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get mediaitem things with success", http.MethodGet, "/v1/mediaItems/:id/things", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM things t LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedThingRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			}, http.StatusOK, thingsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemPeople(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem people bad request", http.MethodGet, "/v1/mediaItems/:id/people", "/v1/mediaItems/bad-uuid/people", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"get mediaitem people with empty table", http.MethodGet, "/v1/mediaItems/:id/people", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			}, http.StatusOK, "[]",
		},
		{
			"get mediaitem people with error", http.MethodGet, "/v1/mediaItems/:id/people", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get mediaitem people with error in scanning", http.MethodGet, "/v1/mediaItems/:id/people", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)).
						AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleBoolTrue, &sampleCoverMediaItemID, &sampleCoverMediaItemID, sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", &sampleCoverMediaItemID, nil, "thumbnail"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get mediaitem people with success", http.MethodGet, "/v1/mediaItems/:id/people", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedPeopleRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			}, http.StatusOK, peopleResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemAlbums(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem albums bad request", http.MethodGet, "/v1/mediaItems/:id/albums", "/v1/mediaItems/bad-uuid/albums", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"get mediaitem albums with empty table", http.MethodGet, "/v1/mediaItems/:id/albums", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums a LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			}, http.StatusOK, "[]",
		},
		{
			"get mediaitem albums with error", http.MethodGet, "/v1/mediaItems/:id/albums", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums a LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get mediaitem albums with error in scanning", http.MethodGet, "/v1/mediaItems/:id/albums", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums a LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, coverMediaItemCols...)).
						AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription, &sampleBoolTrue, &sampleBoolFalse, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get mediaitem albums with success", http.MethodGet, "/v1/mediaItems/:id/albums", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums a LEFT JOIN mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedAlbumRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			}, http.StatusOK, albumsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItem(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem bad request", http.MethodGet, "/v1/mediaItems/:id", "/v1/mediaItems/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"get mediaitem not found", http.MethodGet, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			}, http.StatusNotFound, "mediaitem not found",
		},
		{
			"get mediaitem with error", http.MethodGet, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get mediaitem with error in scanning", http.MethodGet, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get mediaitem with success", http.MethodGet, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRow())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			}, http.StatusOK, mediaitemResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestUpdateMediaItem(t *testing.T) {
	tests := []Test{
		{
			"update mediaitem bad request", http.MethodPut, "/v1/mediaItems/:id", "/v1/mediaItems/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"update mediaitem with no payload", http.MethodPut, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			}, http.StatusBadRequest, "invalid mediaitem",
		},
		{
			"update mediaitem with bad payload", http.MethodPut, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			}, http.StatusBadRequest, "invalid mediaitem",
		},
		{
			"update mediaitem with error", http.MethodPut, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"description":"description","favourite":true,"hidden":true}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(sampleCoverMediaItemID, sampleCoverMediaItemID, &sampleDescription, &sampleBoolTrue, &sampleBoolTrue, pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"update mediaitem with success", http.MethodPut, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"description":"description","favourite":true,"hidden":true}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(sampleCoverMediaItemID, sampleCoverMediaItemID, &sampleDescription, &sampleBoolTrue, &sampleBoolTrue, pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteMediaItem(t *testing.T) {
	tests := []Test{
		{
			"delete mediaitem bad request", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"delete mediaitem with error deleting", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error starting transaction", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{}).WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error getting album new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error scanning album new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "invalid"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"delete mediaitem with error getting place new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error scanning place new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "invalid"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"delete mediaitem with error getting thing new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error scanning thing new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "invalid"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"delete mediaitem with error getting people new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error scanning people new cover mediaitems", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "invalid"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"delete mediaitem with error updating album cover mediaitem", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error updating place cover mediaitem", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error updating thing cover mediaitem", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE things`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error updating people cover mediaitem", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE things`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with error committing transaction", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE things`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete mediaitem with success", http.MethodDelete, "/v1/mediaItems/:id", "/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"album_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"place_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"thing_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"people_id", "mediaitem_id"}).AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE things`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get mediaitems with empty table", http.MethodGet, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			}, http.StatusOK, "[]",
		},
		{
			"get mediaitems with error", http.MethodGet, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get mediaitems with error in scanning", http.MethodGet, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get mediaitems with 2 rows", http.MethodGet, "/v1/mediaItems", "/v1/mediaItems?type=PHOTO&category=PANORAMA", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			}, http.StatusOK, mediaitemsResponseBody,
		},
		{
			"get mediaitems with 2 rows and filters", http.MethodGet, "/v1/mediaItems", "/v1/mediaItems?status=FAILED", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			}, http.StatusOK, mediaitemsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestUploadMediaItems(t *testing.T) {
	sampleFile, contentType := getMockedMediaItemFile(t)
	sampleFile2, contentType2 := getMockedMediaItemFile(t)
	sampleFile3, contentType3 := getMockedMediaItemFile(t)
	sampleFile4, contentType4 := getMockedMediaItemFile(t)
	sampleFile5, contentType5 := getMockedMediaItemFile(t)
	sampleFile6, contentType6 := getMockedMediaItemFile(t)
	tests := []Test{
		{
			"upload mediaitems with invalid command", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				HeaderUploadType: "resumable",
			}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusBadRequest, "invalid command for resumable upload",
		},
		{
			"upload mediaitems with invalid offset", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				HeaderUploadType: "resumable", HeaderUploadCommand: "finish",
			}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusBadRequest, "invalid chunk offset for resumable upload",
		},
		{
			"upload mediaitems with invalid session", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				HeaderUploadType: "resumable", HeaderUploadCommand: "finish", HeaderUploadChunkOffset: "1024",
			}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusBadRequest, "invalid chunk session for resumable upload",
		},
		{
			"upload mediaitems with error uploading for resumable", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				HeaderUploadType: "resumable", HeaderUploadCommand: "start", HeaderUploadChunkOffset: "0",
			}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusBadRequest, "request Content-Type isn't multipart/form-data",
		},
		{
			"upload mediaitems with error uploading", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusBadRequest, "request Content-Type isn't multipart/form-data",
		},
		{
			"upload mediaitems with error inserting mediaitem", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: contentType,
			}, sampleFile, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"upload mediaitems with error saving hash due to duplicates", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: contentType2,
			}, sampleFile2, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("violates unique constraint"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusConflict, "mediaitem already exists",
		},
		{
			"upload mediaitems with error saving hash in mediaitem process", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: contentType3,
			}, sampleFile3, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"upload mediaitems with error queuing for processing", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: contentType4, echo.HeaderAuthorization: "atoken",
			}, sampleFile4, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO queue`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"upload mediaitems successfully", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: contentType5, echo.HeaderAuthorization: "atoken",
			}, sampleFile5, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO queue`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusCreated, `"id"`,
		},
		{
			"upload mediaitems successfully for resumable", http.MethodPost, "/v1/mediaItems", "/v1/mediaItems", []string{}, []string{}, map[string]string{
				HeaderUploadType: "resumable", HeaderUploadCommand: "finish", HeaderUploadChunkOffset: "100", HeaderUploadChunkSession: "4d05b5f6-17c2-475e-87fe-3fc8b9567179", echo.HeaderContentType: contentType6,
			}, sampleFile6, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO queue`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			}, http.StatusNoContent, ``,
		},
	}
	executeTests(t, tests)
}

func getMockedMediaItemRow() *pgxmock.Rows {
	return pgxmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse,
			api.MediaItemStatus_READY.String(), api.MediaItemType_PHOTO.String(), api.MediaItemCategory_DEFAULT.String(), 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime)
}

func getMockedMediaItemRows() *pgxmock.Rows {
	return pgxmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse,
			api.MediaItemStatus_READY.String(), api.MediaItemType_PHOTO.String(), api.MediaItemCategory_DEFAULT.String(), 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolFalse, &sampleBoolTrue, &sampleBoolTrue,
			api.MediaItemStatus_READY.String(), api.MediaItemType_VIDEO.String(), api.MediaItemCategory_DEFAULT.String(), 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime)
}

func getMockedMediaItemFile(t *testing.T) (*io.PipeReader, string) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := writer.CreateFormFile("file", "image.png")
		if err != nil {
			t.Error(err)
		}

		// create sample image
		dim := 10
		upLeft := image.Point{0, 0}
		lowRight := image.Point{dim, dim}
		img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
		for x := 0; x < dim; x++ {
			for y := 0; y < dim; y++ {
				img.Set(x, y, color.White)
			}
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	return pr, writer.FormDataContentType()
}
