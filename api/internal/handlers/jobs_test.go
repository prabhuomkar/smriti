package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
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
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(sqlmock.NewRows(jobCols))
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
			"get job",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
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
		{
			"get job with error",
			http.MethodGet,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
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
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(sqlmock.NewRows(jobCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "jobs"`)).
					WithArgs("RUNNING", sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusNoContent,
			"",
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
			strings.NewReader(`{"status":"PAUSED"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "jobs"`)).
					WithArgs("PAUSED", sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
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
			"update job with error due to job already exists",
			http.MethodPut,
			"/v1/jobs/:id",
			"/v1/jobs/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"status":"RUNNING"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(getMockedJobRow())

			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateJob
			},
			http.StatusConflict,
			"job already exists",
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
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(sqlmock.NewRows(jobCols))
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
			"get jobs with 2 rows",
			http.MethodGet,
			"/v1/jobs",
			"/v1/jobs?sort=updatedAt",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
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
		{
			"get jobs with error",
			http.MethodGet,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
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
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(sqlmock.NewRows(jobCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "jobs"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "SCHEDULED", "search", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusCreated,
			`"status":"SCHEDULED","components":"search",`,
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
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(sqlmock.NewRows(jobCols))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "jobs"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "SCHEDULED", "search", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
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
			"create job with error due to job already exists",
			http.MethodPost,
			"/v1/jobs",
			"/v1/jobs",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"components":"search"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jobs"`)).
					WillReturnRows(getMockedJobRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateJob
			},
			http.StatusConflict,
			"job already exists",
		},
	}
	executeTests(t, tests)
}

func getMockedJobRow() *sqlmock.Rows {
	return sqlmock.NewRows(jobCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "metadata,places", "SCHEDULED",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sampleTime, sampleTime)
}

func getMockedJobRows() *sqlmock.Rows {
	return sqlmock.NewRows(jobCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "metadata,places", "RUNNING",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "faces", "RUNNING",
			"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sampleTime, sampleTime)
}
