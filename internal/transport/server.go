package transport

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	_ "embed"

	"github.com/KoljanTestwamore/go-wb-l0/internal/database"
	"github.com/KoljanTestwamore/go-wb-l0/internal/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type message struct {
	Msg model.Order
	MsgString string
}

type Server struct {
	handler *http.ServeMux
	dbio *database.DBIO
	logger *zerolog.Event
}

func (server *Server) Run(dbio *database.DBIO) {
	server.logger = log.Info().Str("package", "server")

	server.dbio = dbio
	go dbio.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/{id}", server.handleHome)
	mux.HandleFunc("/favicon.ico", http.NotFoundHandler().ServeHTTP)

	server.handler = mux
	var port = os.Getenv("SERVER_PORT")
	var address =":"+port
	
	err := http.ListenAndServe(address, mux)
	if err != nil {
		server.logger.Err(err)
		return
	}
}

func (server *Server) handleHome(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	data := server.dbio.Read_data(id)

	var msg model.Order

	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		server.logger.Err(err)
		return
	}

	tmpl, err := template.ParseFiles("./internal/web/template.html", )
	if err != nil {
		server.logger.Err(err)
		return
	}

	err = tmpl.Execute(w, message{
		Msg: msg,
		MsgString: data,
	})

	if err != nil {
		server.logger.Err(err)
		return
	}

	server.logger.Msg("Successful Hit: homePage")
}