package infracategory

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/Abiti0233/memo-app/domains"
)

var ErrNoRows = errors.New("no rows in result set")

type CategoryRepository interface {
	Create(category *domains.Category) error
	Update(category *domains.Category) error
	Delete(id string) error
	GetByID(id string) (*domains.Category, error)
	ListByUser(userID string) ([]domains.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domains.Category) error {
	category.ID = uuid.New().String()
	query := `INSERT INTO Categories (id, userId, name, createdAt, updatedAt)
						VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.Exec(query, category.ID, category.UserID, category.Name)
	return err
}

func (r *categoryRepository) Update(category *domains.Category) error {
	query := `UPDATE Categories SET name = $1, updatedAt = NOW()
						WHERE id = $2 AND userId = $3`
	res, err := r.db.Exec(query, category.Name, category.ID, category.UserID)
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

func (r *categoryRepository) Delete(id string) error {
	query := `DELETE FROM Categories WHERE id = $1`
	res, err := r.db.Exec(query, id)

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRows
	}
	return nil
}

func (r *categoryRepository) GetByID(id string) (*domains.Category, error) {
	query := `SELECT id, userId, name, createdAt, updatedAt
						FROM Categories WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var category domains.Category
	err := row.Scan(&category.ID, &category.UserID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) ListByUser(userID string) ([]domains.Category, error) {
	// カテゴリーは、作成日付で降順にする
	query := `SELECT id, userId, name, createdAt, updatedAt FROM Categories WHERE userId = $1 ORDER BY createdAt DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var categories []domains.Category
	for rows.Next() {
			var category domains.Category
			err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
			if err != nil {
					return nil, err
			}
			categories = append(categories, category)
	}
	return categories, nil
}

