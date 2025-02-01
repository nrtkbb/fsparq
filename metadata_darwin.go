//go:build darwin

package main

import (
	"os"
	"syscall"
	"time"
)

// getCreationTime returns the creation time of a file on macOS
func getCreationTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	// macOS provides birthtime (creation time)
	return time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
}

// getAccessTime returns the last access time of a file on macOS
func getAccessTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
}

// getModificationTime returns the last modification time of a file on macOS
func getModificationTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec)
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