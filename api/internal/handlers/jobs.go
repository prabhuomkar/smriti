package handlers

import (
	"api/internal/models"
	"api/pkg/services/api"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type ( // JobRequest ...
	JobRequest struct {
		Components []string `json:"components"`
		Status     *string  `json:"status"`
	}
)

const (
	queryGetJob         = `SELECT * FROM jobs WHERE user_id=$1 AND id=$2`
	queryGetJobs        = `SELECT * FROM jobs WHERE user_id=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`
	queryCheckJobExists = `SELECT COUNT(*) FROM jobs WHERE user_id=$1 AND status IN ($2, $3)`
	queryCreateJob      = `INSERT INTO jobs (id, user_id, status, components, created_at, updated_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6)`
	queryUpdateJob             = `UPDATE jobs SET status=$3, updated_at=$4 WHERE user_id=$1 AND id=$2`
	queryInsertQueueMediaItems = `INSERT INTO queue (id, user_id, mediaitem_id, job_id, components, status) ` +
		`SELECT gen_random_uuid(), user_id, id, $2, $3, $4 FROM mediaitems WHERE user_id=$1`
	queryDeleteQueueMediaItems = `DELETE FROM queue WHERE user_id=$1 AND job_id=$2`
)

// GetJob ...
func (h *Handler) GetJob(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getJobID(ctx)
	if err != nil {
		return err
	}
	job := models.Job{}
	jobComponents := ""
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetJob, userID, uid).Scan(&job.ID, &job.UserID,
		&job.Status, &jobComponents, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "job not found")
		}
		slog.Error("error getting job", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	job.Components = strings.Split(jobComponents, ",")

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
	if job.Status == models.JobRunning { //nolint: staticcheck
		existingJobCount := 0
		err = h.DB.QueryRow(ctx.Request().Context(), queryCheckJobExists, userID, string(models.JobPaused),
			string(models.JobRunning)).
			Scan(&existingJobCount)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("error getting existing job count", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if existingJobCount > 0 {
			slog.Error("job already exists", "error", err)

			return echo.NewHTTPError(http.StatusConflict, "job already exists")
		}
	} else if job.Status == models.JobStopped {
		_, err = h.DB.Exec(ctx.Request().Context(), queryDeleteQueueMediaItems, userID, uid)
		if err != nil {
			slog.Error("error deleting job queue mediaitems", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	_, err = h.DB.Exec(ctx.Request().Context(), queryUpdateJob, userID, uid, job.Status, job.UpdatedAt)
	if err != nil {
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
		jobComponents := ""
		err = rows.Scan(&job.ID, &job.UserID, &job.Status, &jobComponents, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			slog.Error("error scanning job", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		job.Components = strings.Split(jobComponents, ",")
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
	job.ID, _ = uuid.NewV7()
	job.UserID = userID
	job.Status = models.JobRunning
	job.CreatedAt = time.Now()
	job.UpdatedAt = job.CreatedAt
	existingJobCount := 0
	err = h.DB.QueryRow(ctx.Request().Context(), queryCheckJobExists, userID,
		string(models.JobPaused), string(models.JobRunning)).Scan(&existingJobCount)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting existing job count", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if existingJobCount > 0 {
		slog.Error("job already exists", "error", err)

		return echo.NewHTTPError(http.StatusConflict, "job already exists")
	}
	_, err = h.DB.Exec(ctx.Request().Context(), queryCreateJob, job.ID, job.UserID, job.Status,
		strings.Join(job.Components, ","), job.CreatedAt, job.UpdatedAt)
	if err != nil {
		slog.Error("error creating job", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = h.DB.Exec(ctx.Request().Context(), queryInsertQueueMediaItems, userID, job.ID,
		strings.Join(job.Components, ","), api.MediaItemStatus_UNSPECIFIED)
	if err != nil {
		slog.Error("error inserting job queue mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, job)
}

func getJobID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
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
	if len(jobRequest.Components) > 0 {
		for _, jobComponent := range jobRequest.Components {
			if _, ok := api.MediaItemComponent_value[jobComponent]; !ok {
				slog.Error("error getting job component", "error", err)

				return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid job component")
			}
		}
		job.Components = jobRequest.Components
	}
	if jobRequest.Status != nil {
		job.Status = models.JobStatus(*jobRequest.Status)
	}
	if reflect.DeepEqual(models.Job{}, job) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid job")
	}

	return &job, nil
}
