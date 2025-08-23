package handler

import (
	"encoding/json"
	"log"
	"myapp/internal/cache"
	"net/http"
	"strings"
)

type OrderHandler struct {
	cache *cache.OrderCache
}

func NewOrderHandler(c *cache.OrderCache) *OrderHandler {
	return &OrderHandler{cache: c}
}

func (h *OrderHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	// Получаем order_uid из URL, например /order/b563feb7b2b84b6test
	uid := strings.TrimPrefix(r.URL.Path, "/order/")
	if uid == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}

	order, found := h.cache.Get(uid)
	if !found {
		log.Printf("Order not found in cache: %s", uid)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Для простоты разработки
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode order", http.StatusInternalServerError)
	}
}
