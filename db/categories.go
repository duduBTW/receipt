package db

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	LucideIconName string    `db:"lucide_icon_name"`
	Hue            int       `db:"hue"`
	Saturation     int       `db:"saturation"`
	Lightness      int       `db:"lightness"`
	CreatedAt      time.Time `db:"created_at"`
}

type CategoryStore interface {
	CreateCategory(ctx context.Context, category *Category) error
	GetCategory(ctx context.Context, id int64) (*Category, error)
	ListCategories(ctx context.Context) ([]Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type SQLiteCategoryStore struct {
	db *sql.DB
}

func NewSQLiteCategoryStore(db *sql.DB) *SQLiteCategoryStore {
	return &SQLiteCategoryStore{db: db}
}

func (s *SQLiteCategoryStore) CreateCategory(ctx context.Context, category *Category) error {
	query := `
        INSERT INTO categories (name, lucide_icon_name, hue, saturation, lightness)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id, created_at`

	return s.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.LucideIconName,
		category.Hue,
		category.Saturation,
		category.Lightness,
	).Scan(&category.ID, &category.CreatedAt)
}

func (s *SQLiteCategoryStore) GetCategory(ctx context.Context, id int64) (*Category, error) {
	category := &Category{}
	query := `SELECT id, name, lucide_icon_name, hue, saturation, lightness, created_at FROM categories WHERE id = ?`

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.LucideIconName,
		&category.Hue,
		&category.Saturation,
		&category.Lightness,
		&category.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func (s *SQLiteCategoryStore) ListCategories(ctx context.Context) ([]Category, error) {
	query := `SELECT id, name, lucide_icon_name, hue, saturation, lightness FROM categories`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.LucideIconName,
			&category.Hue,
			&category.Saturation,
			&category.Lightness,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (s *SQLiteCategoryStore) UpdateCategory(ctx context.Context, category *Category) error {
	query := `
        UPDATE categories 
        SET name = ?, lucide_icon_name = ?, hue = ?, saturation = ?, lightness = ?
        WHERE id = ?`

	result, err := s.db.ExecContext(
		ctx,
		query,
		category.Name,
		category.LucideIconName,
		category.Hue,
		category.Saturation,
		category.Lightness,
		category.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *SQLiteCategoryStore) DeleteCategory(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = ?`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
