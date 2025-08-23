package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "orders"
	broker := "localhost:9092"

	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	// Читаем JSON-модель из файла
	jsonData, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatalf("could not read model.json: %v", err)
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("b563feb7b2b84b6test"), // Используем UID как ключ
			Value: jsonData,
		},
	)
	if err != nil {
		log.Fatalf("failed to write messages: %v", err)
	}

	log.Println("Message sent successfully!")

	// Отправим еще одно сообщение c другим UID для теста
	time.Sleep(1 * time.Second)
	anotherOrder := `{"order_uid": "qazwsxedcrfvtgbtest", "track_number": "TRACKNEW", "entry": "WBIL", "delivery": {"name": "Ivan Ivanov", "phone": "+71998889999", "zip": "123456", "city": "Moscow", "address": "Dom 1", "region": "Moscow", "email": "ivan@gmail.com"}, "payment": {"transaction": "qazwsxedcrfvtgbtest", "currency": "RUB", "provider": "sbp", "amount": 5000, "payment_dt": 1637907727, "bank": "sber", "delivery_cost": 300, "goods_total": 4700, "custom_fee": 0}, "items": [{"chrt_id": 123, "track_number": "TRACKNEW", "price": 4700, "rid": "rid123", "name": "Go Book", "sale": 0, "size": "0", "total_price": 4700, "nm_id": 456, "brand": "O'Reilly", "status": 202}], "locale": "ru", "customer_id": "ivan_test", "delivery_service": "cdek", "shardkey": "1", "sm_id": 10, "date_created": "2023-10-27T10:00:00Z", "oof_shard": "1"}`
	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("qazwsxedcrfvtgb"),
			Value: []byte(anotherOrder),
		},
	)
	if err != nil {
		log.Fatalf("failed to write second message: %v", err)
	}
	log.Println("Second message sent successfully!")
}
