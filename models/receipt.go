package models

import (
	"mime/multipart"
	"time"
)

type Receipt struct {
	ID         int64     `db:"id" json:"id"`
	CategoryID int64     `db:"category_id" json:"category_id"`
	Date       string    `db:"date" json:"date"`
	ImageName  string    `db:"image_name" json:"image_name"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type NewReceiptFormFields struct {
	CategoryID string
	Date       string
	File       string
}

type NewReceipt struct {
	CategoryID int64
	Date       string
	File       multipart.File
	FileName   string
	NewReceiptFormFields
}

var NewReceiptFormFieldsInstance = NewReceiptFormFields{
	CategoryID: "cateogry_id",
	Date:       "date",
	File:       "file",
}