// ユーザーに関するビジネスロジックを定義
package usecases

import (
	"errors"

	"github.com/Abiti0233/memo-app/domains"
	"github.com/Abiti0233/memo-app/infrastructures"
)

var (
	ErrUserNotFound = errors.New("ユーザーが見つかりません。")
	ErrUserAlreadyExists = errors.New("ユーザーは既に存在します。")
)

type UserUseCase interface {
	RegisterUser(user *domains.User) error
	GetUserByID(id string) (*domains.User, error)
	GetUserByEmail(email string) (*domains.User, error)
}

type userUseCase struct {
	userRepo infrastructures.UserRepository
}

func NewUserUseCase(userRepo infrastructures.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) RegisterUser(user *domains.User) error {
	// メールアドレスが既に登録されているかを確認する。
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// 登録されていなかったら、ユーザーを登録する。
	return u.userRepo.CreateUser(user)
}

func (u *userUseCase) GetUserByID (id string) (*domains.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userUseCase) GetUserByEmail(email string) (*domains.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}