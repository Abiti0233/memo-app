package usecases

import (
	"errors"

	"github.com/Abiti0233/memo-app/domains/memo"
	"github.com/Abiti0233/memo-app/infrastructures"
	"github.com/Abiti0233/memo-app/infrastructures/infracategory"
	"github.com/Abiti0233/memo-app/infrastructures/inframemo"
)

var (
	ErrMemoNotFound = errors.New("memo not found")
)

type MemoUseCase interface {
	CreateMemo(userID string, memo *memo.Memo) error
	UpdateMemo(userID string, memo *memo.Memo) error
	DeleteMemo(userID, id string) error
	GetMemoByID(userID, id string) (*memo.Memo, error)
	ListMemos(userID string) ([]memo.Memo, error)
	ArchiveMemo(userID, id string, isArchived bool) error
	AssignCategory(userID, memoID, categoryID string) error
	RemoveCategory(userID, memoID, categoryID string) error
}

type memoUseCase struct {
	memoRepo     inframemo.MemoRepository
	categoryRepo infracategory.CategoryRepository
}

func NewMemoUseCase(memoRepo inframemo.MemoRepository, categoryRepo infracategory.CategoryRepository) MemoUseCase {
	return &memoUseCase{memoRepo: memoRepo, categoryRepo: categoryRepo}
}

func (u *memoUseCase) CreateMemo(userID string, memo *memo.Memo) error {
	memo.UserID = userID
	return u.memoRepo.Create(memo)
}

func (u *memoUseCase) UpdateMemo(userID string, memo *memo.Memo) error {
	memo.UserID = userID
	// メモの更新処理をUpdateメソッドを呼び出して行う。
	// 成功時はnilが返ってくる。
	err := u.memoRepo.Update(memo)
	if err == infrastructures.ErrNoRows {
		return ErrMemoNotFound
	}
	return err
}

func (u *memoUseCase) DeleteMemo(userID, id string) error {
	err := u.memoRepo.Delete(id)
	if err == infrastructures.ErrNoRows {
		return ErrMemoNotFound
	}
	return err
}

func (u *memoUseCase) GetMemoByID(userID, id string) (*memo.Memo, error) {
	memo, err := u.memoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	// メモがないもしくは、ユーザーIDが異なる場合は、NotFoundエラーを返す。
	if memo == nil || memo.UserID != userID {
		return nil, ErrMemoNotFound
	}
	return memo, nil
}

func (u *memoUseCase) ListMemos(userID string) ([]memo.Memo, error) {
	return u.memoRepo.ListByUser(userID)
}

func (u *memoUseCase) ArchiveMemo(userID, id string, isArchived bool) error {
	memo, err := u.memoRepo.GetByID(id)
	if err != nil {
		return err
	}
	if memo == nil || memo.UserID != userID {
		return ErrMemoNotFound
	}
	return u.memoRepo.Archive(id, isArchived)
}

func (u *memoUseCase) AssignCategory(userID, memoID, categoryID string) error {
	// メモが存在するか確認する。
	memo, err := u.memoRepo.GetByID(memoID)
	if err != nil {
		return err
	}
	if memo == nil || memo.UserID != userID {
		return ErrMemoNotFound
	}
	category, err := u.categoryRepo.GetByID(categoryID)
	if err != nil {
		return err
	}
	if category == nil || category.UserID != userID {
		return errors.New("カテゴリーが見つかりません。")
	}
	return u.memoRepo.AssignCategory(memoID, categoryID)
}

func (u *memoUseCase) RemoveCategory(userID, memoID, categoryID string) error {
	memo, err := u.memoRepo.GetByID(memoID)
	if err != nil {
		return err
	}
	if memo == nil || memo.UserID != userID {
		return ErrMemoNotFound
	}
	category, err := u.categoryRepo.GetByID(categoryID)
	if err != nil {
		return err
	}
	if category == nil || category.UserID != userID {
		return errors.New("カテゴリーが見つかりません。")
	}
	return u.memoRepo.RemoveCategory(memoID, categoryID)
}
