# ナレッジ管理ツール - APIスキーマ仕様書

## 1. API概要

### **1.1 基本情報**

```yaml
Base URL: http://localhost:8080
API Version: v1
Content-Type: application/json
Accept: application/json
```

### **1.2 認証方式**

```yaml
認証: なし（ローカル動作のため）
API Key: 不要
Rate Limiting: なし（ローカル動作のため）
```

## 2. エンドポイント一覧

### **2.1 ナレッジ管理API**

#### **POST /api/knowledge**
ナレッジを作成する

**リクエスト**
```json
{
  "seed_type": "url|text",
  "content": "string"
}
```

**レスポンス**
```json
{
  "id": 123,
  "status": "success",
  "knowledge": {
    "id": 123,
    "title": "string",
    "summary": "string",
    "problem": "string",
    "solution": "string",
    "result": "string",
    "constraint": "string",
    "insight": "string",
    "keywords": ["string"],
    "category": "string",
    "source_url": "string",
    "quality_score": 85,
    "created_at": "2024-10-20T10:00:00Z",
    "updated_at": "2024-10-20T10:00:00Z"
  }
}
```

**エラーレスポンス**
```json
{
  "error": true,
  "code": "2001",
  "message": "入力内容に問題があります。内容を確認してください。",
  "details": {
    "field": "content",
    "reason": "required"
  },
  "timestamp": "2024-10-20T10:00:00Z"
}
```

#### **GET /api/knowledge/:id**
指定されたIDのナレッジを取得する

**パラメータ**
- `id` (path): ナレッジID

**レスポンス**
```json
{
  "id": 123,
  "title": "string",
  "summary": "string",
  "problem": "string",
  "solution": "string",
  "result": "string",
  "constraint": "string",
  "insight": "string",
  "keywords": ["string"],
  "category": "string",
  "source_url": "string",
  "quality_score": 85,
  "created_at": "2024-10-20T10:00:00Z",
  "updated_at": "2024-10-20T10:00:00Z"
}
```

#### **PUT /api/knowledge/:id**
指定されたIDのナレッジを更新する

**パラメータ**
- `id` (path): ナレッジID

**リクエスト**
```json
{
  "title": "string",
  "summary": "string",
  "problem": "string",
  "solution": "string",
  "result": "string",
  "constraint": "string",
  "insight": "string",
  "keywords": ["string"],
  "category": "string"
}
```

**レスポンス**
```json
{
  "id": 123,
  "status": "success",
  "knowledge": {
    "id": 123,
    "title": "string",
    "summary": "string",
    "problem": "string",
    "solution": "string",
    "result": "string",
    "constraint": "string",
    "insight": "string",
    "keywords": ["string"],
    "category": "string",
    "source_url": "string",
    "quality_score": 85,
    "created_at": "2024-10-20T10:00:00Z",
    "updated_at": "2024-10-20T10:00:00Z"
  }
}
```

#### **DELETE /api/knowledge/:id**
指定されたIDのナレッジを削除する

**パラメータ**
- `id` (path): ナレッジID

**レスポンス**
```json
{
  "status": "success",
  "message": "ナレッジが削除されました"
}
```

#### **GET /api/knowledge**
ナレッジ一覧を取得する

**クエリパラメータ**
- `limit` (query): 取得件数 (default: 20, max: 100)
- `offset` (query): オフセット (default: 0)
- `sort` (query): ソート項目 (created_at|updated_at|title|category|quality_score)
- `order` (query): ソート順 (asc|desc, default: desc)

