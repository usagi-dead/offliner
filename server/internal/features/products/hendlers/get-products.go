package hendlers

import (
	"log/slog"
	"net/http"
)

type Catalog interface {
	GetProducts(category string)
}

func New(log *slog.Logger, catalog Catalog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
