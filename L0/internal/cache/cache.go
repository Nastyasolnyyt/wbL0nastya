package cache

import (
	"myapp/internal/model"
	"sync"
)

type OrderCache struct {
	store sync.Map
}

func New() *OrderCache {
	return &OrderCache{}
}

func (c *OrderCache) Get(orderUID string) (model.Order, bool) {
	value, found := c.store.Load(orderUID)
	if !found {
		return model.Order{}, false
	}

	order, ok := value.(model.Order)
	if !ok {
		c.store.Delete(orderUID)
		return model.Order{}, false
	}

	return order, true
}

func (c *OrderCache) Set(order model.Order) {
	c.store.Store(order.OrderUID, order)
}

func (c *OrderCache) Load(orders []model.Order) {
	for _, order := range orders {
		c.Set(order)
	}

}
