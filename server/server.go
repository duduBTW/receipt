package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/db"
	"github.com/dudubtw/receipt/models"
	pages "github.com/dudubtw/receipt/renderer/pages"
)

var dbInstance *sql.DB

func HomeHandler() {
	http.HandleFunc(constants.HomeRoute, func(w http.ResponseWriter, r *http.Request) {
		categoriesDb := db.NewSQLiteCategoryStore(dbInstance)
		categories, err := categoriesDb.ListCategories(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		app := pages.App(pages.Home(pages.HomeProps{
			Categories: categories,
		}))
		app.Render(r.Context(), w)
	})
}

func ReceiptHandler() {
	http.HandleFunc(constants.ReceiptRoute, func(w http.ResponseWriter, r *http.Request) {
		categoryID := r.URL.Query().Get(constants.ReceiptSearchParamCategory)
		categoryIDInt, _ := strconv.ParseInt(categoryID, 10, 64)
		receiptDb := db.NewSQLiteReceiptStore(dbInstance)
		gruppedReceipt, err := receiptDb.ListReceiptsByDate(context.Background(), categoryIDInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categoriesDb := db.NewSQLiteCategoryStore(dbInstance)
		categories, _ := categoriesDb.ListCategories(context.Background())

		app := pages.App(pages.ReceiptsPage(pages.ReceiptsPageProps{
			GruppedReceipt:          gruppedReceipt,
			Categories:              categories,
			DefaultCategorySelected: categoryIDInt,
		}))
		app.Render(r.Context(), w)
	})
}

func CategoriesAPIHandler() {
	http.HandleFunc(constants.ApiCategories, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		categoriesDb := db.NewSQLiteCategoryStore(dbInstance)
		categories, err := categoriesDb.ListCategories(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(categories)
	})
}

func ReceiptsAPIHandler() {
	http.HandleFunc(constants.ApiRecepipts, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		receiptsDb := db.NewSQLiteReceiptStore(dbInstance)
		receipts, err := receiptsDb.ListReceipts(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(receipts)
	})
}

func SingleReceiptAPIHandler() {
	http.HandleFunc(constants.ApiRecepipt, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from URL path
		id := strings.TrimPrefix(r.URL.Path, constants.ApiRecepipt)
		receiptID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid receipt ID", http.StatusBadRequest)
			return
		}

		receiptsDb := db.NewSQLiteReceiptStore(dbInstance)
		receipt, err := receiptsDb.GetReceipt(context.Background(), receiptID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if receipt == nil {
			http.Error(w, "Receipt not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(receipt)
	})
}

func PublicStaticHandler() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func ParseNewReceipt(request *http.Request) (models.NewReceipt, error) {
	newReceipt := models.NewReceipt{}
	if err := request.ParseMultipartForm(10 << 20); err != nil {
		return newReceipt, err
	}

	file, handler, err := request.FormFile(models.NewReceiptFormFieldsInstance.File)
	if err != nil {
		return newReceipt, err
	}
	defer file.Close()

	newReceipt.File = file
	newReceipt.FileName = handler.Filename

	categoryId, err := strconv.ParseInt(request.FormValue(models.NewReceiptFormFieldsInstance.CategoryID), 10, 64)
	newReceipt.CategoryID = categoryId

	date := request.FormValue(models.NewReceiptFormFieldsInstance.Date)
	newReceipt.Date = date

	return newReceipt, nil
}

func UploadHandler() {
	http.HandleFunc(constants.ApiUpload, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		newReceipt, err := ParseNewReceipt(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll("public/uploads", 0755); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new file in the uploads directory
		dst, err := os.Create(fmt.Sprintf("public/uploads/%s", newReceipt.FileName))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the destination file
		if _, err := io.Copy(dst, newReceipt.File); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		receipt := &models.Receipt{
			CategoryID: newReceipt.CategoryID,
			Date:       newReceipt.Date,
			ImageName:  dst.Name(),
		}
		db.NewSQLiteReceiptStore(dbInstance).CreateReceipt(context.Background(), receipt)

		// send receipt to the client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(receipt)
	})
}

func UpdateHandler() {
	http.HandleFunc(constants.ApiUpdateReceipt, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var receipt models.Receipt
		err := json.NewDecoder(r.Body).Decode(&receipt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.NewSQLiteReceiptStore(dbInstance).UpdateReceipt(context.Background(), &receipt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(receipt)
	})
}

func InitDb() {
	var err error
	dbInstance, err = db.InitDB("./data.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}

func main() {
	port := "127.0.0.1:8080"
	server := &http.Server{
		Addr: port,
	}

	InitDb()
	defer dbInstance.Close()

	PublicStaticHandler()
	HomeHandler()
	UploadHandler()
	CategoriesAPIHandler()
	ReceiptHandler()
	ReceiptsAPIHandler()
	SingleReceiptAPIHandler()
	UpdateHandler()

	fmt.Println("Server is running on port ", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting the server: %s\n", err)
	}
}
