package health

import (
	"errors"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/handlers"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/repository/adapter"
	httpStatus "github.com/matthiasmohr/ed4-pricechanger-go/utils/http"
	"net/http"
)

type Handler struct {
	handlers.Interface
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Repository: repository,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.Repository.Health() {
		httpStatus.StatusInternalServerError(w, r, errors.New("Relational database not alive"))
		return
	}
	httpStatus.StatusOK(w, r, "Service OK")
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodAllowed(w, r)
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodAllowed(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodAllowed(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodAllowed(w, r)
}
