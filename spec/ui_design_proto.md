# ナレッジ管理ツール - UI/UX設計書

## 1. デザイン概要

### **1.1 デザインコンセプト**

```yaml
コンセプト: "シンプル・直感的・効率的"
- シンプル: 余計な装飾を排除し、機能に集中
- 直感的: 説明なしで操作できるUI
- 効率的: 最小のクリックで目的を達成
```

### **1.2 デザイン原則**

- **一貫性**: 全画面で統一されたデザイン言語
- **アクセシビリティ**: 多様なユーザーに対応
- **レスポンシブ**: デスクトップ中心、将来のモバイル対応
- **パフォーマンス**: 軽量で高速な表示

## 2. 画面構成

### **2.1 画面一覧**

```yaml
メイン画面:
  - ダッシュボード: ナレッジ一覧・統計情報
  - ナレッジ作成: タネ投入・AI処理・編集
  - ナレッジ詳細: 個別ナレッジの表示・編集
  - 検索結果: 検索・フィルタ結果表示
  - 設定: システム設定・API設定

共通要素:
  - ヘッダー: ナビゲーション・検索
  - サイドバー: カテゴリ・フィルタ
  - フッター: システム情報
```

### **2.2 ナビゲーション構造**

```
ダッシュボード
├── ナレッジ一覧
├── 統計情報
└── 最近の活動

ナレッジ管理
├── 新規作成
├── 検索・フィルタ
└── カテゴリ別表示

設定
├── システム設定
├── API設定
└── バックアップ
```

## 3. 画面設計

### **3.1 ダッシュボード画面**

#### **レイアウト**
```
┌─────────────────────────────────────────────────────────┐
│ ヘッダー [ロゴ] [検索] [新規作成] [設定]                │
├─────────────────────────────────────────────────────────┤
│ サイドバー │ メインコンテンツ                            │
│           │ ┌─────────────────────────────────────────┐ │
│ カテゴリ  │ │ 統計情報                                │ │
│ - 全件    │ │ 総数: 150件 今月: 12件 品質平均: 78点   │ │
│ - マーケ  │ └─────────────────────────────────────────┘ │
│ - 開発    │ ┌─────────────────────────────────────────┐ │
│ - 営業    │ │ 最近のナレッジ                          │ │
│ - 改善    │ │ [カード1] [カード2] [カード3] ...      │ │
│           │ └─────────────────────────────────────────┘ │
│ フィルタ  │ ┌─────────────────────────────────────────┐ │
│ - 品質    │ │ 人気のナレッジ                          │ │
│ - 日付    │ │ [カード1] [カード2] [カード3] ...      │ │
│ - タグ    │ └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

#### **コンポーネント**

**統計情報カード**
```html
<div class="stats-card">
  <div class="stat-item">
    <span class="stat-number">150</span>
    <span class="stat-label">総ナレッジ数</span>
  </div>
  <div class="stat-item">
    <span class="stat-number">12</span>
    <span class="stat-label">今月作成</span>
  </div>
  <div class="stat-item">
    <span class="stat-number">78</span>
    <span class="stat-label">品質平均</span>
  </div>
</div>
```

**ナレッジカード**
```html
<div class="knowledge-card">
  <div class="card-header">
    <h3 class="card-title">SNS広告CTR改善事例</h3>
    <span class="card-category">マーケティング</span>
  </div>
  <div class="card-summary">
    ターゲティングを絞り込むことでCTRが0.5%から1.2%に改善...
  </div>
  <div class="card-meta">
    <span class="quality-score">品質: 85点</span>
    <span class="created-date">2024/10/20</span>
  </div>
  <div class="card-actions">
    <button class="btn-view">詳細</button>
    <button class="btn-export">エクスポート</button>
  </div>
