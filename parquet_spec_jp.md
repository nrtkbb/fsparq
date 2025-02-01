# ファイルメタデータ Parquet 仕様書

## 1. ファイル概要

この Parquet ファイルは、Windows、macOS、Linux のファイルシステムから収集したファイルおよびディレクトリのメタデータを保存します。

### 1.1 基本情報

- ファイル形式: Apache Parquet
- 圧縮方式: Snappy
- Row Group サイズ: 128MB
- エンコーディング: PLAIN_DICTIONARY を基本とし、データ型に応じて最適化

## 2. スキーマ定義

### 2.1 カラム一覧

| カラム名              | データ型   | Converted Type | エンコーディング | 説明                                           |
| --------------------- | ---------- | -------------- | ---------------- | ---------------------------------------------- |
| file_path             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | ファイルの絶対パス                             |
| file_name             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | ファイル名                                     |
| directory             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | 親ディレクトリのパス                           |
| size_bytes            | INT64      | -              | PLAIN            | ファイルサイズ（バイト）                       |
| creation_time_utc     | INT64      | -              | PLAIN            | 作成日時（UTC の Unix タイムスタンプ）         |
| modification_time_utc | INT64      | -              | PLAIN            | 更新日時（UTC の Unix タイムスタンプ）         |
| access_time_utc       | INT64      | -              | PLAIN            | 最終アクセス日時（UTC の Unix タイムスタンプ） |
| file_mode             | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | ファイルモード（権限）                         |
| is_directory          | BOOLEAN    | -              | PLAIN            | ディレクトリフラグ                             |
| is_file               | BOOLEAN    | -              | PLAIN            | 通常ファイルフラグ                             |
| is_symlink            | BOOLEAN    | -              | PLAIN            | シンボリックリンクフラグ                       |
| is_hidden             | BOOLEAN    | -              | PLAIN            | 隠しファイルフラグ                             |
| is_system             | BOOLEAN    | -              | PLAIN            | システムファイルフラグ                         |
| is_archive            | BOOLEAN    | -              | PLAIN            | アーカイブフラグ                               |
| is_readonly           | BOOLEAN    | -              | PLAIN            | 読み取り専用フラグ                             |
| file_extension        | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | ファイル拡張子                                 |
| sha256                | BYTE_ARRAY | UTF8           | PLAIN_DICTIONARY | ファイルの SHA256 ハッシュ値                   |

### 2.2 データ型の詳細仕様

#### 文字列フィールド（BYTE_ARRAY, UTF8）

- `file_path`

  - ファイルシステムの絶対パス
  - プラットフォームネイティブのパス区切り文字を保持
    - Windows: バックスラッシュ（\）
    - Unix 系: スラッシュ（/）
  - 最大長: プラットフォーム依存

- `file_name`

  - パス区切り文字を含まないファイル名
  - 最大長: プラットフォーム依存
    - Windows: 255 文字
    - Unix 系: 通常 255 バイト

- `directory`

  - `file_path`の親ディレクトリ部分
  - 末尾のパス区切り文字は含まない

- `file_mode`

  - Unix スタイルのパーミッション文字列
  - 固定長: 10 文字
  - 例: "drwxr-xr-x"

- `file_extension`

  - ピリオドを含む拡張子
  - 拡張子がない場合は空文字列
  - 大文字小文字の区別はプラットフォーム依存

- `sha256`
  - 16 進数の文字列表現
  - 固定長: 64 文字
  - ディレクトリの場合は null

#### 数値フィールド（INT64）

- `size_bytes`

  - 範囲: 0 ～ 2^63-1
  - ディレクトリサイズはプラットフォーム依存

- タイムスタンプフィールド（全て UTC）
  - Unix 時間（秒）
  - 範囲: 0 ～ 2^63-1
  - プラットフォーム別の取得方法：
    - Windows: Win32FileAttributeData
    - macOS: Birthtimespec/Atimespec/Mtimespec
    - Linux: Ctim/Atim/Mtim

#### 真偽値フィールド（BOOLEAN）

- `is_directory`

  - ディレクトリの場合 true

- `is_file`

  - 通常ファイルの場合 true

- `is_symlink`

  - シンボリックリンクの場合 true

- `is_hidden`

  - Windows: FILE_ATTRIBUTE_HIDDEN が設定されている場合 true
  - Unix 系: ファイル名がドット（.）で始まる場合 true

- `is_system`

  - Windows: FILE_ATTRIBUTE_SYSTEM が設定されている場合 true
  - Unix 系: 常に false

- `is_archive`

  - Windows: FILE_ATTRIBUTE_ARCHIVE が設定されている場合 true
  - Unix 系: 常に false

- `is_readonly`
  - Windows: FILE_ATTRIBUTE_READONLY が設定されている場合 true
  - Unix 系: 書き込み権限がない場合 true

## 3. プラットフォーム固有の動作

### 3.1 Windows 固有の動作

- ファイルパスの区切り文字にバックスラッシュを使用
- ファイルシステム属性（hidden, system, archive）を完全サポート
- NTFS タイムスタンプ（100 ナノ秒精度）を UTC に変換

### 3.2 macOS 固有の動作

- ファイルパスの区切り文字にスラッシュを使用
- HFS+/APFS の作成時刻（birthtime）を正確に取得
- 隠しファイルはドットファイルとして判定

### 3.3 Linux 固有の動作

- ファイルパスの区切り文字にスラッシュを使用
- 作成時刻は change time（ctime）をフォールバックとして使用
- システム属性とアーカイブ属性は未サポート（常に false）

## 4. エラー処理とエッジケース

### 4.1 NULL 値の取り扱い

- `sha256`: ディレクトリまたはハッシュ計算エラー時は null
- `file_extension`: 拡張子がない場合は空文字列
- タイムスタンプ: プラットフォームでサポートされていない場合、0 を使用
- その他のフィールド: 常に non-null

### 4.2 エラー時の動作

- アクセス権限エラー: 該当ファイルをスキップ
- ハッシュ計算エラー: SHA256 フィールドを null に設定
- 破損したシンボリックリンク: リンク自体のメタデータを保存

## 5. 性能特性

### 5.1 メモリ使用量

- バッファサイズ: 1000 レコード
- Row Group サイズ: 128MB
- 文字列の重複排除: PLAIN_DICTIONARY エンコーディングによる最適化

### 5.2 ディスク使用量

- Snappy 圧縮による効率的なストレージ使用
- 文字列の辞書圧縮
- 統計情報の保持（Min/Max, Null count）
