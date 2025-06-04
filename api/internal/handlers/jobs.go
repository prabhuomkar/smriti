package handlers

import (
	"api/internal/models"
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type (
	// JobRequest ...
	JobRequest struct {
		Components *string `json:"components"`
		Status     *string `json:"status"`
	}
)

const (
	queryGetJob         = `SELECT * FROM jobs WHERE user_id=$1 AND id=$2`
	queryGetJobs        = `SELECT * FROM jobs WHERE user_id=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`
	queryCheckJobExists = `SELECT COUNT(*) FROM jobs WHERE user_id=$1 AND status IN ($2, $3, $4)`
	queryCreateJob      = `INSERT INTO jobs (id, user_id, status, components, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	queryUpdateJob      = `UPDATE jobs SET status=$3, updated_at=$4 WHERE user_id=$1 AND id=$2`
)

// GetJob ...
func (h *Handler) GetJob(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getJobID(ctx)
	if err != nil {
		return err
	}
	job := models.Job{}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetJob, userID, uid).Scan(
		&job.ID,
		&job.UserID,
		&job.Status,
		&job.Components,
		&job.LastMediItemID,
		&job.CreatedAt,
		&job.UpdatedAt)
	if err != nil {
		slog.Error("error getting job", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "job not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, job)
}

// UpdateJob ...
func (h *Handler) UpdateJob(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getJobID(ctx)
	if err != nil {
		return err
	}
	job, err := getJob(ctx)
	if err != nil {
		return err
	}
	if job.Status == models.JobRunning {
		existingJobCount := 0
		err = h.DB.QueryRow(ctx.Request().Context(), queryCheckJobExists, userID, string(models.JobPaused), string(models.JobScheduled), string(models.JobRunning)).Scan(&existingJobCount)
		if err != nil {
			slog.Error("error getting existing job count", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if existingJobCount > 0 {
			slog.Error("job already exists", "error", err)
			return echo.NewHTTPError(http.StatusConflict, "job already exists")
		}
	}
	result, err := h.DB.Exec(ctx.Request().Context(), queryUpdateJob, userID, uid, job.Status, job.UpdatedAt)
	if !result.Update() || err != nil {
		slog.Error("error updating job", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetJobs ...
func (h *Handler) GetJobs(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	jobs := []models.Job{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetJobs, userID, offset, limit)
	if err != nil {
		slog.Error("error getting jobs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		job := models.Job{}
		err = rows.Scan(&job.ID,
			&job.UserID,
			&job.Status,
			&job.Components,
			&job.LastMediItemID,
			&job.CreatedAt,
			&job.UpdatedAt)
		if err != nil {
			slog.Error("error scanning job", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		jobs = append(jobs, job)
	}
	return ctx.JSON(http.StatusOK, jobs)
}

// CreateJob ...
func (h *Handler) CreateJob(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	job, err := getJob(ctx)
	if err != nil {
		return err
	}
	job.ID = uuid.NewV4()
	job.UserID = userID
	job.Status = models.JobScheduled
	existingJobCount := 0
	err = h.DB.QueryRow(ctx.Request().Context(), queryCheckJobExists, userID, string(models.JobPaused), string(models.JobScheduled), string(models.JobRunning)).Scan(&existingJobCount)
	if err != nil {
		slog.Error("error getting existing job count", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if existingJobCount > 0 {
		slog.Error("job already exists", "error", err)
		return echo.NewHTTPError(http.StatusConflict, "job already exists")
	}
	result, err := h.DB.Exec(ctx.Request().Context(), queryCreateJob,
		job.ID,
		job.UserID,
		job.Status,
		job.Components,
		job.CreatedAt,
		job.UpdatedAt)
	if !result.Insert() || err != nil {
		slog.Error("error creating job", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, job)
}

func getJobID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting job id", "error", err)
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid job id")
	}
	return uid, err
}

func getJob(ctx echo.Context) (*models.Job, error) {
	jobRequest := new(JobRequest)
	err := ctx.Bind(jobRequest)
	if err != nil {
		slog.Error("error getting job", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid job")
	}
	job := models.Job{}
	if jobRequest.Components != nil {
		job.Components = *jobRequest.Components
	}
	if jobRequest.Status != nil {
		job.Status = models.JobStatus(*jobRequest.Status)
	}
	if reflect.DeepEqual(models.Job{}, job) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid job")
	}
	return &job, nil
}