</div>
```

### **3.2 ナレッジ作成画面**

#### **レイアウト**
```
┌─────────────────────────────────────────────────────────┐
│ ヘッダー [ロゴ] [検索] [新規作成] [設定]                │
├─────────────────────────────────────────────────────────┤
│ メインコンテンツ                                        │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ナレッジ作成                                        │ │
│ │                                                     │ │
│ │ タネ投入                                            │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ [URL] [テキスト]                                │ │ │
│ │ │ ┌─────────────────────────────────────────────┐ │ │ │
│ │ │ │ URL: https://example.com/article            │ │ │ │
│ │ │ │ または                                      │ │ │ │
│ │ │ │ テキスト: [テキスト入力エリア]              │ │ │ │
│ │ │ └─────────────────────────────────────────────┘ │ │ │
│ │ │ [作成開始]                                     │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │                                                     │ │
│ │ AI処理中...                                        │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ 処理中... 30秒以内に完了予定                    │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │                                                     │ │
│ │ 生成結果（編集可能）                                │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ タイトル: [SNS広告CTR改善事例]                  │ │ │
│ │ │ 要約: [要約テキスト]                            │ │ │
│ │ │ 問題: [問題テキスト]                            │ │ │
│ │ │ 解決策: [解決策テキスト]                        │ │ │
│ │ │ 結果: [結果テキスト]                            │ │ │
│ │ │ キーワード: [マーケティング, SNS, CTR]          │ │ │
│ │ │ カテゴリ: [マーケティング ▼]                    │ │ │
│ │ │ [保存] [キャンセル]                             │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

#### **ステップフロー**

**ステップ1: タネ投入**
```html
<div class="seed-input">
  <div class="input-type-selector">
    <button class="type-btn active" data-type="url">URL</button>
    <button class="type-btn" data-type="text">テキスト</button>
  </div>
  
  <div class="input-area">
    <input type="url" class="url-input" placeholder="https://example.com/article">
    <textarea class="text-input" placeholder="ナレッジの元となるテキストを入力..."></textarea>
  </div>
  
  <button class="btn-primary" id="create-knowledge">作成開始</button>
</div>
```

**ステップ2: AI処理中**
```html
<div class="processing-status">
  <div class="spinner"></div>
  <p>AI処理中... 30秒以内に完了予定</p>
  <div class="progress-bar">
    <div class="progress-fill"></div>
  </div>
</div>
```

**ステップ3: 編集・保存**
```html
<div class="knowledge-editor">
  <div class="form-group">
    <label>タイトル</label>
    <input type="text" class="form-control" value="SNS広告CTR改善事例">
  </div>
  
  <div class="form-group">
    <label>要約</label>
    <textarea class="form-control" rows="3">ターゲティングを絞り込むことで...</textarea>
  </div>
  
  <div class="form-group">
    <label>問題</label>
    <textarea class="form-control" rows="4">SNS広告のCTRが低く...</textarea>
  </div>
  
  <div class="form-group">
    <label>解決策</label>
    <textarea class="form-control" rows="4">ターゲティングを絞り込み...</textarea>
  </div>
  
  <div class="form-group">
    <label>結果</label>
    <textarea class="form-control" rows="4">CTRが0.5%から1.2%に改善...</textarea>
  </div>
  
  <div class="form-group">
    <label>キーワード</label>
    <div class="keyword-tags">
      <span class="tag">マーケティング</span>
      <span class="tag">SNS</span>
      <span class="tag">CTR</span>
      <input type="text" class="tag-input" placeholder="キーワードを追加...">
    </div>
  </div>
  
  <div class="form-group">
    <label>カテゴリ</label>
    <select class="form-control">
      <option>マーケティング</option>
      <option>開発・技術</option>
      <option>営業・顧客対応</option>
      <option>業務改善・効率化</option>
      <option>その他</option>
    </select>
  </div>
  
  <div class="form-actions">
    <button class="btn-primary">保存</button>
    <button class="btn-secondary">キャンセル</button>
  </div>
</div>
```

### **3.3 ナレッジ詳細画面**

