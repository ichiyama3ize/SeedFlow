# ナレッジ管理ツール - 機能仕様書

## 1. システムアーキテクチャ

### **1.1 全体構成**

```
┌─────────────────┐
│   Web Browser   │
│   (Frontend UI) │
└─────────────────┘
         ↕ HTTP
┌──────────────────┐     ┌─────────────────┐
│   Go Main App    │     │  Python AI      │
│   - Web Server   │←───→│  Service        │
│   - Business     │HTTP │  - LangChain    │
│   - Database     │API  │  - OpenAI       │
└──────────────────┘     │  - Vector DB    │
         ↕               └─────────────────┘
┌─────────────────┐
│    SQLite       │
│   Main Database │
└─────────────────┘
```

### **1.2 プロセス構成**

- **メインプロセス**: Go Web サーバー（ポート8080）
- **AIサービス**: Python FastAPI サーバー（ポート8001）
- **データベース**: SQLiteファイル（ローカルファイルシステム）
- **フロントエンド**: Go サーバーから配信される HTML + JavaScript

### **1.3 通信方式**

- **フロントエンド ↔ Go**: HTTP REST API + Server-Side Rendering
- **Go ↔ Python**: HTTP REST API（JSON通信）
- **Go ↔ SQLite**: GORM（Go ORM）
- **Python ↔ ChromaDB**: 直接API呼び出し

## 2. 技術スタック

### **2.1 Go メインアプリケーション**

```yaml
言語・ランタイム:
  - Go: 1.21+
  - SQLite: 3.x

フレームワーク・ライブラリ:
  - Web Framework: Gin (github.com/gin-gonic/gin)
  - ORM: GORM (gorm.io/gorm)
  - SQLite Driver: gorm.io/driver/sqlite
  - HTTP Client: net/http (標準ライブラリ)
  - JSON処理: encoding/json (標準ライブラリ)

フロントエンド:
  - Template Engine: html/template (標準ライブラリ)
  - JavaScript: htmx (動的UI)
  - CSS Framework: なし（カスタムCSS）
```

### **2.2 Python AIサービス**

```yaml
言語・ランタイム:
  - Python: 3.13
  - pip: パッケージ管理

フレームワーク・ライブラリ:
  - Web Framework: FastAPI
  - ASGI Server: uvicorn
  - AI Framework: LangChain
  - LLM API: OpenAI Python SDK (GPT-4以上、モデル指定可能)
  - Vector Database: ChromaDB
  - Web Scraping: BeautifulSoup4 + requests
  - Data Validation: Pydantic (FastAPI組み込み)
```

### **2.3 データストレージ**

```yaml
メインデータ:
  - SQLite: ナレッジの構造化データ
  - ファイルシステム: アップロードファイル、エクスポート

AIデータ:
  - ChromaDB: ベクトル埋め込み
  - ファイルシステム: ローカルキャッシュ
```

## 3. データフロー仕様

### **3.1 ナレッジ作成フロー**

```json
1. User → Go: タネ投入リクエスト
   POST /api/knowledge
   {
     "seed_type": "url|text",
     "content": "データ内容"
   }

2. Go → Python: AI処理リクエスト
   POST http://localhost:8001/ai/process
   {
     "content": "処理対象テキスト",
     "seed_type": "url|text"
   }

3. Python → Go: AI処理結果
   {
     "title": "生成タイトル",
     "summary": "要約",
     "blocks": {
       "problem": "課題テキスト",
       "solution": "解決策テキスト",
       "result": "結果テキスト"
     },
     "keywords": ["keyword1", "keyword2"],
     "category": "カテゴリ名",
     "quality_score": 85
   }

4. Go → SQLite: データ永続化
   INSERT INTO knowledge (title, summary, problem, ...)

5. Go → User: 作成完了レスポンス
   {
     "id": 123,
     "status": "success",
     "knowledge": {...}
   }
```

### **3.2 検索フロー**

