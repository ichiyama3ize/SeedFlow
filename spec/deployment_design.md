# ナレッジ管理ツール - デプロイメント設計書

## 1. デプロイメント概要

### **1.1 デプロイメント方針**

```yaml
基本方針:
  - ローカル動作: クラウド不要の設計
  - シンプル: 複雑な設定を避ける
  - 自動化: 可能な限り自動化
  - 保守性: 運用・保守の容易さ

デプロイメント方式:
  - 手動起動: ユーザーによる手動起動
  - ローカル実行: ローカル環境での実行
  - 設定ファイル: 設定ファイルベースの管理
```

### **1.2 システム構成**

```yaml
実行環境:
  - OS: Windows 10+, macOS 10.15+, Linux (Ubuntu 20.04+)
  - メモリ: 4GB以上
  - ストレージ: 2GB以上の空き容量
  - CPU: デュアルコア 2GHz以上
  - ネットワーク: AI処理時のインターネット接続

アプリケーション構成:
  - Go メインアプリ: ポート8080
  - Python AIサービス: ポート8001
  - SQLite データベース: ローカルファイル
  - フロントエンド: Go サーバーから配信
```

## 2. インストール手順

### **2.1 前提条件**

#### **必要なソフトウェア**
```yaml
Go:
  - バージョン: 1.21以上
  - インストール: https://golang.org/dl/
  - 確認: go version

Python:
  - バージョン: 3.10以上
  - インストール: https://python.org/downloads/
  - 確認: python --version

Git:
  - バージョン: 2.0以上
  - インストール: https://git-scm.com/downloads
  - 確認: git --version
```

#### **環境変数設定**
```bash
# OpenAI API Key
export OPENAI_API_KEY="sk-..."

# Claude API Key (オプション)
export CLAUDE_API_KEY="sk-ant-..."

# アプリケーション設定
export KNOWLEDGE_APP_PORT="8080"
export AI_SERVICE_PORT="8001"
export DATA_DIR="./data"
export LOG_DIR="./logs"
```

### **2.2 インストール手順**

#### **1. リポジトリのクローン**
```bash
# リポジトリをクローン
git clone https://github.com/your-org/knowledge-management.git
cd knowledge-management

# ディレクトリ構造の確認
ls -la
```

#### **2. Go アプリケーションのセットアップ**
```bash
# Go アプリケーションディレクトリに移動
cd go-app

# 依存関係のインストール
go mod tidy

# ビルド
go build -o knowledge-app

# 実行権限の付与
chmod +x knowledge-app
```

#### **3. Python AIサービスのセットアップ**
```bash
# Python AIサービスディレクトリに移動
cd ../ai-service

# 仮想環境の作成
python -m venv venv

# 仮想環境の有効化
# Windows
venv\Scripts\activate
# macOS/Linux
source venv/bin/activate

# 依存関係のインストール
pip install -r requirements.txt

# サービス起動テスト
python -m uvicorn main:app --host 127.0.0.1 --port 8001
```

#### **4. 設定ファイルの作成**
```bash
# 設定ディレクトリの作成
mkdir -p config

# 設定ファイルの作成
cat > config/config.yaml << EOF
# アプリケーション設定
app:
  name: "KnowledgeFlow"
  version: "1.0.0"
  port: 8080

# AIサービス設定
ai_service:
  host: "127.0.0.1"
  port: 8001
  timeout: 30

# データベース設定
database:
  path: "./data/knowledge.db"
  backup_enabled: true
  backup_interval: "24h"

# ログ設定
logging:
  level: "INFO"
  max_size: "100MB"
  max_files: 30

# セキュリティ設定
security:
  encryption_enabled: false
  file_permissions: 600
EOF
```

#### **5. ディレクトリ構造の作成**
```bash
# 必要なディレクトリの作成
mkdir -p data
mkdir -p logs
mkdir -p exports
mkdir -p backups

# 権限の設定
chmod 700 data
chmod 755 logs
chmod 755 exports
chmod 755 backups
```

## 3. 起動・停止手順

