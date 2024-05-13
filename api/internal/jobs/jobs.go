package jobs

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/worker"
	"api/pkg/storage"
	"log/slog"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Job ...
type Job struct {
	Config  *config.Config
	DB      *gorm.DB
	Storage storage.Provider
	Worker  worker.WorkerClient
}

func (j *Job) StartJobs() {
	ticker := time.NewTicker(j.Config.Job.QueueInterval)
	go func() {
		for {
			<-ticker.C
			j.queueJob()
		}
	}()
}

func (j *Job) queueJob() {
	jobs := []models.Job{}
	result := j.DB.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate, Options: clause.LockingOptionsSkipLocked}).Model(&models.Job{}).
		Where("status=?", string(models.JobScheduled)).
		Find(&jobs)
	if result.Error != nil {
		slog.Error("error getting jobs", "error", result.Error)
		return
	}
	for _, job := range jobs {
		go j.executeJob(job)
	}
}

func (j *Job) executeJob(jobCfg models.Job) {
	slog.Info("starting job", "userId", jobCfg.UserID, "job", jobCfg.ID)
	result := j.DB.Model(&models.Job{UserID: jobCfg.UserID, ID: jobCfg.ID}).Updates(map[string]interface{}{
		"Status": string(models.JobRunning),
	})
	if result.Error != nil {
		slog.Error("error updating job status", "error", result.Error)
		return
	}
	var executorWg sync.WaitGroup
	queue := make(chan models.MediaItem, j.Config.Job.Concurrency)
	results := make(chan uuid.UUID)
	for range j.Config.Job.Concurrency {
		executorWg.Add(1)
		go j.executeJobMediaItem(&executorWg, queue, results)
	}
	mediaItem, err := j.getJobMediaItem(jobCfg, uuid.Nil)
	if err != nil {
		j.updateJobStatus(jobCfg, models.JobCompleted)
	} else {
		queue <- mediaItem
		for result := range results {
			slog.Info("completed item from job queue", "mediaItem", result)
			err := j.updateJobLastMediaItem(jobCfg, result)
			if err != nil {
				j.updateJobStatus(jobCfg, models.JobCompleted)
				break
			}
			jobStatus := j.getJobStatus(jobCfg)
			if jobStatus == models.JobRunning {
				mediaItem, err := j.getJobMediaItem(jobCfg, result)
				if err != nil {
					j.updateJobStatus(jobCfg, models.JobCompleted)
					break
				}
				queue <- mediaItem
			} else {
				slog.Info("stopping job", "userId", jobCfg.UserID, "job", jobCfg.ID, "status", jobStatus)
				break
			}
		}
	}
	close(queue)
	slog.Info("waiting for job to complete", "userId", jobCfg.UserID, "job", jobCfg.ID)
	executorWg.Wait()
	close(results)
	slog.Info("completed job", "userId", jobCfg.UserID, "job", jobCfg.ID)
}

func (j *Job) getJobMediaItem(jobCfg models.Job, lastMediaItemID uuid.UUID) (models.MediaItem, error) {
	mediaItem := models.MediaItem{}
	result := j.DB.Where("user_id=? AND id>?", jobCfg.UserID, lastMediaItemID).Order("created_at").First(&mediaItem)
	if result.Error != nil {
		slog.Error("error getting job mediaitem", "error", result.Error)
		return mediaItem, result.Error
	}
	return mediaItem, nil
}

func (j *Job) updateJobLastMediaItem(jobCfg models.Job, mediaItemID uuid.UUID) error {
	result := j.DB.Model(&models.Job{UserID: jobCfg.UserID, ID: jobCfg.ID}).Updates(map[string]interface{}{
		"LastMediaItemID": mediaItemID,
	})
	if result.Error != nil {
		slog.Error("error updating job last mediaitem", "error", result.Error)
		return result.Error
	}
	return nil
}

func (j *Job) executeJobMediaItem(wg *sync.WaitGroup, queue <-chan models.MediaItem, results chan<- uuid.UUID) {
	defer wg.Done()
	for item := range queue {
		slog.Debug("processing item from job queue", "mediaItem", item.ID.String())
		// send mediaitem to worker for processing
		results <- item.ID
	}
}

func (j *Job) getJobStatus(jobCfg models.Job) models.JobStatus {
	currentJob := models.Job{}
	result := j.DB.Model(&models.Job{}).
		Where("id=? AND user_id=?", jobCfg.ID, jobCfg.UserID).
		First(&currentJob)
	if result.Error != nil {
		slog.Error("error getting job status", "error", result.Error)
		return ""
	}
	return currentJob.Status
}

func (j *Job) updateJobStatus(jobCfg models.Job, status models.JobStatus) {
	result := j.DB.Model(&models.Job{UserID: jobCfg.UserID, ID: jobCfg.ID}).Updates(map[string]interface{}{
		"Status": status,
	})
	if result.Error != nil {
		slog.Error("error updating job status", "error", result.Error)
	}
}
