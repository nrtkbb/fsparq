package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

// FileMetadata represents file metadata to be stored in Parquet format
type FileMetadata struct {
	FilePath            string `parquet:"name=file_path, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	FileName            string `parquet:"name=file_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Directory           string `parquet:"name=directory, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	SizeBytes           int64  `parquet:"name=size_bytes, type=INT64"`
	CreationTimeUTC     int64  `parquet:"name=creation_time_utc, type=INT64"`
	ModificationTimeUTC int64  `parquet:"name=modification_time_utc, type=INT64"`
	AccessTimeUTC       int64  `parquet:"name=access_time_utc, type=INT64"`
	FileMode            string `parquet:"name=file_mode, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	IsDirectory         bool   `parquet:"name=is_directory, type=BOOLEAN"`
	IsFile              bool   `parquet:"name=is_file, type=BOOLEAN"`
	IsSymlink           bool   `parquet:"name=is_symlink, type=BOOLEAN"`
	IsHidden            bool   `parquet:"name=is_hidden, type=BOOLEAN"`
	IsSystem            bool   `parquet:"name=is_system, type=BOOLEAN"`
	IsArchive           bool   `parquet:"name=is_archive, type=BOOLEAN"`
	IsReadOnly          bool   `parquet:"name=is_readonly, type=BOOLEAN"`
	FileExtension       string `parquet:"name=file_extension, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	SHA256              string `parquet:"name=sha256, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ParquetFileName     string `parquet:"name=parquet_file_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

// Config holds application configuration
type Config struct {
	rootDir    string
	outputFile string
	bufferSize int
	workers    int
	flushSize  int
}

// calculateSHA256 computes the SHA256 hash of a file
func calculateSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// timeToUTC converts a timestamp to UTC
func timeToUTC(t time.Time) int64 {
	return t.UTC().Unix()
}

// walkFiles traverses the filesystem and collects metadata
func walkFiles(config Config) <-chan FileMetadata {
	metadataChan := make(chan FileMetadata, config.bufferSize)

	go func() {
		defer close(metadataChan)

		err := filepath.Walk(config.rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Log error and continue processing
				log.Printf("Warning: Cannot access %s: %v", path, err)
				return filepath.SkipDir
			}

			// Calculate hash only for files
			var hash string
			if !info.IsDir() {
				h, err := calculateSHA256(path)
				if err != nil {
					log.Printf("Warning: Failed to calculate hash for %s: %v", path, err)
				} else {
					hash = h
				}
			}

			// Get platform-specific timestamps
			ctime := getCreationTime(info)
			mtime := getModificationTime(info)
			atime := getAccessTime(info)

			// Get platform-specific file attributes
			isHidden, isSystem, isArchive, isReadOnly := getFileAttributes(info)

			metadata := FileMetadata{
				FilePath:            path,
				FileName:            info.Name(),
				Directory:           filepath.Dir(path),
				SizeBytes:           info.Size(),
				CreationTimeUTC:     timeToUTC(ctime),
				ModificationTimeUTC: timeToUTC(mtime),
				AccessTimeUTC:       timeToUTC(atime),
				FileMode:            info.Mode().String(),
				IsDirectory:         info.IsDir(),
				IsFile:              info.Mode().IsRegular(),
				IsSymlink:           info.Mode()&os.ModeSymlink != 0,
				IsHidden:            isHidden,
				IsSystem:            isSystem,
				IsArchive:           isArchive,
				IsReadOnly:          isReadOnly,
				FileExtension:       filepath.Ext(path),
				SHA256:              hash,
				ParquetFileName:     config.outputFile,
			}

			metadataChan <- metadata
			return nil
		})

		if err != nil && err != filepath.SkipDir {
			log.Printf("Error during walk: %v", err)
		}
	}()

	return metadataChan
}

// writeParquet writes metadata to a Parquet file
func writeParquet(config Config, metadataChan <-chan FileMetadata) error {
	fw, err := local.NewLocalFileWriter(config.outputFile)
	if err != nil {
		return fmt.Errorf("failed to create parquet file: %v", err)
	}
	defer fw.Close()

	pw, err := writer.NewParquetWriter(fw, new(FileMetadata), int64(config.workers))
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %v", err)
	}

	// Set performance settings
	pw.RowGroupSize = 128 * 1024 * 1024 // 128MB
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	count := 0
	for metadata := range metadataChan {
		if err := pw.Write(metadata); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}

		count++
		// Flush writer every set number of records
		if count%config.flushSize == 0 {
			if err := pw.Flush(true); err != nil {
				return fmt.Errorf("failed to flush writer: %v", err)
			}
			log.Printf("Processed %d files", count)
		}
	}

	if err := pw.WriteStop(); err != nil {
		return fmt.Errorf("failed to finish writing: %v", err)
	}

	log.Printf("Successfully processed %d files", count)
	return nil
}

func main() {
	config := Config{}

	flag.StringVar(&config.rootDir, "root", "", "Root directory to process")
	flag.StringVar(&config.outputFile, "output", "file_metadata.parquet", "Output Parquet file path")
	flag.IntVar(&config.bufferSize, "buffer", 1000, "Size of the metadata buffer")
	flag.IntVar(&config.workers, "workers", 4, "Number of worker goroutines")
	flag.IntVar(&config.flushSize, "flush", 10000, "Number of records before flushing to disk")
	flag.Parse()

	if config.rootDir == "" {
		log.Fatal("Root directory is required")
	}

	// Execute metadata collection and Parquet writing pipeline
	metadataChan := walkFiles(config)
	if err := writeParquet(config, metadataChan); err != nil {
		log.Fatal(err)
	}
}