### **3.1 手動起動**

#### **開発環境での起動**
```bash
# ターミナル1: Python AIサービス起動
cd ai-service
source venv/bin/activate
python -m uvicorn main:app --host 127.0.0.1 --port 8001

# ターミナル2: Go メインアプリ起動
cd go-app
go run . --config ../config/config.yaml
```

#### **本番環境での起動**
```bash
# バックグラウンド起動
cd ai-service
source venv/bin/activate
nohup python -m uvicorn main:app --host 127.0.0.1 --port 8001 > ../logs/ai-service.log 2>&1 &

cd ../go-app
nohup ./knowledge-app --config ../config/config.yaml > ../logs/go-app.log 2>&1 &
```

### **3.2 停止手順**

#### **手動停止**
```bash
# プロセスIDの確認
ps aux | grep knowledge-app
ps aux | grep uvicorn

# プロセスの停止
kill <PID>

# または、プロセス名で停止
pkill -f knowledge-app
pkill -f uvicorn
```

#### **強制停止**
```bash
# 強制終了
kill -9 <PID>

# または、プロセス名で強制終了
pkill -9 -f knowledge-app
pkill -9 -f uvicorn
```

## 4. 設定管理

### **4.1 設定ファイル**

#### **config.yaml**
```yaml
# アプリケーション設定
app:
  name: "KnowledgeFlow"
  version: "1.0.0"
  port: 8080
  host: "127.0.0.1"
  debug: false

# AIサービス設定
ai_service:
  host: "127.0.0.1"
  port: 8001
  timeout: 30
  retry_count: 3
  models:
    openai: "gpt-4"
    claude: "claude-3-sonnet"

# データベース設定
database:
  path: "./data/knowledge.db"
  backup_enabled: true
  backup_interval: "24h"
  backup_retention: "30d"
  encryption_enabled: false

# ログ設定
logging:
  level: "INFO"
  format: "json"
  max_size: "100MB"
  max_files: 30
  compress: true

# セキュリティ設定
security:
  encryption_enabled: false
  file_permissions: 600
  directory_permissions: 700
  log_sensitive_data: false

# パフォーマンス設定
performance:
  max_concurrent_requests: 10
  request_timeout: 30
  cache_enabled: true
  cache_size: "100MB"
```

### **4.2 環境変数**

#### **必須環境変数**
```bash
# API認証情報
OPENAI_API_KEY="sk-..."
CLAUDE_API_KEY="sk-ant-..."  # オプション

# アプリケーション設定
KNOWLEDGE_APP_PORT="8080"
AI_SERVICE_PORT="8001"
DATA_DIR="./data"
LOG_DIR="./logs"
```

#### **オプション環境変数**
```bash
# セキュリティ設定
SECURITY_ENCRYPTION_ENABLED="false"
SECURITY_LOG_LEVEL="INFO"
SECURITY_FILE_PERMISSIONS="600"

# パフォーマンス設定
MAX_CONCURRENT_REQUESTS="10"
REQUEST_TIMEOUT="30"
CACHE_ENABLED="true"

# ログ設定
LOG_MAX_SIZE="100MB"
LOG_MAX_FILES="30"
LOG_COMPRESS="true"
```

## 5. バックアップ・復旧

### **5.1 バックアップ戦略**

#### **自動バックアップ**
```yaml
バックアップ対象:
  - データベース: knowledge.db
  - 設定ファイル: config.yaml
  - ログファイル: logs/*.log
  - エクスポートファイル: exports/*

バックアップ方式:
  - 頻度: 日次
  - 保持期間: 30日間
  - 圧縮: gzip
  - 暗号化: なし（ローカル動作のため）

バックアップ先:
  - ローカル: ./backups/
  - 外部: 手動コピー推奨
```

#### **手動バックアップ**
```bash
# バックアップスクリプト
#!/bin/bash
# backup.sh

BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="knowledge_backup_${DATE}.tar.gz"

# バックアップディレクトリの作成
mkdir -p $BACKUP_DIR

# バックアップの作成
tar -czf $BACKUP_DIR/$BACKUP_FILE \
  data/knowledge.db \
  config/config.yaml \
  logs/*.log \
  exports/*

echo "バックアップ完了: $BACKUP_FILE"
```