**レスポンス**
```json
{
  "results": [
    {
      "id": 123,
      "title": "string",
      "summary": "string",
      "category": "string",
      "quality_score": 85,
      "created_at": "2024-10-20T10:00:00Z",
      "updated_at": "2024-10-20T10:00:00Z"
    }
  ],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

### **2.2 検索API**

#### **GET /api/knowledge/search**
キーワード検索を実行する

**クエリパラメータ**
- `q` (query): 検索クエリ (required)
- `category` (query): カテゴリフィルタ
- `limit` (query): 取得件数 (default: 20, max: 100)
- `offset` (query): オフセット (default: 0)

**レスポンス**
```json
{
  "results": [
    {
      "id": 123,
      "title": "string",
      "summary": "string",
      "category": "string",
      "quality_score": 85,
      "created_at": "2024-10-20T10:00:00Z",
      "updated_at": "2024-10-20T10:00:00Z",
      "highlight": {
        "title": "検索語を含む<mark>タイトル</mark>",
        "summary": "検索語を含む<mark>要約</mark>"
      }
    }
  ],
  "total": 50,
  "query": "検索語",
  "search_type": "keyword"
}
```

#### **GET /api/knowledge/filter**
フィルタ検索を実行する

**クエリパラメータ**
- `category` (query): カテゴリフィルタ
- `date_from` (query): 作成日開始 (YYYY-MM-DD)
- `date_to` (query): 作成日終了 (YYYY-MM-DD)
- `quality_min` (query): 品質スコア最小値 (0-100)
- `quality_max` (query): 品質スコア最大値 (0-100)
- `limit` (query): 取得件数 (default: 20, max: 100)
- `offset` (query): オフセット (default: 0)

**レスポンス**
```json
{
  "results": [
    {
      "id": 123,
      "title": "string",
      "summary": "string",
      "category": "string",
      "quality_score": 85,
      "created_at": "2024-10-20T10:00:00Z",
      "updated_at": "2024-10-20T10:00:00Z"
    }
  ],
  "total": 30,
  "filters": {
    "category": "マーケティング",
    "date_from": "2024-10-01",
    "quality_min": 70
  }
}
```

#### **GET /api/knowledge/related/:id**
指定されたナレッジに関連するナレッジを取得する

**パラメータ**
- `id` (path): ナレッジID

**クエリパラメータ**
- `limit` (query): 取得件数 (default: 5, max: 10)

**レスポンス**
```json
{
  "related_knowledge": [
    {
      "id": 124,
      "title": "string",
      "summary": "string",
      "category": "string",
      "relevance_score": 0.85,
      "matched_concepts": ["concept1", "concept2"]
    }
  ],
  "source_knowledge_id": 123
}
```

### **2.3 エクスポートAPI**

#### **POST /api/knowledge/:id/export**
指定されたナレッジをエクスポートする

**パラメータ**
- `id` (path): ナレッジID

**リクエスト**
```json
{
  "format": "checklist|summary|report",
  "options": {
    "include_metadata": true,
    "include_source": true
  }
}
```

**レスポンス**
```json
{
  "download_url": "/exports/checklist_20241020_123.md",
  "content": "# チェックリスト\n- [ ] 項目1\n...",
  "format": "markdown",
  "filename": "checklist_20241020_123.md",
  "status": "success"
}
```

#### **POST /api/knowledge/batch-export**
複数のナレッジを統合してエクスポートする

**リクエスト**
```json
{
  "knowledge_ids": [123, 124, 125],
  "format": "report",
  "options": {
    "include_metadata": true,
    "group_by_category": true
  }
}
```

**レスポンス**
```json
{
  "download_url": "/exports/report_20241020.md",
  "content": "# 統合レポート\n## 概要\n...",
  "format": "markdown",
  "filename": "report_20241020.md",
  "status": "success"
}
```

### **2.4 システムAPI**

#### **GET /api/health**
システムのヘルスチェック

**レスポンス**
```json
{
  "status": "healthy",
  "timestamp": "2024-10-20T10:00:00Z",
  "services": {
    "database": "healthy",
    "ai_service": "healthy"
  },
  "version": "1.0.0"
}
```

#### **GET /api/stats**
システムの統計情報を取得する

**レスポンス**
```json
{
  "total_knowledge": 150,
  "categories": {
    "マーケティング": 45,
    "開発・技術": 38,
    "営業・顧客対応": 32,
    "業務改善・効率化": 25,
    "その他": 10
  },
  "quality_distribution": {
    "90-100": 25,
    "80-89": 45,
    "70-79": 35,
    "60-69": 25,
    "50-59": 20
  },
  "recent_activity": {
    "created_today": 3,
    "created_this_week": 12,
    "created_this_month": 45
  }
}
```

## 3. エラーレスポンス

### **3.1 エラーコード一覧**

```yaml
1000番台 - システムエラー:
  1001: INTERNAL_SERVER_ERROR
  1002: SERVICE_UNAVAILABLE
  1003: DATABASE_ERROR
  1004: AI_SERVICE_ERROR

2000番台 - バリデーションエラー:
  2001: INVALID_INPUT
  2002: MISSING_REQUIRED_FIELD
  2003: INVALID_FORMAT
  2004: SIZE_LIMIT_EXCEEDED

