package ongkirroutes

import (
	"../ongkircontroller"
	"goji.io"
	"goji.io/pat"
	"gopkg.in/mgo.v2"
)

func RoutesOngkir(mux *goji.Mux, session *mgo.Session) {
	mux.HandleFunc(pat.Post("/calculate/:metode"), ongkircontroller.CalculateOngkir(session))
}
