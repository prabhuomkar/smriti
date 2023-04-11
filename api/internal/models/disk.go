package models

import (
	"api/config"
	"log"
	"syscall"
)

type (
	// Disk ...
	Disk struct {
		All  uint64 `json:"all,omitempty"`
		Used uint64 `json:"used,omitempty"`
		Free uint64 `json:"free,omitempty"`
	}
)

// GetDisk ...
func GetDisk(cfg *config.Config) *Disk {
	// nolint: nosnakecase
	diskStat := syscall.Statfs_t{}
	err := syscall.Statfs(cfg.StorageDiskRoot, &diskStat)
	if err != nil {
		log.Printf("error getting disk stats: %+v", err)
		return nil
	}
	disk := &Disk{
		All:  diskStat.Blocks * uint64(diskStat.Bsize),
		Free: diskStat.Bfree * uint64(diskStat.Bsize),
	}
	disk.Used = disk.All - disk.Free
	return disk
}
