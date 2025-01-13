package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/db"
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

func PublicStaticHandler() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
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

	fmt.Println("Server is running on port ", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
