package cache

import (
	"fmt"
	"sync"
	"test_wb/internal/models"
	"test_wb/internal/repository"
)

type OrderCache struct {
	cache  map[string]*models.Order
	dbRepo *repository.Repo
	mu     sync.Mutex
}

func NewCache(repo *repository.Repo) *OrderCache {
	return &OrderCache{
		cache:  map[string]*models.Order{},
		dbRepo: repo,
		mu:     sync.Mutex{},
	}
}

func (oc *OrderCache) CreateCache(or models.Order) {
	err := oc.dbRepo.SaveOrder(or)

	if err != nil {
		fmt.Printf("Cannot insert order: %v", err)
	}

	oc.mu.Lock()
	oc.cache[or.OrderUid] = &or
	oc.mu.Unlock()
	fmt.Printf("Cache written: %s\n", or.OrderUid)
}

func (oc *OrderCache) GetOrderByUid(uid string) *models.Order {
	return oc.cache[uid]
}

func (oc *OrderCache) GetOrders() map[string]*models.Order {
	return oc.cache
}