```
1. User → Go: 検索リクエスト
   GET /api/knowledge/search?q=検索語&category=カテゴリ

2. Go → SQLite: 基本検索実行
   SELECT * FROM knowledge
   WHERE title LIKE '%検索語%' OR keywords LIKE '%検索語%'

3. Go → Python: セマンティック検索（オプション）
   POST http://localhost:8001/ai/search
   {
     "query": "検索語",
     "candidates": [...], // SQLite検索結果
     "limit": 10
   }

4. Python → Go: 類似度スコア付き結果
   {
     "results": [
       {
         "id": 123,
         "relevance_score": 0.92,
         "matched_concepts": ["concept1", "concept2"]
       }
     ]
   }

5. Go → User: 統合検索結果
   {
     "results": [...],
     "total": 5,
     "search_type": "semantic"
   }
```

### **3.3 アウトプット生成フロー**

```
1. User → Go: エクスポートリクエスト
   POST /api/knowledge/:id/export
   {
     "format": "checklist|summary|report",
     "options": {...}
   }

2. Go → SQLite: ナレッジデータ取得
   SELECT * FROM knowledge WHERE id = :id

3. Go → Python: フォーマット変換リクエスト
   POST http://localhost:8001/ai/format
   {
     "knowledge": {...},
     "output_format": "checklist",
     "options": {...}
   }

4. Python → Go: 生成コンテンツ
   {
     "content": "# チェックリスト\n- [ ] 項目1\n...",
     "format": "markdown",
     "filename": "checklist_20241020.md"
   }

5. Go → FileSystem: ファイル保存
   /exports/checklist_20241020.md

6. Go → User: ダウンロードURL
   {
     "download_url": "/exports/checklist_20241020.md",
     "content": "...",
     "status": "success"
   }
```

## 4. データモデル仕様

### **4.1 SQLite スキーマ**

```sql
CREATE TABLE knowledge (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    summary TEXT,
    problem TEXT,
    solution TEXT,
    result TEXT,
    constraint_info TEXT,  -- "constraint" は予約語のため
    insight TEXT,
    keywords TEXT,          -- JSON文字列 ["key1","key2"]
    category TEXT,
    source_url TEXT,
    quality_score INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_title ON knowledge(title);
CREATE INDEX idx_category ON knowledge(category);
CREATE INDEX idx_keywords ON knowledge(keywords);
CREATE INDEX idx_created_at ON knowledge(created_at);
```

### **4.2 Go データ構造**

```go
// メインナレッジ構造体
type Knowledge struct {
    ID           int64     `json:"id" gorm:"primaryKey"`
    Title        string    `json:"title" gorm:"size:255;not null"`
    Summary      string    `json:"summary" gorm:"type:text"`
    Problem      string    `json:"problem" gorm:"type:text"`
    Solution     string    `json:"solution" gorm:"type:text"`
    Result       string    `json:"result" gorm:"type:text"`
    ConstraintInfo string  `json:"constraint" gorm:"column:constraint_info;type:text"`
    Insight      string    `json:"insight" gorm:"type:text"`
    Keywords     string    `json:"keywords" gorm:"type:text"` // JSON文字列
    Category     string    `json:"category" gorm:"size:100"`
    SourceURL    string    `json:"source_url" gorm:"size:500"`
    QualityScore int       `json:"quality_score"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// API リクエスト構造体
type CreateKnowledgeRequest struct {
    SeedType string `json:"seed_type" binding:"required,oneof=url text"`
    Content  string `json:"content" binding:"required"`
}

type SearchRequest struct {
    Query    string `form:"q" binding:"required"`
    Category string `form:"category"`
    Limit    int    `form:"limit,default=20"`
}

type ExportRequest struct {
    Format  string `json:"format" binding:"required,oneof=checklist summary report"`
    Options map[string]interface{} `json:"options"`
}
```

### **4.3 Python データ構造**

```python
# AI処理用モデル
class SeedInput(BaseModel):
    content: str
    seed_type: Literal["url", "text"]

class ProcessedBlocks(BaseModel):
    problem: Optional[str] = None
    solution: Optional[str] = None
    result: Optional[str] = None
    constraint: Optional[str] = None
    insight: Optional[str] = None

class ProcessedOutput(BaseModel):
    title: str
    summary: str
    blocks: ProcessedBlocks
    keywords: List[str]
    category: str
    quality_score: int

# 検索用モデル
class SearchCandidate(BaseModel):
    id: int
    title: str
    summary: str
    content: str  # 検索対象テキスト

class SemanticSearchRequest(BaseModel):
    query: str
    candidates: List[SearchCandidate]
    limit: int = 10