3000番台 - ビジネスロジックエラー:
  3001: KNOWLEDGE_NOT_FOUND
  3002: DUPLICATE_KNOWLEDGE
  3003: PROCESSING_FAILED
  3004: EXPORT_FAILED

4000番台 - 外部依存エラー:
  4001: OPENAI_API_ERROR
  4002: URL_ACCESS_ERROR
  4003: NETWORK_ERROR
  4004: RATE_LIMIT_EXCEEDED
```

### **3.2 エラーレスポンス形式**

```json
{
  "error": true,
  "code": "ERROR_CODE",
  "message": "エラーメッセージ",
  "details": {
    "field": "field_name",
    "reason": "validation_error"
  },
  "timestamp": "2024-10-20T10:00:00Z"
}
```

## 4. データ型定義

### **4.1 基本データ型**

```yaml
Knowledge:
  id: integer
  title: string (max: 255)
  summary: string (max: 1000)
  problem: string (max: 5000)
  solution: string (max: 5000)
  result: string (max: 5000)
  constraint: string (max: 2000)
  insight: string (max: 2000)
  keywords: array[string]
  category: string (max: 100)
  source_url: string (max: 500)
  quality_score: integer (0-100)
  created_at: datetime
  updated_at: datetime

SearchResult:
  id: integer
  title: string
  summary: string
  category: string
  quality_score: integer
  created_at: datetime
  updated_at: datetime
  highlight?: object
  relevance_score?: float
  matched_concepts?: array[string]

ExportRequest:
  format: enum[checklist, summary, report]
  options: object
```

### **4.2 バリデーションルール**

```yaml
CreateKnowledgeRequest:
  seed_type: required, enum[url, text]
  content: required, string, min: 1, max: 50000

UpdateKnowledgeRequest:
  title: string, max: 255
  summary: string, max: 1000
  problem: string, max: 5000
  solution: string, max: 5000
  result: string, max: 5000
  constraint: string, max: 2000
  insight: string, max: 2000
  keywords: array[string], max: 10
  category: string, max: 100

SearchRequest:
  q: required, string, min: 1, max: 100
  category: string, max: 100
  limit: integer, min: 1, max: 100
  offset: integer, min: 0

ExportRequest:
  format: required, enum[checklist, summary, report]
  options: object
```

## 5. 使用例

### **5.1 ナレッジ作成の流れ**

```bash
# 1. ナレッジ作成
curl -X POST http://localhost:8080/api/knowledge \
  -H "Content-Type: application/json" \
  -d '{
    "seed_type": "text",
    "content": "SNS広告のCTRが0.5%から1.2%に改善した。ターゲティングを絞り込んだことが効果的だった。"
  }'

# 2. 作成されたナレッジの確認
curl http://localhost:8080/api/knowledge/123

# 3. 検索で確認
curl "http://localhost:8080/api/knowledge/search?q=CTR"
```

### **5.2 エクスポートの流れ**

```bash
# 1. チェックリスト形式でエクスポート
curl -X POST http://localhost:8080/api/knowledge/123/export \
  -H "Content-Type: application/json" \
  -d '{
    "format": "checklist",
    "options": {
      "include_metadata": true
    }
  }'

# 2. ダウンロード
curl http://localhost:8080/exports/checklist_20241020_123.md
```

## 6. パフォーマンス要件

### **6.1 応答時間**

```yaml
基本操作:
  - ナレッジ作成: 30秒以内
  - ナレッジ取得: 1秒以内
  - ナレッジ更新: 1秒以内
  - ナレッジ削除: 1秒以内
  - 一覧取得: 2秒以内

検索操作:
  - キーワード検索: 2秒以内
  - フィルタ検索: 1秒以内
  - 関連ナレッジ取得: 3秒以内

エクスポート操作:
  - 単一ナレッジエクスポート: 5秒以内
  - 複数ナレッジエクスポート: 30秒以内
```

### **6.2 制限事項**

```yaml
リクエスト制限:
  - テキスト入力: 最大50,000文字
  - 一覧取得: 最大100件
  - 検索結果: 最大100件
  - 関連ナレッジ: 最大10件

同時処理制限:
  - ナレッジ作成: 5件まで
  - 検索処理: 10件まで
  - エクスポート: 3件まで
```

---

**文書情報**
- 作成日: 2024年10月20日
- バージョン: 1.0
- 最終更新: 2024年10月20日
- 作成者: システム開発チーム
