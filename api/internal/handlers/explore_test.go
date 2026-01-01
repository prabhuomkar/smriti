package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	samplePostCode = "postcode"
	sampleCountry  = "country"
	sampleLocality = "locality"
	sampleArea     = "area"

	placeCols = []string{
		"id", "user_id", "name", "postcode", "country", "locality", "area", "is_hidden", "cover_mediaitem_id", "created_at", "updated_at",
	}
	peopleCols = []string{
		"id", "user_id", "name", "is_hidden", "cover_mediaitem_id", "cover_mediaitem_face_id", "created_at", "updated_at",
	}
	memoryMediaItemCols          = append(mediaitemCols, "creation_year")
	memoryMediaItemsResponseBody = `[{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","mediaItemCategory":"mediaitem_category","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"megapixels":"18.4","latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30","year":"2023"},{"id":"019b7796-6072-76ee-8be3-485ff2b33fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","favourite":false,"hidden":true,"deleted":true,"status":"status",` +
		`"mediaItemType":"mediaitem_type","mediaItemCategory":"mediaitem_category","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"megapixels":"18.4","latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30","year":"2022"}]`
	coverMediaItemResponseBody = `"coverMediaItem":{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","mediaItemType":"mediaitem_type",` +
		`"mediaItemCategory":"mediaitem_category","width":720,"height":480}`
	coverFaceResponseBody = `"coverMediaItemFace":{"thumbnail":"thumbnail"}`
	placeResponseBody     = `{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7","userId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"name":"name","postcode":"postcode","country":"country","locality":"locality","area":"area","hidden":true,` +
		`"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `}`
	placesResponseBody = `[{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7","userId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"name":"name","postcode":"postcode","country":"country","locality":"locality","area":"area","hidden":true,"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `},{"id":"019b7796-6072-76ee-8be3-485ff2b33fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","name":"name",` +
		`"postcode":"postcode","country":"country","locality":"locality","area":"area","hidden":false,` +
		`"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30",` + coverMediaItemResponseBody + `}]`
	personResponseBody = `{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","name":"name",` +
		`"hidden":true,"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7","coverMediaItemFaceId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverFaceResponseBody + `}`
	peopleResponseBody = `[{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","name":"name",` +
		`"hidden":true,"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7","coverMediaItemFaceId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverFaceResponseBody + `},{"id":"019b7796-6072-76ee-8be3-485ff2b33fd7",` +
		`"userId":"019b7796-6072-76ee-8be3-485ff2b32fd7","name":"name",` +
		`"hidden":false,"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7","coverMediaItemFaceId":"019b7796-6072-76ee-8be3-485ff2b32fd7",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` + coverFaceResponseBody + `}]`
)

func TestGetYearsAgoMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get years ago mediaitems bad request", http.MethodGet, "/v1/explore/yearsAgo/:monthDate/mediaItems", "/v1/explore/yearsAgo/bad-month-date/mediaItems", []string{"monthDate"}, []string{"bad-month-date"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetYearsAgoMediaItems
			}, http.StatusBadRequest, "invalid month and date",
		},
		{
			"get years ago mediaitems not found", http.MethodGet, "/v1/explore/yearsAgo/:monthDate/mediaItems", "/v1/explore/yearsAgo/0403/mediaItems", []string{"monthDate"}, []string{"0403"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), "04", "03").
					WillReturnRows(pgxmock.NewRows(memoryMediaItemCols))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetYearsAgoMediaItems
			}, http.StatusOK, "",
		},
		{
			"get years ago mediaitems with error", http.MethodGet, "/v1/explore/yearsAgo/:monthDate/mediaItems", "/v1/explore/yearsAgo/0403/mediaItems", []string{"monthDate"}, []string{"0403"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), "04", "03").
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetYearsAgoMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get years ago mediaitems with error in scanning", http.MethodGet, "/v1/explore/yearsAgo/:monthDate/mediaItems", "/v1/explore/yearsAgo/0403/mediaItems", []string{"monthDate"}, []string{"0403"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), "04", "03").
					WillReturnRows(pgxmock.NewRows(memoryMediaItemCols).AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, nil, sampleTime, sampleTime, "2023"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetYearsAgoMediaItems
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get years ago mediaitems with 2 years", http.MethodGet, "/v1/explore/yearsAgo/:monthDate/mediaItems", "/v1/explore/yearsAgo/0403/mediaItems", []string{"monthDate"}, []string{"0403"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), "04", "03").
					WillReturnRows(getMockedMemoryMediaItemRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetYearsAgoMediaItems
			}, http.StatusOK, memoryMediaItemsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetPlaces(t *testing.T) {
	tests := []Test{
		{
			"get places with empty table", http.MethodGet, "/v1/explore/places", "/v1/explore/places", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(placeCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			}, http.StatusOK, "[]",
		},
		{
			"get places with error", http.MethodGet, "/v1/explore/places", "/v1/explore/places", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get places with error in scanning", http.MethodGet, "/v1/explore/places", "/v1/explore/places", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(placeCols, coverMediaItemCols...)).
						AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &samplePostCode, &sampleCountry, &sampleLocality, &sampleArea, &sampleBoolTrue, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get places with 2 rows", http.MethodGet, "/v1/explore/places", "/v1/explore/places", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedPlaceRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			}, http.StatusOK, placesResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetPlace(t *testing.T) {
	tests := []Test{
		{
			"get place bad request", http.MethodGet, "/v1/explore/places/:id", "/v1/explore/places/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			}, http.StatusBadRequest, "invalid place id",
		},
		{
			"get place not found", http.MethodGet, "/v1/explore/places/:id", "/v1/explore/places/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(placeCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			}, http.StatusNotFound, "place not found",
		},
		{
			"get place with error", http.MethodGet, "/v1/explore/places/:id", "/v1/explore/places/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get place with error in scanning", http.MethodGet, "/v1/explore/places/:id", "/v1/explore/places/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(placeCols, coverMediaItemCols...)).
						AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &samplePostCode, &sampleCountry, &sampleLocality, &sampleArea, &sampleBoolTrue, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get place with success", http.MethodGet, "/v1/explore/places/:id", "/v1/explore/places/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedPlaceRow())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			}, http.StatusOK, placeResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetPlaceMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get place mediaitems bad request", http.MethodGet, "/v1/places/:id/mediaItems", "/v1/places/bad-uuid/mediaItems", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			}, http.StatusBadRequest, "invalid place id",
		},
		{
			"get place mediaitems not found", http.MethodGet, "/v1/places/:id/mediaItems", "/v1/places/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			}, http.StatusOK, "[]",
		},
		{
			"get place mediaitems with error", http.MethodGet, "/v1/places/:id/mediaItems", "/v1/places/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get place mediaitems with error in scanning", http.MethodGet, "/v1/places/:id/mediaItems", "/v1/places/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, nil, sampleTime, sampleTime))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get place mediaitems with 2 rows", http.MethodGet, "/v1/places/:id/mediaItems", "/v1/places/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			}, http.StatusOK, mediaitemsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestUpdatePerson(t *testing.T) {
	tests := []Test{
		{
			"update people bad request", http.MethodPut, "/v1/people/:id", "/v1/people/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			}, http.StatusBadRequest, "invalid people id",
		},
		{
			"update people with no payload", http.MethodPut, "/v1/people/:id", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			}, http.StatusBadRequest, "invalid people",
		},
		{
			"update people with bad payload", http.MethodPut, "/v1/people/:id", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			}, http.StatusBadRequest, "invalid people",
		},
		{
			"update people with bad cover mediaitem id", http.MethodPut, "/v1/people/:id", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","coverMediaItemId":"bad-mediaitem-id"}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			}, http.StatusBadRequest, "invalid people cover mediaitem id",
		},
		{
			"update people with error", http.MethodPut, "/v1/people/:id", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","hidden":true,"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "name", &sampleBoolTrue, &sampleCoverMediaItemID, pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"update people with success", http.MethodPut, "/v1/people/:id", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","hidden":true,"coverMediaItemId":"019b7796-6072-76ee-8be3-485ff2b32fd7"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "name", &sampleBoolTrue, &sampleCoverMediaItemID, pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetPeople(t *testing.T) {
	tests := []Test{
		{
			"get people with empty table", http.MethodGet, "/v1/explore/people", "/v1/explore/people", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			}, http.StatusOK, "[]",
		},
		{
			"get people with error", http.MethodGet, "/v1/explore/people", "/v1/explore/people", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get people with error in scanning", http.MethodGet, "/v1/explore/people", "/v1/explore/people", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)).AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &sampleBoolTrue, &sampleCoverMediaItemID, &sampleCoverMediaItemID, sampleTime, sampleTime, "019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleCoverMediaItemID, nil, "thumbnail"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get people with 2 rows", http.MethodGet, "/v1/explore/people", "/v1/explore/people", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedPeopleRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			}, http.StatusOK, peopleResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetPerson(t *testing.T) {
	tests := []Test{
		{
			"get person bad request", http.MethodGet, "/v1/explore/people/:id", "/v1/explore/people/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			}, http.StatusBadRequest, "invalid person id",
		},
		{
			"get person not found", http.MethodGet, "/v1/explore/people/:id", "/v1/explore/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			}, http.StatusNotFound, "person not found",
		},
		{
			"get person with error", http.MethodGet, "/v1/explore/people/:id", "/v1/explore/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get person with error in scanning", http.MethodGet, "/v1/explore/people/:id", "/v1/explore/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)).AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &sampleBoolTrue, &sampleCoverMediaItemID, &sampleCoverMediaItemID, sampleTime, sampleTime, "019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleCoverMediaItemID, nil, "thumbnail"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get people with success", http.MethodGet, "/v1/explore/people/:id", "/v1/explore/people/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.*, mf.* FROM people`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedPeopleRow())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			}, http.StatusOK, personResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestGetPersonMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get people mediaitems bad request", http.MethodGet, "/v1/people/:id/mediaItems", "/v1/people/bad-uuid/mediaItems", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPersonMediaItems
			}, http.StatusBadRequest, "invalid people id",
		},
		{
			"get people mediaitems not found", http.MethodGet, "/v1/people/:id/mediaItems", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPersonMediaItems
			}, http.StatusOK, "[]",
		},
		{
			"get people mediaitems with error", http.MethodGet, "/v1/people/:id/mediaItems", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPersonMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get people mediaitems with error in scanning", http.MethodGet, "/v1/people/:id/mediaItems", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).AddRow("invalid", "019b7796-6072-76ee-8be3-485ff2b32fd7", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, nil, sampleTime, sampleTime))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPersonMediaItems
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get people mediaitems with 2 rows", http.MethodGet, "/v1/people/:id/mediaItems", "/v1/people/019b7796-6072-76ee-8be3-485ff2b32fd7/mediaItems", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPersonMediaItems
			}, http.StatusOK, mediaitemsResponseBody,
		},
	}
	executeTests(t, tests)
}

