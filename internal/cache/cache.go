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

func (oc *OrderCache) CreateCache(ord models.Order) {
	err := oc.dbRepo.SaveOrder(ord)

	if err != nil {
		fmt.Printf("Cannot insert order: %v", err)
	}

	oc.mu.Lock()
	oc.cache[ord.OrderUid] = &ord // запись данных в кэш
	oc.mu.Unlock()
	fmt.Printf("Cache written: %s\n", ord.OrderUid)
}

func (oc *OrderCache) Preload() {

	ors, err := oc.dbRepo.GetALl()
	if err != nil {
		fmt.Printf("Error at DB: %v\n", err)
	}
	fmt.Printf("DB returns len: %d\n", len(ors))
	oc.mu.Lock()
	for _, or := range ors {
		// Создаем копию переменной or для каждой итерации цикла
		// и добавляем эту копию в кэш
		copyOrder := or
		oc.cache[copyOrder.OrderUid] = &copyOrder
	}
	oc.mu.Unlock()
	/*ors, err := oc.dbRepo.GetALl()
	if err != nil {
		fmt.Printf("Error at DB: %v\n", err)
	}
	fmt.Printf("DB returns len: %d\n", len(ors))
	oc.mu.Lock()
	for _, or := range ors {
		oc.cache[or.OrderUid] = &or
	}
	oc.mu.Unlock()*/
}

// возвращение заказа
func (oc *OrderCache) GetOrderByUid(uid string) *models.Order {
	return oc.cache[uid]
}

// возвращение всего кэша
func (oc *OrderCache) GetOrders() map[string]*models.Order {
	return oc.cache
}
