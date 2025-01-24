package infrabookmark

import (
	"database/sql"
	"errors"
	"os/user"

	"github.com/Abiti0233/memo-app/domains"
	domains "github.com/Abiti0233/memo-app/domains/bookmark"
	"github.com/google/uuid"
)

var ErrNoRows = errors.New("no rows in result set")

type BookmarkRepository interface {
	Create(bookmark *domains.Bookmark) error
	Delete(id string) error
	GetByID(id string) (*domains.Bookmark, error)
	ListByUser(userID string) ([]domains.Bookmark, error)
	GetByUserAndMemo(userID, memoID string) (*domains.Bookmark, error)
}

type bookmarkRepository struct {
	// dbはデータベースへの接続を表す。
	db *sql.DB
}

func NewBookmarkRepository(db *sql.DB) BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) Create(bookmark *domains.Bookmark) error {
	bookmark.ID  = uuid.New().String()
	query := `INSERT INTO Bookmarks (id, userId, memoId, createdAt)
						VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(query, bookmark.ID, bookmark.UserID, bookmark.MemoID)
	return err
}

func (r *bookmarkRepository) Delete(id string) error {
	query := `DELETE FROM Bookmarks WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
			return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
			return err
	}
	if rowsAffected == 0 {
			return ErrNoRows
	}
	return nil
}

func (r *bookmarkRepository) GetByID(id string) (*domains.Bookmark, error) {
	query := `SELECT id, userId, memoId, createdAt
						FROM Bookmarks WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var bookmark domains.Bookmark
	err := row.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.MemoID, &bookmark.CreatedAt)
	if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
					return nil, nil
			}
			return nil, err
	}

	// Bookmark構造体のポインタ（アドレス？）を返す。
	// この関数を呼び出す側で、Bookmark構造体の値を変更できるようにするため。
	return &bookmark, nil
}

func (r *bookmarkRepository) ListByUser(userID string) ([]domains.Bookmark, error) {
	// ブックマークは更新される＝データから消えるため作成日時の降順でも最新状態を表せるため、
	// 作成日時の降順で取得するクエリを実行する。
	query := `SELECT id, userId, memoId, createdAt FROM Bookmarks WHERE userId = $1 ORDER BY createdAt DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []domains.Bookmark
	for rows.Next() {
		var bookmark domains.Bookmark
		err := rows.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.MemoID, &bookmark.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	return bookmarks, nil
}

func (r *bookmarkRepository) GetByUserAndMemo(userID, memoID string) (*domains.Bookmark, error) {
	query := `SELECT id, userId, memoId, createdAt FROM Bookmarks WHERE userId = $1 AND memoId = $2`
	row := r.db.QueryRow(query, userID, memoID)

	var bookmark domains.Bookmark
	err := row.Scan(&bookmark.ID, &bookmark.UserID, &bookmark.MemoID, &bookmark.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &bookmark, nil
}

