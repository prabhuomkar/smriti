package handlers

import (
	"api/internal/models"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
)

var (
	jobCols = []string{
		"id", "user_id", "components", "status", "last_mediaitem_id", "created_at", "updated_at",
	}
	jobResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"status":"SCHEDULED","components":"metadata,places","lastMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}`
	jobsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"status":"RUNNING","components":"metadata,places","lastMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"},` +
		`{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"status":"RUNNING","components":"faces","lastMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetJob(t *testing.T) {
	tests := []Test{
		{
			"get job bad request",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJob
			},
			http.StatusBadRequest,
			"invalid job id",
		},
		{
			"get job not found",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(jobCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJob
			},
			http.StatusNotFound,
			"job not found",
		},
		{
			"get job with error",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJob
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"get job with error in scanning",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(jobCols).AddRow(pgxmock.AnyArg(), pgxmock.AnyArg(), "SCHEDULED", "metadata,places",
						pgxmock.AnyArg(), sampleTime, sampleTime))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJob
			},
			http.StatusInternalServerError,
			"Scanning value error",
		},
		{
			"get job with success",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedJobRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJob
			},
			http.StatusOK,
			jobResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestUpdateJob(t *testing.T) {
	tests := []Test{
		{
			"update job bad request",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusBadRequest,
			"invalid job id",
		},
		{
			"update job with no payload",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusBadRequest,
			"invalid job",
		},
		{
			"update job with bad payload",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusBadRequest,
			"invalid job",
		},
		{
			"update job with error getting existing job count",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"status":"RUNNING"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"update job with error due to existing job",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"status":"RUNNING"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusConflict,
			"job already exists",
		},
		{
			"update job with error",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"status":"RUNNING"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE jobs`)).
					WithArgs(sampleCoverMediaItemID, sampleCoverMediaItemID, models.JobRunning, pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"update job with success",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"status":"RUNNING"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE jobs`)).
					WithArgs(sampleCoverMediaItemID, sampleCoverMediaItemID, models.JobRunning, pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusNoContent,
			"",
		},
	}
	executeTests(t, tests)
}

func TestGetJobs(t *testing.T) {
	tests := []Test{
		{
			"get jobs with empty table",
			http.MethodGet,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(jobCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJobs
			},
			http.StatusOK,
			"[]",
		},
		{
			"get jobs with error",
			http.MethodGet,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJobs
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"get jobs with error in scanning",
			http.MethodGet,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(jobCols).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "SCHEDULED", "metadata,places",
						&sampleCoverMediaItemID, sampleTime, sampleTime))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJobs
			},
			http.StatusInternalServerError,
			"Scanning value error",
		},
		{
			"get jobs with 2 rows",
			http.MethodGet,
			"/v1/jobs",
			"/v1/jobs?sort=updatedAt",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedJobRows())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetJobs
			},
			http.StatusOK,
			jobsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestCreateJob(t *testing.T) {
	tests := []Test{
		{
			"create job with bad payload",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request}`),
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusBadRequest,
			"invalid job",
		},
		{
			"create job with no payload",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusBadRequest,
			"invalid job",
		},
		{
			"create job with error getting existing job count",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"components":"search"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"create job with error due to existing job",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"components":"search"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusConflict,
			"job already exists",
		},
		{
			"create job with error",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},

			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"components":"search"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), models.JobScheduled, "search", pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"create job with success",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"components":"search"}`),
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO jobs`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), models.JobScheduled, "search", pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusCreated,
			`"status":"SCHEDULED","components":"search",`,
		},
	}
	executeTests(t, tests)
}

func getMockedJobRow() *pgxmock.Rows {
	return pgxmock.NewRows(jobCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "SCHEDULED", "metadata,places",
			&sampleCoverMediaItemID, sampleTime, sampleTime)
}

func getMockedJobRows() *pgxmock.Rows {
	return pgxmock.NewRows(jobCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "RUNNING", "metadata,places",
			&sampleCoverMediaItemID, sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "RUNNING", "faces",
			&sampleCoverMediaItemID, sampleTime, sampleTime)
}
