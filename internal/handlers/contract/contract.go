package contract

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	contract "github.com/matthiasmohr/ed4-pricechanger-go/internal/controllers"
	EntityContract "github.com/matthiasmohr/ed4-pricechanger-go/internal/entities/contract"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/handlers"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/repository/adapter"
	Rules "github.com/matthiasmohr/ed4-pricechanger-go/internal/rules"
	RulesContract "github.com/matthiasmohr/ed4-pricechanger-go/internal/rules/contract"
	httpStatus "github.com/matthiasmohr/ed4-pricechanger-go/utils/http"
	"net/http"
	"time"
)

type Handler struct {
	handlers.Interface
	Controller contract.Interface
	Rules      Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Controller: contract.NewController(repository),
		Rules:      RulesContract.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != "" {
		h.GetOne(w, r)
	} else {
		h.GetAll(w, r)
	}
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		httpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	response, err := h.Controller.ListOne(ID)

	if err != nil {
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}

	httpStatus.StatusOK(w, r, response)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}
	httpStatus.StatusOK(w, r, response)
}
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	contractBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		httpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(contractBody)

	if err != nil {
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}
	httpStatus.StatusOK(w, r, map[string]interface{}{"id": ID.String()})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		httpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	productBody, err := h.getBodyAndValidate(r, ID)

	if err != nil {
		httpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, productBody); err != nil {
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}

	httpStatus.StatusNoContent(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		httpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Remove(ID); err != nil {
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}

	httpStatus.StatusNoContent(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityContract.Contract, error) {
	contractBody := &EntityContract.Contract{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, contractBody)
	if err != nil {
		return &EntityContract.Contract{}, errors.New("body is required")
	}

	contractParsed, err := EntityContract.InterfaceToModel(body)

	if err != nil {
		return &EntityContract.Contract{}, errors.New("Errors on converting body to model")
	}

	setDefaultValues(contractParsed, ID)

	return contractParsed, h.Rules.Validate(contractParsed)
}

func setDefaultValues(contract *EntityContract.Contract, ID uuid.UUID) {
	contract.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		contract.ID = uuid.New()
		contract.CreatedAt = time.Now()
	} else {
		contract.ID = ID
	}
}
