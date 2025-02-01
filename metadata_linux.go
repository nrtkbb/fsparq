//go:build linux

package main

import (
	"os"
	"syscall"
	"time"
)

// getCreationTime returns the creation time of a file on Linux
func getCreationTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	// Linux doesn't always provide birth time (creation time),
	// so we fall back to change time (ctime)
	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
}

// getAccessTime returns the last access time of a file on Linux
func getAccessTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
}

// getModificationTime returns the last modification time of a file on Linux
func getModificationTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
}

func getPlatformSpecificAttributes(info os.FileInfo) (isHidden, isSystem, isArchive, isReadOnly bool) {
	// Hidden files in Unix start with a dot
	isHidden = len(info.Name()) > 0 && info.Name()[0] == '.'
	
	// Read-only if no write permission
	isReadOnly = info.Mode().Perm()&0200 == 0
	
	// System and Archive attributes are not supported on Unix systems
	isSystem = false
	isArchive = false
	
	return
}