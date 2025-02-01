// file_attributes.go
package main

import "os"

// getFileAttributes retrieves platform-specific file attributes
func getFileAttributes(info os.FileInfo) (isHidden, isSystem, isArchive, isReadOnly bool) {
	return getPlatformSpecificAttributes(info)
}