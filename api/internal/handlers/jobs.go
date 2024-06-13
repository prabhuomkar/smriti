package handlers

import (
	"api/internal/models"
	"errors"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type (
	// JobRequest ...
	JobRequest struct {
		Components *string `json:"components"`
		Status     *string `json:"status"`
	}
)

// GetJob ...
func (h *Handler) GetJob(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getJobID(ctx)
	if err != nil {
		return err
	}
	job := models.Job{}
	result := h.DB.Model(&models.Job{}).
		Where("id=? AND user_id=?", uid, userID).
		First(&job)
	if result.Error != nil {
		slog.Error("error getting job", "error", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "job not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
		existingJob := models.Job{}
		result := h.DB.Model(&models.Job{}).
			Where("user_id=? AND status IN (?, ?, ?)", userID, string(models.JobPaused), string(models.JobScheduled), string(models.JobRunning)).
			First(&existingJob)
		if (models.Job{} != existingJob) || (result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound)) {
			slog.Error("job already exists", "error", result.Error)
			return echo.NewHTTPError(http.StatusConflict, "job already exists")
		}
	}
	result := h.DB.Model(&models.Job{UserID: userID, ID: uid}).Updates(map[string]interface{}{
		"Status": job.Status,
	})
	if result.Error != nil {
		slog.Error("error updating job", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetJobs ...
func (h *Handler) GetJobs(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	jobs := []models.Job{}
	result := h.DB.Model(&models.Job{}).
		Where("user_id=?", userID).
		Order("created_at desc").
		Find(&jobs).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		slog.Error("error getting jobs", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	existingJob := models.Job{}
	result := h.DB.Model(&models.Job{}).
		Where("user_id=? AND status IN (?, ?, ?)", userID, string(models.JobPaused), string(models.JobScheduled), string(models.JobRunning)).
		First(&existingJob)
	if (models.Job{} != existingJob) || (result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound)) {
		slog.Error("job already exists", "error", result.Error)
		return echo.NewHTTPError(http.StatusConflict, "job already exists")
	}
	if result := h.DB.Create(&job); result.Error != nil {
		slog.Error("error creating job", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
