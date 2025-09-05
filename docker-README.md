# SeedFlow - Docker環境セットアップ

SeedFlow ナレッジ管理ツールをDocker環境で実行するためのガイドです。

## 📋 前提条件

- Docker Engine 20.10.0以上
- Docker Compose 2.0以上
- OpenAI API キー（必須）
- Claude API キー（オプション）

### 動作確認

```bash
docker --version
docker-compose --version
# または
docker compose version
```

## 🚀 クイックスタート

### 1. 環境変数の設定

```bash
# .env.example をコピーして設定ファイルを作成
cp .env.example .env

# .envファイルを編集してAPIキーを設定
nano .env
```

最低限必要な設定:
```env
OPENAI_API_KEY=sk-your-openai-api-key-here
```

### 2. アプリケーションの起動

```bash
# 簡単起動（推奨）
./scripts/start.sh

# または手動起動
docker-compose up --build -d
```

### 3. アクセス確認

- **Web UI**: http://localhost:8080
- **AI Service**: http://localhost:8001

## 📁 ディレクトリ構造

```
seedflow/
├── docker-compose.yml          # Docker Compose設定
├── Dockerfile.go               # Go アプリケーション用Dockerfile
├── Dockerfile.ai               # Python AIサービス用Dockerfile
├── .env.example                # 環境変数テンプレート
├── .env                        # 環境変数設定（作成が必要）
├── .dockerignore               # Docker無視ファイル
├── config/
│   └── config.yaml             # アプリケーション設定
├── scripts/                    # 管理スクリプト
│   ├── start.sh               # 起動スクリプト
│   ├── stop.sh                # 停止スクリプト
│   └── backup.sh              # バックアップスクリプト
├── data/                       # データベース（自動作成）
├── logs/                       # ログファイル（自動作成）
├── exports/                    # エクスポートファイル（自動作成）
└── backups/                    # バックアップファイル（自動作成）
```

## 🛠️ 管理コマンド

### 起動・停止

```bash
# 起動
./scripts/start.sh

# 停止
./scripts/stop.sh

# 再起動
docker-compose restart
```

### ログ確認

```bash
# 全サービスのログ
docker-compose logs -f

# 特定サービスのログ
docker-compose logs -f go-app
docker-compose logs -f ai-service
```

### バックアップ

```bash
# バックアップ作成
./scripts/backup.sh

# 手動バックアップ
docker-compose exec go-app sqlite3 /app/data/knowledge.db ".backup /app/backups/manual_backup.db"
```

### メンテナンス

```bash
# コンテナの状態確認
docker-compose ps

# リソース使用状況
docker-compose top

# コンテナ内でコマンド実行
docker-compose exec go-app sh
docker-compose exec ai-service bash
```

## ⚙️ 設定

### 環境変数（.env）

```env
# 必須設定
OPENAI_API_KEY=sk-your-openai-api-key-here

# オプション設定
CLAUDE_API_KEY=sk-ant-your-claude-api-key-here
KNOWLEDGE_APP_PORT=8080
AI_SERVICE_PORT=8001
LOG_LEVEL=INFO
MAX_CONCURRENT_REQUESTS=10
```

### アプリケーション設定（config/config.yaml）

```yaml
app:
  port: 8080
  debug: false

ai_service:
  host: "ai-service"
  port: 8001
  timeout: 30

database:
  path: "/app/data/knowledge.db"
  backup_enabled: true
```

## 🔧 トラブルシューティング

### よくある問題

#### 1. ポート競合エラー

```bash
# ポート使用状況確認
netstat -tulpn | grep :8080
netstat -tulpn | grep :8001

# 別のポートを使用する場合
export KNOWLEDGE_APP_PORT=8081
export AI_SERVICE_PORT=8002
docker-compose up -d
```

#### 2. API キーエラー

```bash
# 環境変数確認
docker-compose exec go-app env | grep API_KEY
docker-compose exec ai-service env | grep API_KEY

# .envファイルの確認
cat .env | grep API_KEY
```

#### 3. データベース権限エラー

```bash
# データディレクトリの権限確認
ls -la data/

# 権限修正
sudo chown -R $USER:$USER data/
chmod 700 data/
```

#### 4. コンテナが起動しない

```bash
# ビルドログ確認
docker-compose build --no-cache

# 詳細ログ確認
docker-compose up --build

# コンテナ状態確認
docker-compose ps
docker inspect <container_name>
```

### ヘルスチェック

```bash
# Go アプリケーション
curl http://localhost:8080/api/health

# AI サービス
curl http://localhost:8001/ai/health

# コンテナ内でのヘルスチェック
docker-compose exec go-app wget -q -O- http://localhost:8080/api/health
docker-compose exec ai-service curl http://localhost:8001/ai/health
```

## 📊 監視

### リソース監視

```bash
# コンテナリソース使用状況
docker stats

# ディスク使用状況
docker system df

# ログサイズ確認
du -sh logs/
```

### パフォーマンス監視

```bash
# レスポンス時間測定
time curl http://localhost:8080/api/health

# データベースサイズ確認
docker-compose exec go-app ls -lh /app/data/
```

## 🔄 アップデート

### アプリケーションアップデート

```bash
# 1. バックアップ作成
./scripts/backup.sh

# 2. 最新コードを取得
git pull origin main

# 3. コンテナ再ビルド・再起動
docker-compose down
docker-compose up --build -d
```

### Docker イメージアップデート

```bash
# ベースイメージ更新
docker-compose pull
docker-compose up -d --force-recreate
```

## 🗑️ クリーンアップ

### 通常のクリーンアップ

```bash
# コンテナ停止・削除
docker-compose down

# 未使用リソース削除
docker system prune
```

### 完全クリーンアップ

```bash
# 全データ削除（注意：データベースも削除されます）
docker-compose down --rmi all --volumes
rm -rf data/ logs/ exports/

# Dockerシステム全体クリーンアップ
docker system prune -a --volumes
```

## 🛡️ セキュリティ

### 基本的なセキュリティ対策

1. **APIキーの管理**
   - `.env`ファイルはバージョン管理に含めない
   - 本番環境では環境変数で設定

2. **ネットワークセキュリティ**
   - 必要なポートのみ公開
   - プライベートネットワーク使用

3. **データ保護**
   - 定期バックアップ
   - データディレクトリの権限制限

### 本番環境での追加対策

```bash
# ファイアウォール設定（例：Ubuntu）
sudo ufw allow 8080/tcp
sudo ufw enable

# Let's Encrypt SSL証明書（リバースプロキシ使用時）
# nginx-proxy + letsencrypt-nginx-proxy-companion の利用を推奨
```

## 📚 追加情報

### 関連ドキュメント

- [API仕様書](spec/api_schema.md)
- [データベース設計](spec/database_design.md)
- [デプロイメント設計](spec/deployment_design.md)

### サポート

問題が発生した場合は、以下の情報を添えてお問い合わせください：

```bash
# システム情報収集
echo "=== System Info ===" > debug_info.txt
uname -a >> debug_info.txt
docker --version >> debug_info.txt
docker-compose --version >> debug_info.txt
echo "" >> debug_info.txt

echo "=== Container Status ===" >> debug_info.txt
docker-compose ps >> debug_info.txt
echo "" >> debug_info.txt

echo "=== Recent Logs ===" >> debug_info.txt
docker-compose logs --tail=50 >> debug_info.txt
```