#### **レイアウト**
```
┌─────────────────────────────────────────────────────────┐
│ ヘッダー [ロゴ] [検索] [新規作成] [設定]                │
├─────────────────────────────────────────────────────────┤
│ メインコンテンツ                                        │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ナレッジ詳細                                        │ │
│ │                                                     │ │
│ │ タイトル: SNS広告CTR改善事例                        │ │
│ │ カテゴリ: マーケティング  品質: 85点  作成: 10/20  │ │
│ │                                                     │ │
│ │ 要約                                                │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ ターゲティングを絞り込むことでCTRが0.5%から1.2% │ │ │
│ │ │ に改善した。具体的な手法と効果をまとめた。      │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │                                                     │ │
│ │ 問題                                                │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ SNS広告のCTRが0.5%と低く、広告費の効率が悪い   │ │ │
│ │ │ 状態だった。ターゲティングが広すぎることが     │ │ │
│ │ │ 原因と推測された。                              │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │                                                     │ │
│ │ 解決策                                              │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ 1. 年齢層を25-35歳に絞り込み                    │ │ │
│ │ │ 2. 興味関心を「マーケティング」「SNS」に限定    │ │ │
│ │ │ 3. 地域を主要都市に限定                         │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │                                                     │ │
│ │ 結果                                                │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │ CTRが0.5%から1.2%に改善（2.4倍向上）            │ │ │
│ │ │ 広告費効率が大幅に改善                          │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │                                                     │ │
│ │ キーワード: [マーケティング] [SNS] [CTR] [改善]    │ │
│ │                                                     │ │
│ │ アクション                                          │ │
│ │ [編集] [エクスポート] [削除] [関連ナレッジ]        │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### **3.4 検索結果画面**

#### **レイアウト**
```
┌─────────────────────────────────────────────────────────┐
│ ヘッダー [ロゴ] [検索: CTR] [新規作成] [設定]           │
├─────────────────────────────────────────────────────────┤
│ サイドバー │ 検索結果                                    │
│           │ ┌─────────────────────────────────────────┐ │
│ フィルタ  │ │ 検索結果: "CTR" (15件)                  │ │
│ カテゴリ  │ │ [関連度順] [作成日順] [品質順]          │ │
│ - 全件    │ └─────────────────────────────────────────┘ │
│ - マーケ  │ ┌─────────────────────────────────────────┐ │
│ - 開発    │ │ [ナレッジカード1]                      │ │
│ - 営業    │ │ [ナレッジカード2]                      │ │
│ - 改善    │ │ [ナレッジカード3]                      │ │
│           │ │ ...                                    │ │
│ 品質      │ └─────────────────────────────────────────┘ │
│ - 90-100  │ [前へ] 1 2 3 4 5 [次へ]                    │
│ - 80-89   │                                             │
│ - 70-79   │                                             │
│ - 60-69   │                                             │
│ - 50-59   │                                             │
│           │                                             │
│ 日付      │                                             │
│ - 今日    │                                             │
│ - 今週    │                                             │
│ - 今月    │                                             │
│ - 3ヶ月   │                                             │
└─────────────────────────────────────────────────────────┘
```

## 4. コンポーネント設計

### **4.1 共通コンポーネント**

#### **ヘッダー**
```html
<header class="main-header">
  <div class="header-left">
    <h1 class="logo">KnowledgeFlow</h1>
  </div>
  <div class="header-center">
    <div class="search-box">
      <input type="text" placeholder="ナレッジを検索..." class="search-input">
      <button class="search-btn">🔍</button>
    </div>
  </div>
  <div class="header-right">
    <button class="btn-primary">新規作成</button>
    <button class="btn-secondary">設定</button>
  </div>
</header>
```

#### **サイドバー**
```html
<aside class="sidebar">
  <div class="sidebar-section">
    <h3>カテゴリ</h3>
    <ul class="category-list">
      <li><a href="#" class="category-link active">全件 (150)</a></li>
      <li><a href="#" class="category-link">マーケティング (45)</a></li>
      <li><a href="#" class="category-link">開発・技術 (38)</a></li>
      <li><a href="#" class="category-link">営業・顧客対応 (32)</a></li>
      <li><a href="#" class="category-link">業務改善・効率化 (25)</a></li>
      <li><a href="#" class="category-link">その他 (10)</a></li>
    </ul>
  </div>
  
  <div class="sidebar-section">
    <h3>フィルタ</h3>
    <div class="filter-group">
      <label>品質スコア</label>
      <input type="range" min="0" max="100" value="50" class="quality-slider">
      <span class="quality-value">50点以上</span>
    </div>
    
    <div class="filter-group">
      <label>作成日</label>
      <select class="date-filter">
        <option>すべて</option>
        <option>今日</option>
        <option>今週</option>
        <option>今月</option>
        <option>3ヶ月</option>
      </select>
    </div>
  </div>
</aside>
```

#### **ナレッジカード**
```html
<div class="knowledge-card">
  <div class="card-header">
    <h3 class="card-title">SNS広告CTR改善事例</h3>
    <span class="card-category">マーケティング</span>
  </div>
  
  <div class="card-summary">
    ターゲティングを絞り込むことでCTRが0.5%から1.2%に改善...
  </div>
  
  <div class="card-meta">
    <span class="quality-score">
      <span class="score-icon">⭐</span>
      品質: 85点
    </span>
    <span class="created-date">2024/10/20</span>
    <span class="view-count">閲覧: 12回</span>
  </div>
  
  <div class="card-tags">
    <span class="tag">マーケティング</span>
    <span class="tag">SNS</span>
    <span class="tag">CTR</span>
  </div>
  
  <div class="card-actions">
    <button class="btn-view">詳細</button>
    <button class="btn-export">エクスポート</button>
    <button class="btn-edit">編集</button>
  </div>
