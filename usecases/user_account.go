// ユーザーに関するビジネスロジックを定義
package usecases

import (
	"errors"
	"log"

	"github.com/Abiti0233/memo-app/domains/user"
	"github.com/Abiti0233/memo-app/infrastructures/infrauser"
)

var (
	ErrUserNotFound      = errors.New("ユーザーが見つかりません。")
	ErrUserAlreadyExists = errors.New("ユーザーは既に存在します。")
)

type UserUseCase interface {
	RegisterUser(user *user.User) error
	GetUserByID(id string) (*user.User, error)
	GetUserByEmail(email string) (*user.User, error)
}

type userUseCase struct {
	userRepo infrauser.UserRepository
}

func NewUserUseCase(userRepo infrauser.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) RegisterUser(user *user.User) error {
	// メールアドレスが既に登録されているかを確認する。
	existingUser, err := u.userRepo.GetByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// 登録されていなかったら、ユーザーを登録する。
	return u.userRepo.Create(user)
}

func (u *userUseCase) GetUserByID(id string) (*user.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		log.Printf("userRepo.GetByID()がエラーを返しました: %v\n", err)
		return nil, err
	}
	if user == nil {
		log.Printf("userがnilです: %v\n", user)
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userUseCase) GetUserByEmail(email string) (*user.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