### **5.2 復旧手順**

#### **データベース復旧**
```bash
# 1. アプリケーションの停止
pkill -f knowledge-app
pkill -f uvicorn

# 2. バックアップファイルの展開
tar -xzf backups/knowledge_backup_20241020_120000.tar.gz

# 3. データベースファイルの復元
cp data/knowledge.db.backup data/knowledge.db

# 4. 権限の設定
chmod 600 data/knowledge.db

# 5. アプリケーションの再起動
cd ai-service && source venv/bin/activate && nohup python -m uvicorn main:app --host 127.0.0.1 --port 8001 > ../logs/ai-service.log 2>&1 &
cd ../go-app && nohup ./knowledge-app --config ../config/config.yaml > ../logs/go-app.log 2>&1 &
```

#### **設定ファイル復旧**
```bash
# 1. 設定ファイルの復元
cp config/config.yaml.backup config/config.yaml

# 2. 権限の設定
chmod 600 config/config.yaml

# 3. アプリケーションの再起動
pkill -f knowledge-app
cd go-app && nohup ./knowledge-app --config ../config/config.yaml > ../logs/go-app.log 2>&1 &
```

## 6. 監視・ログ

### **6.1 ログ管理**

#### **ログファイル構成**
```yaml
ログディレクトリ: ./logs/
ファイル構成:
  - go-app.log: Go アプリケーションログ
  - ai-service.log: Python AIサービスログ
  - error.log: エラーログ
  - access.log: アクセスログ
  - performance.log: パフォーマンスログ

ログ形式:
  - 形式: JSON
  - レベル: DEBUG, INFO, WARN, ERROR
  - ローテーション: 日次
  - 圧縮: gzip
```

#### **ログ監視スクリプト**
```bash
#!/bin/bash
# monitor.sh

LOG_DIR="./logs"
ERROR_THRESHOLD=10
WARN_THRESHOLD=50

# エラー率の監視
ERROR_COUNT=$(grep -c '"level":"error"' $LOG_DIR/*.log)
WARN_COUNT=$(grep -c '"level":"warn"' $LOG_DIR/*.log)

if [ $ERROR_COUNT -gt $ERROR_THRESHOLD ]; then
    echo "警告: エラー数が閾値を超過 ($ERROR_COUNT > $ERROR_THRESHOLD)"
fi

if [ $WARN_COUNT -gt $WARN_THRESHOLD ]; then
    echo "警告: 警告数が閾値を超過 ($WARN_COUNT > $WARN_THRESHOLD)"
fi

# ディスク使用量の監視
DISK_USAGE=$(df -h . | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $DISK_USAGE -gt 90 ]; then
    echo "警告: ディスク使用量が90%を超過 ($DISK_USAGE%)"
fi
```

### **6.2 ヘルスチェック**

#### **ヘルスチェックスクリプト**
```bash
#!/bin/bash
# health_check.sh

# Go アプリケーションのヘルスチェック
GO_HEALTH=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health)
if [ $GO_HEALTH -ne 200 ]; then
    echo "エラー: Go アプリケーションが応答しません (HTTP $GO_HEALTH)"
    exit 1
fi

# Python AIサービスのヘルスチェック
AI_HEALTH=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8001/ai/health)
if [ $AI_HEALTH -ne 200 ]; then
    echo "エラー: AIサービスが応答しません (HTTP $AI_HEALTH)"
    exit 1
fi

# データベースのヘルスチェック
if [ ! -f "./data/knowledge.db" ]; then
    echo "エラー: データベースファイルが見つかりません"
    exit 1
fi

echo "ヘルスチェック完了: 全サービス正常"
```

## 7. パフォーマンス最適化

### **7.1 システム最適化**

#### **OS設定**
```bash
# ファイルディスクリプタ制限の増加
ulimit -n 65536

# メモリ制限の設定
ulimit -m unlimited

# プロセス制限の設定
ulimit -p unlimited
```

