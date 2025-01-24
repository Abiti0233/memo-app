package usecases

import (
	"errors"

	"github.com/Abiti0233/memo-app/domains"
	"github.com/Abiti0233/memo-app/infrastructures"
)

var (
	ErrBookmarkNotFound = errors.New("ブックマークが見つかりません。")
	ErrBookmarkAlreadyExists = errors.New("ブックマークは既に存在します。")
)

type BookmarkUseCase interface {
	CreateBookmark(userID string, memoID string) (*domains.Bookmark, error)
	DeleteBookmark(userID, id string) error
	GetBookmarkByID(userID, id string) (*domains.Bookmark, error)
	ListBookmarks(userID string) ([]domains.Bookmark, error)
}

type bookmarkUseCase struct {
	bookmarkRepo infrastructures.BookmarkRepository
}

func NewBookmarkUseCase(bookmarkRepo infrastructures.BookmarkRepository) BookmarkUseCase {
	return &bookmarkUseCase{bookmarkRepo: bookmarkRepo}
}

func (u *bookmarkUseCase) CreateBookmark(userID string, memoID string) (*domains.Bookmark, error) {
	// ブックマークが既に存在するかを確認する。
	existing, err := u.bookmarkRepo.GetByUserAndMemo(userID, memoID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrBookmarkAlreadyExists
	}

	// ブックマークが存在しなかったら、ブックマークを作成する。
	bookmark := &domains.Bookmark{
		UserID: userID,
		MemoID: memoID,
	}
	err = u.bookmarkRepo.Create(bookmark)
	if err != nil {
		return nil, err
	}
	return bookmark, nil
}

func (u *bookmarkUseCase) DeleteBookmark(userID, id string) error {
	// ブックマークが存在するかを確認する。
	bookmark, err := u.bookmarkRepo.GetByID(id)
	if err != nil {
		return err
	}
	// ブックマークが存在しない場合、またはユーザーが異なる場合はエラーを返す。
	if bookmark == nil || bookmark.UserID != userID {
		return ErrBookmarkNotFound
	}

	return u.bookmarkRepo.Delete(id)
}

func (u *bookmarkUseCase) GetBookmarkByID(userID, id string) (*domains.Bookmark, error) {
	bookmark, err := u.bookmarkRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if bookmark == nil || bookmark.UserID != userID {
		return nil, ErrBookmarkNotFound
	}
	return bookmark, nil
}

func (u *bookmarkUseCase) ListBookmarks(userID string) ([]domains.Bookmark, error) {
	return u.bookmarkRepo.ListByUser(userID)
}