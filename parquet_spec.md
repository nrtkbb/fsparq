# File Metadata Parquet Specification

## 1. File Overview

This Parquet file stores filesystem metadata collected from Windows, macOS, and Linux file systems.

### 1.1 Basic Information

- File Format: Apache Parquet
- Compression: Snappy
- Row Group Size: 128MB
- Encoding: PLAIN_DICTIONARY (default) with type-specific optimizations

## 2. Schema Definition

### 2.1 Column List

| Column Name           | Data Type  | Converted Type | Encoding         | Description                  |
| --------------------- | ---------- | -------------- | ---------------- | ---------------------------- |
| file_path             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | Absolute file path           |
| file_name             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | File name                    |
| directory             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | Parent directory path        |
| size_bytes            | INT64      | -              | PLAIN            | File size in bytes           |
| creation_time_utc     | INT64      | -              | PLAIN            | Creation timestamp (UTC)     |
| modification_time_utc | INT64      | -              | PLAIN            | Modification timestamp (UTC) |
| access_time_utc       | INT64      | -              | PLAIN            | Last access timestamp (UTC)  |
| file_mode             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | File mode (permissions)      |
| is_directory          | BOOLEAN    | -              | PLAIN            | Directory flag               |
| is_file               | BOOLEAN    | -              | PLAIN            | Regular file flag            |
| is_symlink            | BOOLEAN    | -              | PLAIN            | Symbolic link flag           |
| is_hidden             | BOOLEAN    | -              | PLAIN            | Hidden file flag             |
| is_system             | BOOLEAN    | -              | PLAIN            | System file flag             |
| is_archive            | BOOLEAN    | -              | PLAIN            | Archive flag                 |
| is_readonly           | BOOLEAN    | -              | PLAIN            | Read-only flag               |
| file_extension        | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | File extension               |
| sha256                | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | SHA256 hash of file          |

### 2.2 Detailed Data Type Specifications

#### String Fields (BYTE_ARRAY, UTF8)

- `file_path`

  - Absolute filesystem path
  - Platform-native path separators:
    - Windows: backslash (\)
    - Unix-like: forward slash (/)
  - Maximum length: Platform-dependent

- `file_name`

  - File name without path separators
  - Maximum length: Platform-dependent
    - Windows: 255 characters
    - Unix-like: typically 255 bytes

- `directory`

  - Parent directory portion of `file_path`
  - No trailing path separator

- `file_mode`

  - Unix-style permission string
  - Fixed length: 10 characters
  - Example: "drwxr-xr-x"

- `file_extension`

  - File extension including dot
  - Empty string if no extension
  - Case sensitivity is platform-dependent

- `sha256`
  - Hexadecimal string representation
  - Fixed length: 64 characters
  - Null for directories

#### Numeric Fields (INT64)

- `size_bytes`

  - Range: 0 to 2^63-1
  - Directory size is platform-dependent

- Timestamp Fields (all UTC)
  - Unix timestamp (seconds)
  - Range: 0 to 2^63-1
  - Platform-specific retrieval methods:
    - Windows: Win32FileAttributeData
    - macOS: Birthtimespec/Atimespec/Mtimespec
    - Linux: Ctim/Atim/Mtim

#### Boolean Fields (BOOLEAN)

- `is_directory`

  - True if entry is a directory

- `is_file`

  - True if entry is a regular file

- `is_symlink`

  - True if entry is a symbolic link

- `is_hidden`

  - Windows: True if FILE_ATTRIBUTE_HIDDEN is set
  - Unix-like: True if filename begins with dot (.)

- `is_system`

  - Windows: True if FILE_ATTRIBUTE_SYSTEM is set
  - Unix-like: Always false

- `is_archive`

  - Windows: True if FILE_ATTRIBUTE_ARCHIVE is set
  - Unix-like: Always false

- `is_readonly`
  - Windows: True if FILE_ATTRIBUTE_READONLY is set
  - Unix-like: True if write permission is not set

## 3. Platform-Specific Behavior

### 3.1 Windows-Specific

- Uses backslash as path separator
- Full support for filesystem attributes (hidden, system, archive)
- NTFS timestamps (100-nanosecond precision) converted to UTC

### 3.2 macOS-Specific

- Uses forward slash as path separator
- Accurate birth time (creation time) from HFS+/APFS
- Hidden files determined by dot prefix

### 3.3 Linux-Specific

- Uses forward slash as path separator
- Falls back to change time (ctime) for creation time
- System and archive attributes unsupported (always false)

## 4. Error Handling and Edge Cases

### 4.1 NULL Value Handling

- `sha256`: Null for directories or hash calculation errors
- `file_extension`: Empty string for no extension (not null)
- Timestamps: 0 if unsupported on platform
- Other fields: Always non-null

### 4.2 Error Behavior

- Permission errors: Skip affected file
- Hash calculation errors: Set SHA256 field to null
- Broken symbolic links: Store metadata of the link itself

## 5. Performance Characteristics

### 5.1 Memory Usage

- Buffer size: 1000 records
- Row Group size: 128MB
- String deduplication via PLAIN_DICTIONARY encoding

### 5.2 Disk Usage

- Efficient storage with Snappy compression
- Dictionary compression for strings
- Statistics retention (Min/Max, Null count)

### 5.3 Processing Optimizations

- Streaming processing for memory efficiency
- Concurrent file processing
- Configurable buffer sizes and flush intervals
- Optimized platform-specific system calls