class SearchResult(BaseModel):
    id: int
    relevance_score: float
    matched_concepts: List[str]

# フォーマット変換用モデル
class FormatRequest(BaseModel):
    knowledge: dict  # Knowledge構造体のJSON
    output_format: Literal["checklist", "summary", "report"]
    options: Dict[str, Any] = {}

class FormatResponse(BaseModel):
    content: str
    format: str
    filename: str
```

## 5. API エンドポイント仕様

### **5.1 Go メインアプリ API**

#### **ナレッジ管理**

```http
# 作成
POST /api/knowledge
Content-Type: application/json
Body: CreateKnowledgeRequest

# 取得
GET /api/knowledge/:id
Response: Knowledge

# 一覧
GET /api/knowledge?limit=20&offset=0
Response: {results: Knowledge[], total: number}

# 更新
PUT /api/knowledge/:id
Content-Type: application/json
Body: Partial<Knowledge>

# 削除
DELETE /api/knowledge/:id
Response: {status: "success"}

# 検索
GET /api/knowledge/search?q=query&category=cat&limit=20
Response: {results: Knowledge[], total: number, search_type: string}

# エクスポート
POST /api/knowledge/:id/export
Content-Type: application/json
Body: ExportRequest
Response: {download_url: string, content: string, status: string}
```

#### **フロントエンド**

```http
# メインページ
GET /
Response: HTML (index.html)

# ナレッジ詳細ページ
GET /knowledge/:id
Response: HTML (detail.html)

