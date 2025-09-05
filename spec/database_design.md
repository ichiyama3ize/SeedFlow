# ナレッジ管理ツール - データベース設計書

## 1. データベース概要

### **1.1 基本情報**

```yaml
データベース: SQLite 3.x
ファイル名: knowledge.db
場所: ./data/knowledge.db
文字エンコーディング: UTF-8
```

### **1.2 設計方針**

- **正規化**: 第3正規形まで適用
- **インデックス**: 検索性能を重視した設計
- **論理削除**: データの復旧可能性を確保
- **拡張性**: 将来の機能追加に対応

## 2. テーブル設計

### **2.1 メインテーブル**

#### **knowledge テーブル**
ナレッジの基本情報を格納

```sql
CREATE TABLE knowledge (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL CHECK(length(title) <= 255),
    summary TEXT CHECK(length(summary) <= 1000),
    problem TEXT CHECK(length(problem) <= 5000),
    solution TEXT CHECK(length(solution) <= 5000),
    result TEXT CHECK(length(result) <= 5000),
    constraint_info TEXT CHECK(length(constraint_info) <= 2000),
    insight TEXT CHECK(length(insight) <= 2000),
    keywords TEXT CHECK(length(keywords) <= 1000), -- JSON配列
    category TEXT CHECK(length(category) <= 100),
    source_url TEXT CHECK(length(source_url) <= 500),
    quality_score INTEGER CHECK(quality_score >= 0 AND quality_score <= 100),
    estimated_expiry DATE,
    view_count INTEGER DEFAULT 0 CHECK(view_count >= 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

-- インデックス
CREATE INDEX idx_knowledge_title ON knowledge(title);
CREATE INDEX idx_knowledge_category ON knowledge(category);
CREATE INDEX idx_knowledge_keywords ON knowledge(keywords);
CREATE INDEX idx_knowledge_created_at ON knowledge(created_at);
CREATE INDEX idx_knowledge_updated_at ON knowledge(updated_at);
CREATE INDEX idx_knowledge_quality_score ON knowledge(quality_score);
CREATE INDEX idx_knowledge_deleted_at ON knowledge(deleted_at);
CREATE INDEX idx_knowledge_source_url ON knowledge(source_url);
```

#### **knowledge_metadata テーブル**
ナレッジのメタデータを格納

```sql
CREATE TABLE knowledge_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    knowledge_id INTEGER NOT NULL,
    source_type TEXT NOT NULL CHECK(source_type IN ('url', 'text', 'manual')),
    processing_time_ms INTEGER CHECK(processing_time_ms >= 0),
    ai_model_version TEXT,
    last_accessed_at DATETIME,
    access_count INTEGER DEFAULT 0 CHECK(access_count >= 0),
    related_knowledge_ids TEXT, -- JSON配列
    tags TEXT, -- JSON配列
    custom_fields TEXT, -- JSON形式
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX idx_metadata_knowledge_id ON knowledge_metadata(knowledge_id);
CREATE INDEX idx_metadata_source_type ON knowledge_metadata(source_type);
CREATE INDEX idx_metadata_last_accessed ON knowledge_metadata(last_accessed_at);
CREATE INDEX idx_metadata_access_count ON knowledge_metadata(access_count);
```

### **2.2 拡張テーブル**

#### **knowledge_relations テーブル**
ナレッジ間の関係性を格納

```sql
CREATE TABLE knowledge_relations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_knowledge_id INTEGER NOT NULL,
    target_knowledge_id INTEGER NOT NULL,
    relation_type TEXT NOT NULL CHECK(relation_type IN ('similar', 'complementary', 'causal', 'conflicting', 'evolution')),
    relevance_score REAL CHECK(relevance_score >= 0.0 AND relevance_score <= 1.0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE,
    FOREIGN KEY (target_knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE,
    UNIQUE(source_knowledge_id, target_knowledge_id, relation_type)
);

-- インデックス
CREATE INDEX idx_relations_source ON knowledge_relations(source_knowledge_id);
CREATE INDEX idx_relations_target ON knowledge_relations(target_knowledge_id);
CREATE INDEX idx_relations_type ON knowledge_relations(relation_type);
CREATE INDEX idx_relations_score ON knowledge_relations(relevance_score);
```

#### **knowledge_usage_log テーブル**
ナレッジの利用履歴を格納

```sql
CREATE TABLE knowledge_usage_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    knowledge_id INTEGER NOT NULL,
    action_type TEXT NOT NULL CHECK(action_type IN ('view', 'search', 'export', 'edit')),
    user_context TEXT, -- 利用時の文脈情報
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX idx_usage_knowledge_id ON knowledge_usage_log(knowledge_id);
CREATE INDEX idx_usage_action_type ON knowledge_usage_log(action_type);
CREATE INDEX idx_usage_created_at ON knowledge_usage_log(created_at);
```

