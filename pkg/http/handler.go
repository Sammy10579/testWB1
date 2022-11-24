package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"testWB/pkg/user"
)

type Handler struct {
	storage   *user.Storage
	validator *validator.Validate
}

func NewHandler(storage *user.Storage) *Handler {
	return &Handler{storage: storage, validator: validator.New()}
}

type SetUserGradeRequest struct {
	UserId        string `json:"user_id"  validate:"required"`
	PostpaidLimit *int   `json:"postpaid_limit,omitempty"`
	Spp           *int   `json:"spp,omitempty"`
	ShippingFee   *int   `json:"shipping_fee,omitempty"`
	ReturnFee     *int   `json:"return_fee,omitempty"`
}

func (h *Handler) SetUserGrade(responseWriter http.ResponseWriter, req *http.Request) {
	var request SetUserGradeRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.storage.Get(request.UserId)
	if errors.Is(err, user.ErrNotFound) {
		u = user.UserGrade{UserId: request.UserId}
	}
	if request.Spp != nil {
		u.Spp = *request.Spp
	}
	if request.PostpaidLimit != nil {
		u.PostpaidLimit = *request.PostpaidLimit
	}
	if request.ShippingFee != nil {
		u.ShippingFee = *request.ShippingFee
	}
	if request.ReturnFee != nil {
		u.ReturnFee = *request.ReturnFee
	}
	h.storage.Set(u)

	jsonResponse(responseWriter, u)

}

func (h *Handler) GetUserGrade(responseWriter http.ResponseWriter, request *http.Request) {
	userID := request.URL.Query().Get("user_id")

	if err := h.validator.Var(userID, "required"); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	userGrade, err := h.storage.Get(userID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			http.Error(responseWriter, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(responseWriter, userGrade)
}

func jsonResponse(responseWriter http.ResponseWriter, any any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(responseWriter).Encode(any); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}