# 静的ファイル
GET /static/*
Response: CSS, JS files

# エクスポートファイル
GET /exports/*
Response: Generated files (markdown, etc.)
```

### **5.2 Python AIサービス API**

```http
# コンテンツ処理
POST /ai/process
Content-Type: application/json
Body: SeedInput
Response: ProcessedOutput

# セマンティック検索
POST /ai/search
Content-Type: application/json
Body: SemanticSearchRequest
Response: {results: SearchResult[]}

# フォーマット変換
POST /ai/format
Content-Type: application/json
Body: FormatRequest
Response: FormatResponse

# ヘルスチェック
GET /ai/health
Response: {status: "ok", timestamp: string}
```

## 6. エラーハンドリング

### **6.1 Go エラーレスポンス形式**

```json
{
  "error": true,
  "message": "エラーメッセージ",
  "code": "ERROR_CODE",
  "details": {...},
  "timestamp": "2024-10-20T10:00:00Z"
}
```

### **6.2 Python エラーレスポンス形式**

```json
{
  "detail": "エラーメッセージ",
  "error_type": "ValidationError|ProcessingError|ExternalAPIError",
  "traceback": "...", // 開発時のみ
  "timestamp": "2024-10-20T10:00:00Z"
}
```

### **6.3 通信エラー処理**

- **Go → Python 通信失敗**: 3回リトライ後、フォールバック処理
- **タイムアウト**: 30秒でタイムアウト、エラーレスポンス
- **Python サービス停止**: Go 側でエラー検知、手動処理モードに切り替え

## 7. パフォーマンス要件

### **7.1 応答時間**

```yaml
Go メインAPI:
  - 基本CRUD: < 100ms
  - 検索: < 500ms
  - ページ表示: < 200ms

Python AIサービス:
  - テキスト処理: < 10秒
  - URL処理: < 15秒
  - セマンティック検索: < 2秒
  - フォーマット変換: < 5秒
```

### **7.2 データ制限**

```yaml
リクエストサイズ:
  - テキスト入力: 最大50KB
  - URL処理: 最大1MB (取得後)
  - ファイルアップロード: 最大10MB

データベース:
  - ナレッジ件数: 最大10,000件
  - SQLiteファイル: 最大1GB
```

## 8. 起動・停止手順

### **8.1 開発環境起動**

```bash
# 1. Python AIサービス起動
cd ai-service
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt
uvicorn main:app --host 127.0.0.1 --port 8001

# 2. Go メインアプリ起動（別ターミナル）
cd go-app
go mod tidy
go run .
```

### **8.2 プロダクション起動**

```bash
# バックグラウンド起動
nohup python -m uvicorn ai-service.main:app --host 127.0.0.1 --port 8001 &
nohup go run go-app/. &

# または systemd service として設定
```

### **8.3 ヘルスチェック**

```bash
# Go サービス確認
curl http://localhost:8080/api/health

# Python サービス確認
curl http://localhost:8001/ai/health

# 統合動作確認
curl -X POST http://localhost:8080/api/knowledge \
  -H "Content-Type: application/json" \
  -d '{"seed_type":"text","content":"テストデータ"}'
```

## 9. ログ機能仕様

### **9.1 Go ログ設計**

#### **ログレベル定義**

```go
// logger.go
import (
    "github.com/sirupsen/logrus"
    "os"
)

type LogLevel string
const (
    DEBUG LogLevel = "debug"
    INFO  LogLevel = "info"
    WARN  LogLevel = "warn"
    ERROR LogLevel = "error"
)

type Logger struct {
    *logrus.Logger
    requestID string
}

func NewLogger() *Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetOutput(os.Stdout)

    return &Logger{Logger: logger}
}
```

#### **構造化ログ出力**

```go
// リクエスト処理ログ
func (l *Logger) LogRequest(method, path string, duration time.Duration, status int) {
    l.WithFields(logrus.Fields{
        "method":     method,
        "path":       path,
        "duration_ms": duration.Milliseconds(),
        "status":     status,
        "request_id": l.requestID,
        "component":  "go-main",
    }).Info("HTTP Request")
}

// AI呼び出しログ
func (l *Logger) LogAIRequest(endpoint string, payload interface{}, duration time.Duration, success bool) {
    l.WithFields(logrus.Fields{
        "endpoint":    endpoint,
        "payload_size": len(fmt.Sprintf("%v", payload)),
        "duration_ms": duration.Milliseconds(),
        "success":     success,
        "request_id":  l.requestID,
        "component":   "ai-client",
    }).Info("AI Service Call")
}

// データベース操作ログ
func (l *Logger) LogDatabase(operation string, table string, recordID interface{}, duration time.Duration) {
    l.WithFields(logrus.Fields{
        "operation":   operation, // CREATE, READ, UPDATE, DELETE
        "table":       table,
        "record_id":   recordID,
        "duration_ms": duration.Milliseconds(),
        "request_id":  l.requestID,
        "component":   "database",
    }).Info("Database Operation")
}
```

#### **ログ出力例**

```json
{
  "level": "info",
  "msg": "HTTP Request",
  "method": "POST",
  "path": "/api/knowledge",
  "duration_ms": 1250,
  "status": 201,
  "request_id": "req_abc123",
  "component": "go-main",
  "timestamp": "2024-10-20T10:00:00Z"
}
```

### **9.2 Python ログ設計**

#### **構造化ログ設定**

```python
# logger.py
import logging
import json
import time
from typing import Dict, Any

class StructuredLogger:
    def __init__(self, name: str):
        self.logger = logging.getLogger(name)
        handler = logging.StreamHandler()
        formatter = logging.Formatter('%(message)s')
        handler.setFormatter(formatter)
        self.logger.addHandler(handler)
        self.logger.setLevel(logging.INFO)

    def log(self, level: str, message: str, **kwargs):
        log_data = {
            "level": level,
            "message": message,
            "timestamp": time.time(),
            "component": "python-ai",
            **kwargs
        }
        self.logger.info(json.dumps(log_data))

# AI処理ログ
def log_ai_processing(func):
    def wrapper(*args, **kwargs):
        start_time = time.time()
        try:
            result = func(*args, **kwargs)
            duration = (time.time() - start_time) * 1000

            logger.log("info", "AI Processing Success",
                function=func.__name__,
                duration_ms=duration,
                input_size=len(str(args[0])) if args else 0
            )
            return result
        except Exception as e:
            duration = (time.time() - start_time) * 1000
            logger.log("error", "AI Processing Failed",
                function=func.__name__,
                duration_ms=duration,
                error=str(e)
            )
            raise
    return wrapper
```

#### **LangChain処理ログ**

```python
# chains.py
@log_ai_processing
async def process_content_with_langchain(content: str, seed_type: str) -> ProcessedOutput:
    logger.log("info", "LangChain Processing Start",
        content_length=len(content),
        seed_type=seed_type
    )

    # LangChain実行
    chain = create_processing_chain()
    result = await chain.arun(content=content)

    logger.log("info", "LangChain Processing Complete",
        blocks_generated=len(result.get("blocks", {})),
        keywords_count=len(result.get("keywords", [])),
        quality_score=result.get("quality_score", 0)
    )

    return result
```

### **9.3 ログ集約・分析**

#### **ログファイル出力**

```yaml
# Go アプリ
logs/go-app/
  ├── app_2024-10-20.log      # 日付別ログファイル
  ├── error_2024-10-20.log    # エラーログ分離
  └── access_2024-10-20.log   # アクセスログ

# Python AI サービス
logs/ai-service/
  ├── ai_2024-10-20.log       # AI処理ログ
  ├── error_2024-10-20.log    # エラーログ
  └── performance_2024-10-20.log  # パフォーマンスログ
```

#### **ログ分析用スクリプト**

```bash
# tools/log_analyzer.sh
#!/bin/bash

# エラー率の監視
echo "=== エラー率分析 ==="
grep '"level":"error"' logs/*/*.log | wc -l