#### **アプリケーション設定**
```yaml
# パフォーマンス設定
performance:
  max_concurrent_requests: 10
  request_timeout: 30
  cache_enabled: true
  cache_size: "100MB"
  connection_pool_size: 10
  keep_alive_timeout: 30
```

### **7.2 データベース最適化**

#### **SQLite設定**
```sql
-- パフォーマンス設定
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 10000;
PRAGMA temp_store = MEMORY;
PRAGMA mmap_size = 268435456;

-- インデックス最適化
ANALYZE;
REINDEX;
```

## 8. トラブルシューティング

### **8.1 よくある問題**

#### **ポート競合**
```bash
# ポート使用状況の確認
netstat -tulpn | grep :8080
netstat -tulpn | grep :8001

# プロセスの確認
lsof -i :8080
lsof -i :8001

# プロセスの停止
kill -9 <PID>
```

#### **権限エラー**
```bash
# ファイル権限の確認
ls -la data/
ls -la config/

# 権限の修正
chmod 600 data/knowledge.db
chmod 600 config/config.yaml
chmod 700 data/
chmod 700 config/
```

#### **メモリ不足**
```bash
# メモリ使用状況の確認
free -h
ps aux --sort=-%mem | head

# プロセスの再起動
pkill -f knowledge-app
pkill -f uvicorn
# 再起動
```

### **8.2 ログ分析**

#### **エラーログの分析**
```bash
# エラーログの確認
grep -i error logs/*.log

# 警告ログの確認
grep -i warn logs/*.log

# 特定の時間帯のログ
grep "2024-10-20 10:" logs/*.log
```

#### **パフォーマンスログの分析**
```bash
# 応答時間の分析
grep "duration_ms" logs/*.log | jq '.duration_ms' | sort -n

# リクエスト数の確認
grep "HTTP Request" logs/*.log | wc -l

# エラー率の計算
ERROR_COUNT=$(grep -c '"level":"error"' logs/*.log)
TOTAL_COUNT=$(grep -c '"level":"info"' logs/*.log)
ERROR_RATE=$(echo "scale=2; $ERROR_COUNT * 100 / $TOTAL_COUNT" | bc)
echo "エラー率: $ERROR_RATE%"
```

## 9. アップデート手順

### **9.1 アプリケーションアップデート**

#### **アップデート手順**
```bash
# 1. 現在のバージョンの確認
./go-app/knowledge-app --version

# 2. バックアップの作成
./backup.sh

# 3. アプリケーションの停止
pkill -f knowledge-app
pkill -f uvicorn

# 4. 新しいバージョンのダウンロード
git pull origin main

# 5. 依存関係の更新
cd go-app && go mod tidy
cd ../ai-service && pip install -r requirements.txt

# 6. アプリケーションの再ビルド
cd ../go-app && go build -o knowledge-app

# 7. アプリケーションの再起動
cd ../ai-service && source venv/bin/activate && nohup python -m uvicorn main:app --host 127.0.0.1 --port 8001 > ../logs/ai-service.log 2>&1 &
cd ../go-app && nohup ./knowledge-app --config ../config/config.yaml > ../logs/go-app.log 2>&1 &

# 8. 動作確認
./health_check.sh
```

### **9.2 データベースマイグレーション**

#### **マイグレーション手順**
```bash
# 1. データベースのバックアップ
cp data/knowledge.db data/knowledge.db.backup

# 2. マイグレーションスクリプトの実行
sqlite3 data/knowledge.db < migrations/001_add_new_column.sql

# 3. データ整合性の確認
sqlite3 data/knowledge.db "PRAGMA integrity_check;"

# 4. アプリケーションの再起動
pkill -f knowledge-app
cd go-app && nohup ./knowledge-app --config ../config/config.yaml > ../logs/go-app.log 2>&1 &
```

---

**文書情報**
- 作成日: 2024年10月20日
- バージョン: 1.0
- 最終更新: 2024年10月20日
- 作成者: システム開発チーム