#### **system_config テーブル**
システム設定を格納

```sql
CREATE TABLE system_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key TEXT NOT NULL UNIQUE,
    config_value TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- インデックス
CREATE INDEX idx_config_key ON system_config(config_key);
```

## 3. データ型定義

### **3.1 基本データ型**

```yaml
Knowledge:
  id: INTEGER PRIMARY KEY AUTOINCREMENT
  title: TEXT NOT NULL (最大255文字)
  summary: TEXT (最大1000文字)
  problem: TEXT (最大5000文字)
  solution: TEXT (最大5000文字)
  result: TEXT (最大5000文字)
  constraint_info: TEXT (最大2000文字)
  insight: TEXT (最大2000文字)
  keywords: TEXT (JSON配列、最大1000文字)
  category: TEXT (最大100文字)
  source_url: TEXT (最大500文字)
  quality_score: INTEGER (0-100)
  estimated_expiry: DATE
  view_count: INTEGER DEFAULT 0
  created_at: DATETIME DEFAULT CURRENT_TIMESTAMP
  updated_at: DATETIME DEFAULT CURRENT_TIMESTAMP
  deleted_at: DATETIME NULL

KnowledgeMetadata:
  id: INTEGER PRIMARY KEY AUTOINCREMENT
  knowledge_id: INTEGER (外部キー)
  source_type: TEXT (url|text|manual)
  processing_time_ms: INTEGER
  ai_model_version: TEXT
  last_accessed_at: DATETIME
  access_count: INTEGER DEFAULT 0
  related_knowledge_ids: TEXT (JSON配列)
  tags: TEXT (JSON配列)
  custom_fields: TEXT (JSON形式)
  created_at: DATETIME DEFAULT CURRENT_TIMESTAMP
  updated_at: DATETIME DEFAULT CURRENT_TIMESTAMP

KnowledgeRelations:
  id: INTEGER PRIMARY KEY AUTOINCREMENT
  source_knowledge_id: INTEGER (外部キー)
  target_knowledge_id: INTEGER (外部キー)
  relation_type: TEXT (similar|complementary|causal|conflicting|evolution)
  relevance_score: REAL (0.0-1.0)
  created_at: DATETIME DEFAULT CURRENT_TIMESTAMP

KnowledgeUsageLog:
  id: INTEGER PRIMARY KEY AUTOINCREMENT
  knowledge_id: INTEGER (外部キー)
  action_type: TEXT (view|search|export|edit)
  user_context: TEXT
  created_at: DATETIME DEFAULT CURRENT_TIMESTAMP

SystemConfig:
  id: INTEGER PRIMARY KEY AUTOINCREMENT
  config_key: TEXT UNIQUE
  config_value: TEXT
  description: TEXT
  created_at: DATETIME DEFAULT CURRENT_TIMESTAMP
  updated_at: DATETIME DEFAULT CURRENT_TIMESTAMP
```

### **3.2 JSON形式のデータ**

#### **keywords フィールド**
```json
["マーケティング", "SNS", "CTR", "改善"]
```

#### **related_knowledge_ids フィールド**
```json
[123, 124, 125]
```

#### **tags フィールド**
```json
["重要", "実績あり", "再現可能"]
```

#### **custom_fields フィールド**
```json
{
  "industry": "IT",
  "company_size": "中小企業",
  "difficulty": "中級"
}
```

## 4. インデックス設計

### **4.1 検索性能最適化**

```sql
-- 全文検索用インデックス
CREATE INDEX idx_knowledge_title_search ON knowledge(title);
CREATE INDEX idx_knowledge_summary_search ON knowledge(summary);
CREATE INDEX idx_knowledge_problem_search ON knowledge(problem);
CREATE INDEX idx_knowledge_solution_search ON knowledge(solution);

-- 複合インデックス
CREATE INDEX idx_knowledge_category_created ON knowledge(category, created_at);
CREATE INDEX idx_knowledge_quality_created ON knowledge(quality_score, created_at);
CREATE INDEX idx_knowledge_active ON knowledge(deleted_at, created_at);
```

### **4.2 パフォーマンス最適化**

```sql
-- 統計情報更新
ANALYZE knowledge;
ANALYZE knowledge_metadata;
ANALYZE knowledge_relations;
ANALYZE knowledge_usage_log;
```

## 5. 制約・ルール

### **5.1 データ整合性制約**

