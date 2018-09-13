package beratroutes

import (
	"../beratcontroller"
	"goji.io"
	"goji.io/pat"
	"gopkg.in/mgo.v2"
)

func RoutesBerat(mux *goji.Mux, session *mgo.Session) {
	mux.HandleFunc(pat.Get("/index"), beratcontroller.AllBerat(session))
	mux.HandleFunc(pat.Get("/show/:idberat"), beratcontroller.OneBerat(session))
	mux.HandleFunc(pat.Post("/add"), beratcontroller.AddBerat(session))
	mux.HandleFunc(pat.Put("/update/:idberat"), beratcontroller.UpdateBerat(session))
	mux.HandleFunc(pat.Delete("/delete/:idberat"), beratcontroller.DeleteBerat(session))
}
