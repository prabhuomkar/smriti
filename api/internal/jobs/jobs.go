package jobs

import (
	"api/config"
	"api/internal/models"
	"api/pkg/services/worker"
	"api/pkg/storage"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
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
	init := true
	ticker := time.NewTicker(j.Config.Job.QueueInterval)
	go func() {
		for {
			<-ticker.C
			j.queueJob(init)
			init = false
		}
	}()
}

func (j *Job) queueJob(init bool) {
	jobs := []models.Job{}
	filter := "status='" + string(models.JobScheduled) + "'"
	if init {
		filter += " OR status='" + string(models.JobRunning) + "'"
	}
	result := j.DB.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate, Options: clause.LockingOptionsSkipLocked}).Model(&models.Job{}).
		Where(filter).
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
		go j.executeJobMediaItem(&executorWg, jobCfg, queue, results)
	}
	mediaItem, err := j.getJobMediaItem(jobCfg, uuid.Nil)
	if err != nil { //nolint: nestif
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
		"last_mediaitem_id": mediaItemID,
	})
	if result.Error != nil {
		slog.Error("error updating job last mediaitem", "error", result.Error)
		return result.Error
	}
	return nil
}

func (j *Job) executeJobMediaItem(wg *sync.WaitGroup, jobCfg models.Job, queue <-chan models.MediaItem, results chan<- uuid.UUID) { //nolint: cyclop
	defer wg.Done()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("error creating file watcher", "error", err)
	}
	defer watcher.Close()
	for item := range queue {
		slog.Debug("processing item from job queue", "mediaItem", item.ID.String())
		// download the mediaitem to disk root depending on the type of the job
		fileType, fileName := getFileTypeAndFileName(jobCfg.Components, item.ID.String(), string(item.MediaItemType))
		filePath := fmt.Sprintf("%s/%s", j.Config.Storage.DiskRoot, fileName)
		err := j.Storage.Download(filePath, fileType, item.ID.String())
		if err != nil {
			slog.Error("error downloading mediaitem for processing", "error", err)
			continue
		}

		// send to worker for processing
		_, err = j.Worker.MediaItemProcess(context.Background(), &worker.MediaItemProcessRequest{
			UserId:     item.UserID.String(),
			Id:         item.ID.String(),
			FilePath:   j.Config.Storage.DiskRoot,
			Components: strings.Split(jobCfg.Components, ","),
		})
		if err != nil {
			slog.Error("error sending mediaitem for processing", "error", err)
			continue
		}

		// start a file watcher to notify when the file is removed
		err = watcher.Add(filePath)
		if err != nil {
			slog.Error("error adding file to watcher", "file", fileName, "error", err)
		}

	watcherLoop:
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					break watcherLoop
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					slog.Debug("finished processing item from job queue", "mediaItem", item.ID.String())
					results <- item.ID
					break watcherLoop
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					break watcherLoop
				}
				slog.Error("error watching file", "file", fileName, "error", err)
			}
		}
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

func getFileTypeAndFileName(components, mediaItemID, mediaItemType string) (string, string) {
	if strings.Contains(components, "metadata") {
		return "originals", mediaItemID
	}
	suffix := "-preview"
	if mediaItemType == "video" {
		suffix += ".mp4"
	}
	return "previews", mediaItemID + suffix
}
