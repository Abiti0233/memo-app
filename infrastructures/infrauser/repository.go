package infrauser

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/Abiti0233/memo-app/domains"
)

var ErrNoRows = errors.New("no rows in result set")

type UserRepository interface {
	Create(user *domains.User) error
	GetByID(id string) (*domains.User, error)
	GetByEmail(email string) (*domains.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domains.User) error {
	user.ID = uuid.New().String()
	query := `INSERT INTO Users (id, name, email, emailLowerCase, createdAt, updatedAt)
						VALUES ($1, $2, $3, LOWER($3), NOW(), NOW())`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email)
	return err
}

func (r *userRepository) GetByID(id string) (*domains.User, error) {
	query := `SELECT id, name, email, emailLowerCase, emailVerified, createdAt, updatedAt
						FROM Users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user domains.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EmailLowerCase, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domains.User, error) {
	query := `SELECT id, name, email, emailLowerCase, emailVerified, createdAt, updatedAt
						FROM Users WHERE emailLowerCase = LOWER($1)`
	row := r.db.QueryRow(query, email)

	var user domains.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EmailLowerCase, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}