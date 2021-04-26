package sysMonitor

import (
	"fmt"
	"syscall"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// 获取磁盘信息
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	var toG = func(val uint64) float32 {
		return float32(val / 1024 / 1024 / 1024)
	}
	fmt.Printf("总空间=%fG;使用空间=%fG;剩余空间=%fG", toG(disk.All), toG(disk.Used), toG(disk.Free))
	return
}
