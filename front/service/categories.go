package service

import (
	"encoding/json"
	"net/http"

	"github.com/dudubtw/receipt/models"
)

func FetchCategories() ([]models.Category, error) {
	ch := make(chan struct {
		categories []models.Category
		err        error
	})

	go func() {
		resp, err := http.Get("/api/categories")
		if err != nil {
			ch <- struct {
				categories []models.Category
				err        error
			}{nil, err}
			return
		}
		defer resp.Body.Close()

		var categories []models.Category
		if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
			ch <- struct {
				categories []models.Category
				err        error
			}{nil, err}
			return
		}

		ch <- struct {
			categories []models.Category
			err        error
		}{categories, nil}
	}()

	result := <-ch
	return result.categories, result.err
}
