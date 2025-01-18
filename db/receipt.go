package db

import (
	"context"
	"database/sql"

	"github.com/dudubtw/receipt/models"
)

type ReceiptStore interface {
	CreateReceipt(ctx context.Context, receipt *models.Receipt) error
	GetReceipt(ctx context.Context, id int64) (*models.Receipt, error)
	ListReceipts(ctx context.Context) ([]models.Receipt, error)
	UpdateReceipt(ctx context.Context, receipt *models.Receipt) error
	DeleteReceipt(ctx context.Context, id int64) error
}

type SQLiteReceiptStore struct {
	db *sql.DB
}

func NewSQLiteReceiptStore(db *sql.DB) *SQLiteReceiptStore {
	return &SQLiteReceiptStore{db: db}
}

func (s *SQLiteReceiptStore) CreateReceipt(ctx context.Context, receipt *models.Receipt) error {
	query := `
        INSERT INTO receipts (category_id, date, image_name, created_at)
        VALUES (?, ?, ?, ?)
        RETURNING id, created_at`

	return s.db.QueryRowContext(
		ctx,
		query,
		receipt.CategoryID,
		receipt.Date,
		receipt.ImageName,
		receipt.CreatedAt,
	).Scan(&receipt.ID, &receipt.CreatedAt)
}

func (s *SQLiteReceiptStore) GetReceipt(ctx context.Context, id int64) (*models.Receipt, error) {
	receipt := &models.Receipt{}
	query := `SELECT id, category_id, date, image_name, created_at FROM receipts WHERE id = ?`

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&receipt.ID,
		&receipt.CategoryID,
		&receipt.Date,
		&receipt.ImageName,
		&receipt.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return receipt, err
}

func (s *SQLiteReceiptStore) ListReceipts(ctx context.Context) ([]models.Receipt, error) {
	query := `SELECT id, category_id, date, image_name, created_at FROM receipts`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		if err := rows.Scan(
			&receipt.ID,
			&receipt.CategoryID,
			&receipt.Date,
			&receipt.ImageName,
			&receipt.CreatedAt,
		); err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, rows.Err()
}

func (s *SQLiteReceiptStore) UpdateReceipt(ctx context.Context, receipt *models.Receipt) error {
	query := `
        UPDATE receipts 
        SET category_id = ?, date = ?, image_name = ?
        WHERE id = ?`

	result, err := s.db.ExecContext(
		ctx,
		query,
		receipt.CategoryID,
		receipt.Date,
		receipt.ImageName,
		receipt.ID,
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

func (s *SQLiteReceiptStore) DeleteReceipt(ctx context.Context, id int64) error {
	query := `DELETE FROM receipts WHERE id = ?`

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
