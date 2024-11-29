package services

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/IBM/sarama"
	"github.com/KoljanTestwamore/go-wb-l0/internal/database"
	"github.com/KoljanTestwamore/go-wb-l0/internal/model"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Consumer struct {
	dbio *database.DBIO
	logger *zerolog.Event
}

func (c *Consumer) Run(dbio *database.DBIO) {
	c.dbio = dbio
	c.logger = log.Info().Str("package", "consumer")

	var port = os.Getenv("KAFKA_PORT")

	consumer, err := sarama.NewConsumer([]string{"kafka:"+port}, nil)
	if (err != nil) {
		c.logger.Err(err)
		return
	}

	partConsumer, err := consumer.ConsumePartition("orders", 0, sarama.OffsetNewest)
	if err != nil {
		c.logger.Err(err)
		return
	}
	defer partConsumer.Close()

	for msg := range partConsumer.Messages() {
		// Проверка, что данные корректные
		// Сделано в слое с бизнес логикой, так как
		// проверка данных из канала должна быть тут
		var receivedMessage model.Order
		err := json.Unmarshal(msg.Value, &receivedMessage)

		if err != nil {
			c.logger.Err(err)
			continue
		}

		err = c.dbio.Write_data(string(msg.Value))

		if err != nil {
			c.logger.Err(err)
			continue
		}
	}
}