```sql
-- 外部キー制約
ALTER TABLE knowledge_metadata 
ADD CONSTRAINT fk_metadata_knowledge 
FOREIGN KEY (knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE;

ALTER TABLE knowledge_relations 
ADD CONSTRAINT fk_relations_source 
FOREIGN KEY (source_knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE;

ALTER TABLE knowledge_relations 
ADD CONSTRAINT fk_relations_target 
FOREIGN KEY (target_knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE;

ALTER TABLE knowledge_usage_log 
ADD CONSTRAINT fk_usage_knowledge 
FOREIGN KEY (knowledge_id) REFERENCES knowledge(id) ON DELETE CASCADE;

-- チェック制約
ALTER TABLE knowledge 
ADD CONSTRAINT chk_quality_score 
CHECK (quality_score >= 0 AND quality_score <= 100);

ALTER TABLE knowledge 
ADD CONSTRAINT chk_view_count 
CHECK (view_count >= 0);

ALTER TABLE knowledge_metadata 
ADD CONSTRAINT chk_access_count 
CHECK (access_count >= 0);

ALTER TABLE knowledge_relations 
ADD CONSTRAINT chk_relevance_score 
CHECK (relevance_score >= 0.0 AND relevance_score <= 1.0);
```

### **5.2 ビジネスルール**

```sql
-- トリガー: updated_at自動更新
CREATE TRIGGER update_knowledge_timestamp 
AFTER UPDATE ON knowledge 
BEGIN
    UPDATE knowledge SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_metadata_timestamp 
AFTER UPDATE ON knowledge_metadata 
BEGIN
    UPDATE knowledge_metadata SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- トリガー: 論理削除時のメタデータ更新
CREATE TRIGGER soft_delete_knowledge 
AFTER UPDATE OF deleted_at ON knowledge 
BEGIN
    UPDATE knowledge_metadata 
    SET last_accessed_at = CURRENT_TIMESTAMP 
    WHERE knowledge_id = NEW.id AND NEW.deleted_at IS NOT NULL;
END;
```

## 6. 初期データ

### **6.1 システム設定**

```sql
INSERT INTO system_config (config_key, config_value, description) VALUES
('ai_model_version', 'gpt-4', '使用するAIモデルのバージョン'),
('max_knowledge_count', '10000', '最大ナレッジ数'),
('default_quality_threshold', '50', 'デフォルト品質閾値'),
('auto_cleanup_days', '365', '自動クリーンアップ日数'),
('backup_enabled', 'true', 'バックアップ有効フラグ'),
('backup_interval_hours', '24', 'バックアップ間隔（時間）');
```

### **6.2 カテゴリ定義**

```sql
-- カテゴリは動的に管理されるため、初期データは不要
-- システム定義カテゴリ:
-- - マーケティング
-- - 開発・技術
-- - 営業・顧客対応
-- - 業務改善・効率化
-- - その他
```

## 7. バックアップ・復旧

### **7.1 バックアップ戦略**

```sql
-- 完全バックアップ
.backup main backup_20241020.db

-- 差分バックアップ（SQLiteの制限により手動実装）
-- 1. 前回バックアップ以降の変更をエクスポート
-- 2. 変更分のみをバックアップファイルに追加
```

### **7.2 復旧手順**

```sql
-- 1. バックアップファイルから復元
.restore backup_20241020.db

-- 2. データ整合性チェック
PRAGMA integrity_check;

-- 3. インデックス再構築
REINDEX;
```

## 8. パフォーマンス監視

### **8.1 クエリ最適化**

```sql
-- 実行計画確認
EXPLAIN QUERY PLAN 
SELECT * FROM knowledge 
WHERE category = 'マーケティング' 
AND created_at >= '2024-10-01' 
ORDER BY quality_score DESC 
LIMIT 20;

-- 統計情報確認
SELECT name, sql FROM sqlite_master WHERE type = 'index';
```

### **8.2 監視クエリ**

```sql
-- テーブルサイズ確認
SELECT 
    name,
    (SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = t.name) as row_count
FROM sqlite_master t 
WHERE type = 'table';

-- インデックス使用状況
SELECT * FROM sqlite_stat1;

-- データベースサイズ
SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();
```

## 9. 拡張性考慮

### **9.1 将来の拡張**

```sql
-- ユーザー管理テーブル（将来機能）
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 権限管理テーブル（将来機能）
CREATE TABLE permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    knowledge_id INTEGER,
    permission_type TEXT CHECK(permission_type IN ('read', 'write', 'admin')),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (knowledge_id) REFERENCES knowledge(id)
);
```

### **9.2 マイグレーション戦略**

```sql
-- バージョン管理テーブル
CREATE TABLE schema_version (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- 初期バージョン
INSERT INTO schema_version (version, description) VALUES (1, 'Initial schema');
```

---

**文書情報**
- 作成日: 2024年10月20日
- バージョン: 1.0
- 最終更新: 2024年10月20日
- 作成者: システム開発チーム