func getMockedPlaceRow() *pgxmock.Rows {
	return pgxmock.NewRows(append(placeCols, coverMediaItemCols...)).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &samplePostCode, &sampleCountry, &sampleLocality, &sampleArea, &sampleBoolTrue, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight)
}

func getMockedPlaceRows() *pgxmock.Rows {
	return pgxmock.NewRows(append(placeCols, coverMediaItemCols...)).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &samplePostCode, &sampleCountry, &sampleLocality, &sampleArea, &sampleBoolTrue, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight).
		AddRow("019b7796-6072-76ee-8be3-485ff2b33fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &samplePostCode, &sampleCountry, &sampleLocality, &sampleArea, &sampleBoolFalse, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight)
}

func getMockedPeopleRow() *pgxmock.Rows {
	return pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &sampleBoolTrue, &sampleCoverMediaItemID, &sampleCoverMediaItemID, sampleTime, sampleTime, "019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleCoverMediaItemID, nil, "thumbnail")
}

func getMockedPeopleRows() *pgxmock.Rows {
	return pgxmock.NewRows(append(peopleCols, mediaitemFaceCols...)).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &sampleBoolTrue, &sampleCoverMediaItemID, &sampleCoverMediaItemID, sampleTime, sampleTime, "019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleCoverMediaItemID, nil, "thumbnail").
		AddRow("019b7796-6072-76ee-8be3-485ff2b33fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "name", &sampleBoolFalse, &sampleCoverMediaItemID, &sampleCoverMediaItemID, sampleTime, sampleTime, "019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", &sampleCoverMediaItemID, nil, "thumbnail")
}

func getMockedMemoryMediaItemRows() *pgxmock.Rows {
	return pgxmock.NewRows(memoryMediaItemCols).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, nil, sampleTime, sampleTime, "2023").
		AddRow("019b7796-6072-76ee-8be3-485ff2b33fd7", "019b7796-6072-76ee-8be3-485ff2b32fd7", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolFalse, &sampleBoolTrue, &sampleBoolTrue, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, nil, sampleTime, sampleTime, "2022")
}