</div>
```

### **4.2 フォームコンポーネント**

#### **入力フィールド**
```html
<div class="form-group">
  <label for="title" class="form-label">タイトル</label>
  <input type="text" id="title" class="form-control" placeholder="ナレッジのタイトルを入力...">
  <div class="form-help">255文字以内で入力してください</div>
</div>

<div class="form-group">
  <label for="summary" class="form-label">要約</label>
  <textarea id="summary" class="form-control" rows="3" placeholder="ナレッジの要約を入力..."></textarea>
  <div class="form-help">1000文字以内で入力してください</div>
</div>
```

#### **キーワードタグ**
```html
<div class="keyword-input">
  <div class="keyword-tags">
    <span class="tag removable">
      マーケティング
      <button class="tag-remove">×</button>
    </span>
    <span class="tag removable">
      SNS
      <button class="tag-remove">×</button>
    </span>
    <input type="text" class="tag-input" placeholder="キーワードを追加...">
  </div>
</div>
```

## 5. スタイル設計

### **5.1 カラーパレット**

```css
:root {
  /* プライマリカラー */
  --primary-color: #2563eb;
  --primary-hover: #1d4ed8;
  --primary-light: #dbeafe;
  
  /* セカンダリカラー */
  --secondary-color: #64748b;
  --secondary-hover: #475569;
  --secondary-light: #f1f5f9;
  
  /* アクセントカラー */
  --accent-color: #059669;
  --accent-hover: #047857;
  --accent-light: #d1fae5;
  
  /* 警告カラー */
  --warning-color: #d97706;
  --warning-hover: #b45309;
  --warning-light: #fef3c7;
  
  /* エラーカラー */
  --error-color: #dc2626;
  --error-hover: #b91c1c;
  --error-light: #fee2e2;
  
  /* ニュートラルカラー */
  --gray-50: #f9fafb;
  --gray-100: #f3f4f6;
  --gray-200: #e5e7eb;
  --gray-300: #d1d5db;
  --gray-400: #9ca3af;
  --gray-500: #6b7280;
  --gray-600: #4b5563;
  --gray-700: #374151;
  --gray-800: #1f2937;
  --gray-900: #111827;
  
  /* 背景カラー */
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;
  
  /* テキストカラー */
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --text-tertiary: #9ca3af;
}
```

### **5.2 タイポグラフィ**

```css
/* フォント設定 */
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  color: var(--text-primary);
}

/* 見出し */
h1 { font-size: 2.25rem; font-weight: 700; line-height: 1.2; }
h2 { font-size: 1.875rem; font-weight: 600; line-height: 1.3; }
h3 { font-size: 1.5rem; font-weight: 600; line-height: 1.4; }
h4 { font-size: 1.25rem; font-weight: 500; line-height: 1.4; }
h5 { font-size: 1.125rem; font-weight: 500; line-height: 1.4; }
h6 { font-size: 1rem; font-weight: 500; line-height: 1.4; }

/* 本文 */
p { margin-bottom: 1rem; }
small { font-size: 0.875rem; color: var(--text-secondary); }
```

### **5.3 レイアウト**

```css
/* グリッドシステム */
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.grid {
  display: grid;
  gap: 1.5rem;
}

.grid-cols-1 { grid-template-columns: repeat(1, 1fr); }
.grid-cols-2 { grid-template-columns: repeat(2, 1fr); }
.grid-cols-3 { grid-template-columns: repeat(3, 1fr); }
.grid-cols-4 { grid-template-columns: repeat(4, 1fr); }

/* フレックスボックス */
.flex { display: flex; }
.flex-col { flex-direction: column; }
.items-center { align-items: center; }
.justify-between { justify-content: space-between; }
.justify-center { justify-content: center; }

/* スペーシング */
.p-1 { padding: 0.25rem; }
.p-2 { padding: 0.5rem; }
.p-3 { padding: 0.75rem; }
.p-4 { padding: 1rem; }
.p-6 { padding: 1.5rem; }
.p-8 { padding: 2rem; }

.m-1 { margin: 0.25rem; }
.m-2 { margin: 0.5rem; }
.m-3 { margin: 0.75rem; }
.m-4 { margin: 1rem; }
.m-6 { margin: 1.5rem; }
.m-8 { margin: 2rem; }
```

### **5.4 コンポーネントスタイル**

```css
/* ボタン */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: var(--primary-hover);
}

