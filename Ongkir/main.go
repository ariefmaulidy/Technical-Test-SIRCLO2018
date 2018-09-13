package main

import (
	"log"
	"net/http"

	"./ongkirroutes"
	"github.com/rs/cors"
	"goji.io"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	mux := goji.NewMux()

	ongkirroutes.RoutesOngkir(mux, session)

	handler := cors.New(cors.Options{AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, AllowCredentials: true}).Handler(mux)
	log.Println("Starting Listen server....")
	http.ListenAndServe("localhost:8080", handler)

}
