-- 1. Users テーブル
CREATE TABLE Users (
  id VARCHAR(64) PRIMARY KEY, -- ユーザーID（Google ID）
  name VARCHAR(255), -- ユーザー名
  email VARCHAR(512) NOT NULL UNIQUE, -- ログインに使用する Email アドレス
  emailLowerCase VARCHAR(512) GENERATED ALWAYS AS (LOWER(email)) STORED NOT NULL, -- 小文字化したログインに使用する Email アドレス
  emailVerified TIMESTAMP, -- Email確認日時
  createdAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 作成日時
  updatedAt TIMESTAMP NOT NULL DEFAULT NOW() -- 更新日時
);

-- インデックスの作成
CREATE UNIQUE INDEX UsersEmailLowerCase ON Users(emailLowerCase);

-- 2. Categories テーブル
CREATE TABLE Categories (
  id VARCHAR(64) PRIMARY KEY, -- カテゴリID
  userId VARCHAR(64) NOT NULL, -- ユーザーID
  name VARCHAR(255) NOT NULL, -- カテゴリ名
  createdAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 作成日時
  updatedAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 更新日時
  FOREIGN KEY (userId) REFERENCES Users(id) ON DELETE CASCADE
);

-- ユニーク制約
CREATE UNIQUE INDEX CategoriesUserCategoryNameUnique ON Categories(userId, name);

-- 3. Memos テーブル
CREATE TABLE Memos (
  id VARCHAR(64) PRIMARY KEY, -- メモID
  userId VARCHAR(64) NOT NULL, -- ユーザーID
  title VARCHAR(255) NOT NULL, -- メモタイトル
  content TEXT NOT NULL, -- メモ内容
  is_archived BOOLEAN NOT NULL DEFAULT FALSE, -- アーカイブ状態
  createdAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 作成日時
  updatedAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 更新日時
  FOREIGN KEY (userId) REFERENCES Users(id) ON DELETE CASCADE
);

-- 4. MemoCategories テーブル (中間テーブル)
CREATE TABLE MemoCategories (
  memoId VARCHAR(64) NOT NULL, -- メモID
  categoryId VARCHAR(64) NOT NULL, -- カテゴリID
  PRIMARY KEY (memoId, categoryId),
  FOREIGN KEY (memoId) REFERENCES Memos(id) ON DELETE CASCADE,
  FOREIGN KEY (categoryId) REFERENCES Categories(id) ON DELETE CASCADE
);

-- インデックスの作成
CREATE INDEX MemoCategoriesMemoCategoryIndex ON MemoCategories(memoId, categoryId);

-- 5. Bookmarks テーブル
CREATE TABLE Bookmarks (
  id VARCHAR(64) PRIMARY KEY, -- ブックマークID
  userId VARCHAR(64) NOT NULL, -- ユーザーID
  memoId VARCHAR(64) NOT NULL, -- メモID
  createdAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 作成日時
  FOREIGN KEY (userId) REFERENCES Users(id) ON DELETE CASCADE,
  FOREIGN KEY (memoId) REFERENCES Memos(id) ON DELETE CASCADE
);

-- インデックスの作成
CREATE UNIQUE INDEX BookmarksUserMemoIndex ON Bookmarks(userId, memoId) WHERE userId IS NOT NULL AND memoId IS NOT NULL;

-- 6. UserSettings テーブル
CREATE TABLE UserSettings (
  userId VARCHAR(64) PRIMARY KEY, -- ユーザーID
  theme VARCHAR(64) NOT NULL DEFAULT 'light', -- テーマ設定
  updatedAt TIMESTAMP NOT NULL DEFAULT NOW(), -- 更新日時
  FOREIGN KEY (userId) REFERENCES Users(id) ON DELETE CASCADE
);
