# fsparq

`fsparq`は、ファイルシステムをスキャンし、詳細なメタデータを Parquet 形式でアーカイブする高性能なクロスプラットフォームツールです。大規模なディレクトリ構造を効率的に処理しながら、メモリ使用量を最小限に抑えるように設計されています。

## 主な機能

- 🚀 **高性能**: ストリーミング処理によるメモリ効率の最適化
- 🔄 **クロスプラットフォーム**: Windows、macOS、Linux を完全サポート
- 🔍 **豊富なメタデータ**: 以下を含む包括的なファイル属性を取得
  - 正確なタイムスタンプ（作成、変更、アクセス）を UTC で記録
  - ファイルパーミッションとモード
  - プラットフォーム固有の属性（隠しファイル、システムファイル、アーカイブフラグ）
  - ファイルの SHA256 ハッシュ
- 📦 **効率的なストレージ**: Snappy 圧縮と辞書エンコーディングを使用
- ⚡ **並行処理**: goroutines を活用した並列処理

## インストール

```bash
# go installを使用
go install github.com/nrtkbb/fsparq@latest

# またはクローンしてビルド
git clone https://github.com/nrtkbb/fsparq.git
cd fsparq
go build
```

## 使用方法

基本的な使用法：

```bash
fsparq -root /path/to/scan -output metadata.parquet
```

詳細なオプション：

```bash
fsparq \
  -root /path/to/scan \
  -output metadata.parquet \
  -buffer 2000 \        # メタデータレコードのバッファサイズ
  -workers 8 \          # ワーカーgoroutineの数
  -flush 20000          # ディスクへの書き込み間隔（レコード数）
```

## 出力形式

生成される Parquet ファイルは以下のカラムを含みます：

| カラム名              | 型      | 説明                                |
| --------------------- | ------- | ----------------------------------- |
| file_path             | STRING  | ファイルの絶対パス                  |
| file_name             | STRING  | ファイル名                          |
| directory             | STRING  | 親ディレクトリのパス                |
| size_bytes            | INT64   | ファイルサイズ（バイト）            |
| creation_time_utc     | INT64   | ファイル作成日時（UTC）             |
| modification_time_utc | INT64   | 最終更新日時（UTC）                 |
| access_time_utc       | INT64   | 最終アクセス日時（UTC）             |
| file_mode             | STRING  | ファイルパーミッション（Unix 形式） |
| is_directory          | BOOLEAN | ディレクトリフラグ                  |
| is_file               | BOOLEAN | 通常ファイルフラグ                  |
| is_symlink            | BOOLEAN | シンボリックリンクフラグ            |
| is_hidden             | BOOLEAN | 隠しファイルフラグ                  |
| is_system             | BOOLEAN | システムファイルフラグ（Windows）   |
| is_archive            | BOOLEAN | アーカイブフラグ（Windows）         |
| is_readonly           | BOOLEAN | 読み取り専用フラグ                  |
| file_extension        | STRING  | ファイル拡張子（ドット付き）        |
| sha256                | STRING  | SHA256 ハッシュ（ファイルのみ）     |

## プラットフォーム固有の動作

### Windows

- Win32 API を使用して正確なファイル属性を取得
- NTFS タイムスタンプと特殊フラグ（隠し、システム、アーカイブ）をサポート
- ファイルパスはバックスラッシュ区切り

### macOS

- 利用可能な場合は正確な作成時刻を使用
- ドットで始まるファイルを隠しファイルとして判定
- ファイルパスはスラッシュ区切り

### Linux

- 作成時刻が取得できない場合は ctime をフォールバックとして使用
- システム属性とアーカイブ属性は常に false
- ファイルパスはスラッシュ区切り

## パフォーマンスに関する考慮事項

- ストリーミング処理によるメモリ使用量の最小化
- 設定可能なフラッシュ間隔によるバッファ書き込み
- Snappy 圧縮による効率的なストレージ使用
- 繰り返し文字列の辞書エンコーディング
- 設定可能なワーカー数による並行処理

## エラー処理

- パーミッションエラーが発生しても処理を継続
- アクセス不能ファイルに関する警告をログ出力
- 読み取り不能ファイルのハッシュ計算をスキップ
- 処理エラーの記録を維持

## ビルド方法

各プラットフォーム向けのビルド：

```bash
# Windows向け
GOOS=windows GOARCH=amd64 go build -o fsparq.exe

# macOS向け
GOOS=darwin GOARCH=amd64 go build -o fsparq-mac

# Linux向け
GOOS=linux GOARCH=amd64 go build -o fsparq-linux
```

## コントリビューション

コントリビューションを歓迎します！Pull Request をお気軽にご提出ください。

## ライセンス

MIT License - 詳細は[LICENSE](LICENSE)をご覧ください。
