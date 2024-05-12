package jobs

import (
	"api/config"
	"api/internal/models"
	"api/pkg/storage"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// Job ...
type Job struct {
	Config  *config.Config
	DB      *gorm.DB
	Storage storage.Provider
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
	result := j.DB.Model(&models.Job{}).
		Where("status=?", string(models.JobScheduled)).
		Find(jobs)
	if result.Error != nil {
		slog.Error("error getting jobs", "error", result.Error)
		return
	}
	for _, job := range jobs {
		go j.executeJob(job)
	}
}

func (j *Job) executeJob(jobCfg models.Job) {
	// get status of job if should proceed or not, exit if paused or stopped
	// loop and parallely get mediaitems to send it to worker
	// update the last mediaitem processed
}