.btn-secondary {
  background-color: var(--secondary-color);
  color: white;
}

.btn-secondary:hover {
  background-color: var(--secondary-hover);
}

/* フォーム */
.form-control {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--gray-300);
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.form-label {
  display: block;
  margin-bottom: 0.25rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

/* カード */
.card {
  background-color: var(--bg-primary);
  border: 1px solid var(--gray-200);
  border-radius: 0.5rem;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card-header {
  margin-bottom: 1rem;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.card-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

/* タグ */
.tag {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.5rem;
  background-color: var(--primary-light);
  color: var(--primary-color);
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.tag-remove {
  margin-left: 0.25rem;
  background: none;
  border: none;
  color: var(--primary-color);
  cursor: pointer;
  font-size: 0.875rem;
}
```

## 6. レスポンシブデザイン

### **6.1 ブレークポイント**

```css
/* モバイル */
@media (max-width: 640px) {
  .container { padding: 0 0.5rem; }
  .grid-cols-2 { grid-template-columns: 1fr; }
  .grid-cols-3 { grid-template-columns: 1fr; }
  .grid-cols-4 { grid-template-columns: 1fr; }
  
  .sidebar { display: none; }
  .main-content { width: 100%; }
}

/* タブレット */
@media (min-width: 641px) and (max-width: 1024px) {
  .grid-cols-3 { grid-template-columns: repeat(2, 1fr); }
  .grid-cols-4 { grid-template-columns: repeat(2, 1fr); }
}

/* デスクトップ */
@media (min-width: 1025px) {
  .grid-cols-3 { grid-template-columns: repeat(3, 1fr); }
  .grid-cols-4 { grid-template-columns: repeat(4, 1fr); }
}
```

### **6.2 モバイル対応**

```css
/* モバイルナビゲーション */
.mobile-nav {
  display: none;
}

@media (max-width: 640px) {
  .mobile-nav {
    display: block;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: var(--bg-primary);
    border-top: 1px solid var(--gray-200);
    padding: 0.5rem;
  }
  
  .mobile-nav-items {
    display: flex;
    justify-content: space-around;
  }
  
  .mobile-nav-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0.5rem;
    text-decoration: none;
    color: var(--text-secondary);
    font-size: 0.75rem;
  }
  
  .mobile-nav-item.active {
    color: var(--primary-color);
  }
}
```

## 7. アニメーション・トランジション

### **7.1 基本アニメーション**

```css
/* フェードイン */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.fade-in {
  animation: fadeIn 0.3s ease-in-out;
}

/* スライドアップ */
@keyframes slideUp {
  from { 
    opacity: 0;
    transform: translateY(20px);
  }
  to { 
    opacity: 1;
    transform: translateY(0);
  }
}

.slide-up {
  animation: slideUp 0.3s ease-out;
}

/* ローディングスピナー */
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--gray-300);
  border-top: 2px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}
```

### **7.2 ホバーエフェクト**

```css
/* カードホバー */
.knowledge-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.knowledge-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* ボタンホバー */
.btn {
  transition: all 0.2s;
}

.btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
```

## 8. アクセシビリティ

### **8.1 キーボードナビゲーション**

```css
/* フォーカス表示 */
.btn:focus,
.form-control:focus,
.nav-link:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* スキップリンク */
.skip-link {
  position: absolute;
  top: -40px;
  left: 6px;
  background: var(--primary-color);
  color: white;
  padding: 8px;
  text-decoration: none;
  z-index: 1000;
}

.skip-link:focus {
  top: 6px;
}
```

### **8.2 スクリーンリーダー対応**

```html
<!-- セマンティックHTML -->
<main role="main">
  <section aria-labelledby="knowledge-list">
    <h2 id="knowledge-list">ナレッジ一覧</h2>
    <div role="list">
      <article role="listitem">
        <h3>SNS広告CTR改善事例</h3>
        <p>ターゲティングを絞り込むことで...</p>
      </article>
    </div>
  </section>
</main>

<!-- ARIA属性 -->
<button aria-expanded="false" aria-controls="filter-menu">
  フィルタ
</button>

<div id="filter-menu" aria-hidden="true">
  <!-- フィルタメニュー -->
</div>
```

---

**文書情報**
- 作成日: 2025年9月4日
- バージョン: 0.1
- 最終更新: 2025年9月4日
- 作成者: ich