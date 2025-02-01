# fsparq

`fsparq` is a high-performance, cross-platform tool that scans your filesystem and archives detailed metadata into Parquet format. It's designed to efficiently handle large directory structures while maintaining a small memory footprint.

## Features

- 🚀 **High Performance**: Streams data with minimal memory usage
- 🔄 **Cross-Platform**: Full support for Windows, macOS, and Linux
- 🔍 **Rich Metadata**: Captures comprehensive file attributes including:
  - Precise timestamps (creation, modification, access) in UTC
  - File permissions and modes
  - Platform-specific attributes (hidden, system, archive flags)
  - SHA256 hashes for files
- 📦 **Efficient Storage**: Uses Snappy compression and dictionary encoding
- ⚡ **Concurrent Processing**: Leverages goroutines for parallel operations

## Installation

```bash
# Using go install
go install github.com/nrtkbb/fsparq@latest

# Or clone and build
git clone https://github.com/nrtkbb/fsparq.git
cd fsparq
go build
```

## Usage

Basic usage:

```bash
fsparq -root /path/to/scan -output metadata.parquet
```

Advanced options:

```bash
fsparq \
  -root /path/to/scan \
  -output metadata.parquet \
  -buffer 2000 \        # Buffer size for metadata records
  -workers 8 \          # Number of worker goroutines
  -flush 20000          # Records per flush to disk
```

## Output Format

The generated Parquet file includes the following columns:

| Column                | Type    | Description                       |
| --------------------- | ------- | --------------------------------- |
| file_path             | STRING  | Absolute file path                |
| file_name             | STRING  | Base name of the file             |
| directory             | STRING  | Parent directory path             |
| size_bytes            | INT64   | File size in bytes                |
| creation_time_utc     | INT64   | File creation timestamp (UTC)     |
| modification_time_utc | INT64   | Last modification timestamp (UTC) |
| access_time_utc       | INT64   | Last access timestamp (UTC)       |
| file_mode             | STRING  | File permissions (Unix-style)     |
| is_directory          | BOOLEAN | Directory flag                    |
| is_file               | BOOLEAN | Regular file flag                 |
| is_symlink            | BOOLEAN | Symbolic link flag                |
| is_hidden             | BOOLEAN | Hidden file flag                  |
| is_system             | BOOLEAN | System file flag (Windows)        |
| is_archive            | BOOLEAN | Archive flag (Windows)            |
| is_readonly           | BOOLEAN | Read-only flag                    |
| file_extension        | STRING  | File extension (with dot)         |
| sha256                | STRING  | SHA256 hash (files only)          |

## Platform-Specific Behavior

### Windows

- Uses native Win32 API for accurate file attributes
- Supports NTFS timestamps and special flags (hidden, system, archive)
- File paths use backslash separator

### macOS

- Uses native birth time when available
- Hidden files determined by dot prefix
- File paths use forward slash separator

### Linux

- Uses ctime as fallback when birth time unavailable
- System and archive attributes always false
- File paths use forward slash separator

## Performance Considerations

- Streams data to minimize memory usage
- Uses buffered writes with configurable flush intervals
- Employs Snappy compression for efficient storage
- Utilizes dictionary encoding for repeated strings
- Processes files concurrently with configurable worker count

## Error Handling

- Continues processing on permission errors
- Logs warnings for inaccessible files
- Skips hash calculation for unreadable files
- Maintains a record of processing errors

## Build

Build for different platforms:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o fsparq.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o fsparq-mac

# Linux
GOOS=linux GOARCH=amd64 go build -o fsparq-linux
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
