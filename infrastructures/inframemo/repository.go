package inframemo

import (
	"database/sql"
	"errors"

	"github.com/Abiti0233/memo-app/domains/memo"
	"github.com/Abiti0233/memo-app/infrastructures"
	"github.com/google/uuid"
)

type MemoRepository interface {
	Create(memo *memo.Memo) error
	Update(memo *memo.Memo) error
	Delete(id string) error
	GetByID(id string) (*memo.Memo, error)
	ListByUser(userID string) ([]memo.Memo, error)
	Archive(id string, isArchived bool) error
	AssignCategory(memoID, categoryID string) error
	RemoveCategory(memoID, categoryID string) error
}

type memoRepository struct {
	db *sql.DB
}

func NewMemoRepository(db *sql.DB) MemoRepository {
	return &memoRepository{db: db}
}

func (r *memoRepository) Create(memo *memo.Memo) error {
	memo.ID = uuid.New().String()
	query := `INSERT INTO Memos (id, userId, title, content, is_archived, createdAt, updatedAt)
						VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`
	_, err := r.db.Exec(query, memo.ID, memo.UserID, memo.Title, memo.Content, memo.IsArchived)
	return err
}

func (r *memoRepository) Update(memo *memo.Memo) error {
	query := `UPDATE Memos SET title = $1, content = $2, is_archived = $3, updatedAt = NOW()
						WHERE id = $4 AND userId = $5`
	res, err := r.db.Exec(query, memo.Title, memo.Content, memo.IsArchived, memo.ID, memo.UserID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return infrastructures.ErrNoRows
	}
	return nil
}

func (r *memoRepository) Delete(id string) error {
	query := `DELETE FROM Memos WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return infrastructures.ErrNoRows
	}
	return nil
}

func (r *memoRepository) GetByID(id string) (*memo.Memo, error) {
	query := `SELECT id, userId, title, content, is_archived, createdAt, updatedAt
						FROM Memos WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var memo memo.Memo
	err := row.Scan(&memo.ID, &memo.UserID, &memo.Title, &memo.Content, &memo.IsArchived,
		&memo.CreatedAt, &memo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &memo, nil
}

func (r *memoRepository) ListByUser(userID string) ([]memo.Memo, error) {
	// Memosテーブルから、指定したユーザーIDのメモを取得するクエリ。
	// メモ一覧は作成された日付ではなくて更新された日付で降順にする
	query := `SELECT id, userId, title, content, is_archived, createdAt, updatedAt
						FROM Memos WHERE userId = $1 ORDER BY updatedAt DESC`

	// r.db.Query() は、クエリを実行して結果を取得するもの。第一引数にクエリ、第二引数以降にクエリのプレースホルダに入れる値を指定する。
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []memo.Memo
	// rows.Next() は次の行があれば true を返し、次の行にカーソルを移動するもの。なければfor文を抜ける。
	for rows.Next() {
		var memo memo.Memo
		// rows.Scan() はカーソルの位置の行からデータを取得して、引数に渡した変数に格納するもの。（内の変数の数と型は、SELECT文で取得するカラムの数と型と一致している必要がある。）
		err := rows.Scan(&memo.ID, &memo.UserID, &memo.Title, &memo.Content, &memo.IsArchived, &memo.CreatedAt, &memo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

func (r *memoRepository) Archive(id string, isArchived bool) error {
	// Memosテーブルのis_archivedカラムを更新するクエリ。
	query := `UPDATE Memos SET is_archived = $1, updatedAt = NOW()
						WHERE id = $2`
	res, err := r.db.Exec(query, isArchived, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil
	}
	if rowsAffected == 0 {
		return infrastructures.ErrNoRows
	}
	return nil
}

func (r *memoRepository) AssignCategory(memoID, categoryID string) error {
	// ON CONFLICT DO NOTHING は、重複エラーが発生した場合に何もしないようにするためのもの。
	query := `INSERT INTO MemoCategories (memoId, categoryId) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	// r.db.Exec() は、クエリを実行するもの。第一引数にクエリ、第二引数以降にクエリのプレースホルダに入れる値を指定する。
	// $1, $2 は、プレースホルダ。$1 には memoID が、$2 には categoryID が入る想定。
	_, err := r.db.Exec(query, memoID, categoryID)
	return err
}

func (r *memoRepository) RemoveCategory(memoID, categoryID string) error {
	query := `DELETE FROM MemoCategories WHERE memoId = &1 AND cateogyId = $2`
	res, err := r.db.Exec(query, memoID, categoryID)
	if err != nil {
		return err
	}
	// res.RowsAffected() は、クエリの実行結果の行数を取得するもの。
	// このメソッドの場合は、削除された行の数を確認する。
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return infrastructures.ErrNoRows
	}
	return nil
}