# AI処理時間の分析
echo "=== AI処理時間分析 ==="
grep "AI Processing" logs/ai-service/*.log | \
  jq '.duration_ms' | \
  awk '{sum+=$1; count++} END {print "平均:", sum/count "ms"}'

# よく使われる検索クエリ
echo "=== 検索クエリ分析 ==="
grep '"path":"/api/knowledge/search"' logs/go-app/*.log | \
  jq -r '.query' | sort | uniq -c | sort -nr
```

## 10. テスト駆動設計

### **10.1 Go テスト構造**

#### **テストファイル構成**

```
go-app/
├── handlers_test.go        # API ハンドラーのテスト
├── models_test.go         # データモデルのテスト
├── database_test.go       # データベース操作のテスト
├── ai_client_test.go      # AI クライアントのテスト
├── integration_test.go    # 統合テスト
└── testdata/
    ├── sample_knowledge.json
    └── mock_responses.json
```

#### **単体テスト例**

```go
// handlers_test.go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"))
    db.AutoMigrate(&Knowledge{})
    return db
}

func TestCreateKnowledge(t *testing.T) {
    // テスト用DB・ルーター設定
    db := setupTestDB()
    router := setupRouter(db, &MockAIClient{})

    // テストケース
    tests := []struct {
        name     string
        payload  CreateKnowledgeRequest
        expected int // HTTP status
    }{
        {
            name: "正常なテキスト投入",
            payload: CreateKnowledgeRequest{
                SeedType: "text",
                Content:  "テスト用のナレッジ内容",
            },
            expected: 201,
        },
        {
            name: "無効な種別",
            payload: CreateKnowledgeRequest{
                SeedType: "invalid",
                Content:  "テスト",
            },
            expected: 400,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jsonPayload, _ := json.Marshal(tt.payload)
            req := httptest.NewRequest("POST", "/api/knowledge",
                bytes.NewBuffer(jsonPayload))
            req.Header.Set("Content-Type", "application/json")

            recorder := httptest.NewRecorder()
            router.ServeHTTP(recorder, req)

            assert.Equal(t, tt.expected, recorder.Code)

            if recorder.Code == 201 {
                var response Knowledge
                json.NewDecoder(recorder.Body).Decode(&response)
                assert.NotEmpty(t, response.Title)
                assert.NotEmpty(t, response.ID)
            }
        })
    }
}
```

#### **モック AI クライアント**

```go
// ai_client_test.go
type MockAIClient struct{}

func (m *MockAIClient) ProcessSeed(content, seedType string) (*ProcessedOutput, error) {
    return &ProcessedOutput{
        Title:    "モックタイトル",
        Summary:  "モック要約",
        Blocks: ProcessedBlocks{
            Problem:  "モック課題",
            Solution: "モック解決策",
        },
        Keywords:   []string{"mock", "test"},
        Category:   "テスト",
        Confidence: 0.85,
    }, nil
}

func (m *MockAIClient) Search(query string, candidates []Knowledge) ([]SearchResult, error) {
    return []SearchResult{
        {ID: 1, RelevanceScore: 0.9, MatchedConcepts: []string{"mock"}},
    }, nil
}
```

### **10.2 Python テスト構造**

#### **テストファイル構成**

```
ai-service/
├── test_processors.py     # AI処理ロジックのテスト
├── test_chains.py         # LangChain設定のテスト
├── test_api.py           # API エンドポイントのテスト
├── conftest.py           # pytest 設定・フィクスチャ
└── test_data/
    ├── sample_urls.json
    ├── sample_texts.json
    └── expected_outputs.json
```

#### **pytest テスト例**

```python
# test_processors.py
import pytest
from unittest.mock import Mock, patch
from processors import KnowledgeProcessor

@pytest.fixture
def processor():
    return KnowledgeProcessor()

@pytest.fixture
def mock_llm_response():
    return {
        "title": "テストタイトル",
        "summary": "テスト要約",
        "blocks": {
            "problem": "テスト課題",
            "solution": "テスト解決策"
        },
        "keywords": ["test", "mock"],
        "category": "テスト"
    }

class TestKnowledgeProcessor:
    @patch('processors.OpenAI')
    async def test_process_text_success(self, mock_openai, processor, mock_llm_response):
        # モックの設定
        mock_openai.return_value.arun.return_value = mock_llm_response

        # テスト実行
        result = await processor.process_text("テスト入力")

        # アサーション
        assert result.title == "テストタイトル"
        assert result.quality_score > 0
        assert "test" in result.keywords

    @patch('requests.get')
    async def test_process_url_success(self, mock_requests, processor):
        # モックHTTPレスポンス
        mock_response = Mock()
        mock_response.content = "<html><title>Test</title><body>Test content</body></html>"
        mock_requests.return_value = mock_response

        # テスト実行
        result = await processor.process_url("https://example.com")

        # アサーション
        assert result.title is not None
        assert result.blocks.problem is not None or result.blocks.solution is not None

class TestAPIEndpoints:
    @pytest.mark.asyncio
    async def test_process_endpoint(self, client):
        response = await client.post("/ai/process", json={
            "content": "CTRが0.5%から1.2%に改善した",
            "seed_type": "text"
        })

        assert response.status_code == 200
        data = response.json()
        assert data["title"] != ""
        assert data["quality_score"] > 0
```

#### **テスト用フィクスチャ**

```python
# conftest.py
import pytest
from fastapi.testclient import TestClient
from main import app

@pytest.fixture
def client():
    return TestClient(app)

@pytest.fixture
def sample_knowledge_data():
    return {
        "title": "サンプルナレッジ",
        "summary": "テスト用のサンプルデータ",
        "blocks": {
            "problem": "テスト課題",
            "solution": "テスト解決策",
            "result": "テスト結果"
        }
    }

@pytest.fixture
def mock_openai_response():
    return """
    タイトル: テスト記事の分析
    要約: この記事ではAIを活用した業務改善について述べている
    課題/問題: 手動処理による時間コスト
    解決策: AI自動化ツールの導入
    結果/効果: 作業時間50%削減
    """
```

### **10.3 統合テスト設計**

#### **エンドツーエンドテストシナリオ**

```go
// integration_test.go
func TestFullWorkflow(t *testing.T) {
    // テスト用サービス起動
    testDB := setupTestDB()
    mockAI := &MockAIClient{}
    server := startTestServer(testDB, mockAI)
    defer server.Close()

    // シナリオ1: ナレッジ作成→検索→取得
    t.Run("Complete Knowledge Creation Workflow", func(t *testing.T) {
        // 1. ナレッジ作成
        createResp := createTestKnowledge(t, server, "テスト内容")
        assert.Equal(t, 201, createResp.StatusCode)

        var knowledge Knowledge
        json.NewDecoder(createResp.Body).Decode(&knowledge)
        knowledgeID := knowledge.ID

        // 2. 検索で見つかることを確認
        searchResp := searchTestKnowledge(t, server, "テスト")
        assert.Equal(t, 200, searchResp.StatusCode)

        var searchResult SearchResponse
        json.NewDecoder(searchResp.Body).Decode(&searchResult)
        assert.True(t, len(searchResult.Results) > 0)
        assert.Equal(t, knowledgeID, searchResult.Results[0].ID)

        // 3. 詳細取得
        detailResp := getTestKnowledge(t, server, knowledgeID)
        assert.Equal(t, 200, detailResp.StatusCode)

        var detail Knowledge
        json.NewDecoder(detailResp.Body).Decode(&detail)
        assert.Equal(t, knowledge.Title, detail.Title)
    })
}
```

### **10.4 パフォーマンステスト**

#### **負荷テスト**

```go
// performance_test.go
func TestConcurrentKnowledgeCreation(t *testing.T) {
    server := startTestServer()
    defer server.Close()

    concurrency := 10
    requests := 50

    results := make(chan TestResult, requests)

    // 並行リクエスト実行
    for i := 0; i < concurrency; i++ {
        go func() {
            for j := 0; j < requests/concurrency; j++ {
                start := time.Now()
                resp := createTestKnowledge(server, fmt.Sprintf("テスト%d", j))
                duration := time.Since(start)

                results <- TestResult{
                    StatusCode: resp.StatusCode,
                    Duration:   duration,
                }
            }
        }()
    }

    // 結果検証
    var totalDuration time.Duration
    successCount := 0

    for i := 0; i < requests; i++ {
        result := <-results
        totalDuration += result.Duration
        if result.StatusCode == 201 {
            successCount++
        }
    }

    avgDuration := totalDuration / time.Duration(requests)
    successRate := float64(successCount) / float64(requests)

    assert.True(t, avgDuration < 5*time.Second, "平均応答時間が5秒を超過")
    assert.True(t, successRate >= 0.95, "成功率が95%を下回る")
}
```

### **10.5 継続的インテグレーション**

#### **テスト実行スクリプト**

```bash
#!/bin/bash
# scripts/run_tests.sh

echo "=== Go テスト実行 ==="
cd go-app
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo "=== Python テスト実行 ==="
cd ../ai-service
python -m pytest -v --cov=. --cov-report=html

echo "=== 統合テスト実行 ==="
cd ../go-app
go test -v -tags=integration ./...

echo "=== ログ分析 ==="
../tools/log_analyzer.sh
```

#### **Makefile**

```makefile
# Makefile
.PHONY: test test-go test-python test-integration

test: test-go test-python test-integration

test-go:
	cd go-app && go test -v ./...

test-python:
	cd ai-service && python -m pytest -v

test-integration:
	cd go-app && go test -v -tags=integration ./...

start-services:
	cd ai-service && uvicorn main:app --host 127.0.0.1 --port 8001 &
	cd go-app && go run .

clean:
	rm -rf logs/
	rm -rf go-app/coverage.*
	rm -rf ai-service/htmlcov/
```

## 11. 開発ワークフロー

### **11.1 TDD サイクル**

```
1. テスト作成 (Red)
   └ 期待する動作の定義

2. 最小実装 (Green)
   └ テストが通る最小コード

3. リファクタリング (Refactor)
   └ コード品質向上

4. ログ確認
   └ 動作状況の詳細確認
```

### **11.2 開発フロー例**

```bash
# 1. 新機能のテスト作成
echo "新機能のテストを作成"
go test -v ./handlers_test.go -run TestNewFeature

# 2. テスト実行（失敗確認）
make test-go

# 3. 実装
echo "最小実装でテストを通す"

# 4. テスト実行（成功確認）
make test-go

# 5. ログ確認
tail -f logs/go-app/app_$(date +%Y-%m-%d).log

# 6. リファクタリング + テスト
make test
```

### **11.3 デバッグ支援**

#### **ログベースデバッグ**

```bash
# リアルタイムログ監視
tail -f logs/go-app/app_$(date +%Y-%m-%d).log | jq .

# エラーログのみフィルタ
tail -f logs/*/error_$(date +%Y-%m-%d).log

# AI処理のパフォーマンス監視
grep "AI Processing" logs/ai-service/*.log | tail -20 | jq '.duration_ms'
```

#### **テスト用ユーティリティ**

```go
// test_utils.go
func LogTestResult(t *testing.T, operation string, duration time.Duration, success bool) {
    logger.WithFields(logrus.Fields{
        "test_name":   t.Name(),
        "operation":   operation,
        "duration_ms": duration.Milliseconds(),
        "success":     success,
        "component":   "test",
    }).Info("Test Execution")
}

func AssertWithLog(t *testing.T, condition bool, message string) {
    if !condition {
        logger.WithFields(logrus.Fields{
            "test_name": t.Name(),
            "assertion": message,
            "component": "test",
        }).Error("Assertion Failed")
    }
    assert.True(t, condition, message)
}
```

この設計により、**開発中の動作状況が詳細に把握でき**、**テスト駆動で品質を保証**できるシステムになります。
