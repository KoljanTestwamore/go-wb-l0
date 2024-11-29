package app

import (
	"github.com/KoljanTestwamore/go-wb-l0/internal/database"
	"github.com/KoljanTestwamore/go-wb-l0/internal/services"
	"github.com/KoljanTestwamore/go-wb-l0/internal/transport"
)


var server   = transport.Server{}
var dbio     = database.DBIO{}
var consumer = services.Consumer{}
// var producer = services.Producer{}

func Run() {
	
	dbio.Init()
	go consumer.Run(&dbio)
	// Продюсер нужен для генерации тестовых сообщений в кафку.
	// go producer.Run()
	server.Run(&dbio)
}