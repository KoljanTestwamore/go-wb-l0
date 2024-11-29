package services

import (
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/hashicorp/go-uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Класс для проверки
type Producer struct {
	producer sarama.SyncProducer 
	logger *zerolog.Event
}

func (p *Producer) Run() {
	var port = os.Getenv("KAFKA_PORT")
	p.logger = log.Info().Str("package", "producer")

	producer, err := sarama.NewSyncProducer([]string{"kafka:"+port}, nil)
	if (err != nil) {
		p.logger.Err(err)
		return
	}

	p.producer = producer

	p.SendMessage()
}

func (p *Producer) SendMessage() {
	id, err := uuid.GenerateUUID()
	if err != nil {
		p.logger.Err(err)
		return
	}

	value := fmt.Sprintf(`{
   "order_uid": "%s",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`, id)

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Key: sarama.StringEncoder(id),
		Value: sarama.ByteEncoder(value),
		Partition: 0,
	}

	_, _, err = p.producer.SendMessage(msg)
	if err!=nil {
		p.logger.Err(err)
		return
	}

	p.logger.Msg(fmt.Sprintf("Published message with uuid %s", id))
}