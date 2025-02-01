//go:build windows

package main

import (
	"os"
	"syscall"
	"time"
)

// getCreationTime returns the creation time of a file on Windows
func getCreationTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, stat.CreationTime.Nanoseconds())
}

// getAccessTime returns the last access time of a file on Windows
func getAccessTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, stat.LastAccessTime.Nanoseconds())
}

// getModificationTime returns the last modification time of a file on Windows
func getModificationTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, stat.LastWriteTime.Nanoseconds())
}

func getPlatformSpecificAttributes(info os.FileInfo) (isHidden, isSystem, isArchive, isReadOnly bool) {
	if stat, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		isHidden = stat.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
		isSystem = stat.FileAttributes&syscall.FILE_ATTRIBUTE_SYSTEM != 0
		isArchive = stat.FileAttributes&syscall.FILE_ATTRIBUTE_ARCHIVE != 0
		isReadOnly = stat.FileAttributes&syscall.FILE_ATTRIBUTE_READONLY != 0
	}
	return
}