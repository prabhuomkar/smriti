package models

import (
	"api/config"
	"syscall"

	"golang.org/x/exp/slog"
)

type (
	// Disk ...
	Disk struct {
		Total uint64 `json:"total,omitempty"`
		Used  uint64 `json:"used,omitempty"`
		Free  uint64 `json:"free,omitempty"`
	}
)

// GetDisk ...
func GetDisk(cfg *config.Config) *Disk {
	diskStat := syscall.Statfs_t{}
	err := syscall.Statfs(cfg.Storage.DiskRoot, &diskStat)
	if err != nil {
		slog.Error("error getting disk stats", slog.Any("error", err))
		return nil
	}
	disk := &Disk{
		Total: diskStat.Blocks * uint64(diskStat.Bsize),
		Free:  diskStat.Bfree * uint64(diskStat.Bsize),
	}
	disk.Used = disk.Total - disk.Free
	return disk
}
