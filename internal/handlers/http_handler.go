package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tvbondar/go-server/internal/usecases"
)

type HTTPHandler struct {
	usecase *usecases.GetOrderUseCase
}

func NewHTTPHandler(usecase *usecases.GetOrderUseCase) *HTTPHandler {
	return &HTTPHandler{usecase: usecase}
}

func (h *HTTPHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	order, err := h.usecase.Execute(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
