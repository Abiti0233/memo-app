package infrauser

import (
	"database/sql"
	"errors"

	"github.com/Abiti0233/memo-app/domains/user"
)

type UserRepository interface {
	Create(user *user.User) error
	GetByID(id string) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *user.User) error {
	query := `INSERT INTO Users (id, name, email, emailVerified, createdAt, updatedAt)
						VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.EmailVerified)
	return err
}

func (r *userRepository) GetByID(id string) (*user.User, error) {
	query := `SELECT id, name, email, emailLowerCase, emailVerified, createdAt, updatedAt
						FROM Users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user user.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EmailLowerCase, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*user.User, error) {
	query := `SELECT id, name, email, emailLowerCase, emailVerified, createdAt, updatedAt
						FROM Users WHERE emailLowerCase = LOWER($1)`
	row := r.db.QueryRow(query, email)

	var user user.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EmailLowerCase, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
