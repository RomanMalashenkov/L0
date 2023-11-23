package generate

import (
	"fmt"
	"math/rand"
	"strconv"
	"test_wb/internal/models"
	"time"

	"github.com/google/uuid" // Импорт пакета для генерации UUID
)

func GenerateOrder() *models.Order {
	orderUid := uuid.NewString() // Генерация уникального идентификатора заказа

	var itemCount = 1 + rand.Intn(5) // Генерация случайного количества товаров от 1 до 5
	items := make([]models.Item, itemCount)
	for i := range items {
		items[i] = models.Item{
			ChrtId:      0,
			TrackNumber: "",
			Price:       0,
			Rid:         "",
			Name:        "",
			Sale:        0,
			Size:        "",
			TotalPrice:  0,
			NmId:        0,
			Brand:       "",
			Status:      0,
		}
	}

	// Формирование заказа
	order := models.Order{
		OrderUid:    orderUid,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery:    generateOrderDelivery(),
		Payment: models.Payment{
			Transaction:  "",
			RequestId:    "",
			Currency:     "",
			Provider:     "",
			Amount:       0,
			PaymentDt:    0,
			Bank:         "",
			DeliveryCost: 0,
			GoodsTotal:   0,
			CustomFee:    0,
		},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmId:              99,
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          "1",
	}

	generateOrderItems(&order)
	generateOrderPayment(&order)

	return &order
}

// Генерация случайной информации о доставке
func generateOrderDelivery() models.Delivery {
	names := []string{"Test Testov", "Roman Malashenckov", "Ivan Ivanov"}
	addresse := []string{"Ploshad Mira 15", "Naberegnaja 243", "Prospect Randoma 2"}

	delivery := &models.Delivery{
		Name:    names[rand.Intn(len(names))],
		Phone:   "+7" + strconv.Itoa(1000000000+rand.Intn(9000000000)),
		Zip:     strconv.Itoa(100000 + rand.Intn(900000)),
		City:    "Kiryat Mozkin",
		Address: addresse[rand.Intn(len(addresse))],
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	}

	return *delivery
}

// Генерация случайной информации о платеже
func generateOrderPayment(order *models.Order) {
	currency := []string{"EUR", "RUB", "USD"}
	bank := []string{"alpha", "sber", "tink"}

	var amount float64 = 0
	for i := range order.Items {
		amount += order.Items[i].TotalPrice
	}

	deliveryCost := float64(1 + rand.Intn(1500))

	order.Payment = models.Payment{
		Transaction:  order.OrderUid,
		RequestId:    "",
		Currency:     currency[rand.Intn(len(currency))],
		Provider:     "wbpay",
		Amount:       amount + deliveryCost,
		PaymentDt:    uint32(1000000000 + rand.Intn(1000000000)),
		Bank:         bank[rand.Intn(len(bank))],
		DeliveryCost: uint32(deliveryCost),
		GoodsTotal:   amount,
		CustomFee:    0,
	}
}

// Генерация случайных товаров
func generateOrderItems(order *models.Order) {

	for i := range order.Items {
		price := float64(rand.Intn(10000))
		sale := float64(rand.Intn(100)) // Случайная скидка от 0 до 100
		totalPrice := price * ((100 - sale) / 100.0)

		order.Items[i] = models.Item{
			ChrtId:      uint32(rand.Intn(10000000)),
			TrackNumber: order.TrackNumber,
			Price:       price,
			Rid:         uuid.New().String(),
			Name:        fmt.Sprintf("Item%d", i+1),
			Sale:        uint8(sale),
			Size:        "0",
			TotalPrice:  totalPrice,
			NmId:        uint32(rand.Intn(10000000)),
			Brand:       "Vivienne Sabo",
			Status:      202,
		}
	}
}
