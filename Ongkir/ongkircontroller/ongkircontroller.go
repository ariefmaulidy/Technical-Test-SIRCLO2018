package ongkircontroller

import (
	"encoding/json"
	"net/http"

	"../jsonhandler"
	"../ongkirmodel"
	"goji.io/pat"
	mgo "gopkg.in/mgo.v2"
)

func CalculateOngkir(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		metode := pat.Param(r, "metode")
		var ongkir ongkirmodel.Ongkir
		var varmap map[string]interface{}
		varmap = make(map[string]interface{})
		in := []byte(`{}`)
		json.Unmarshal(in, &varmap)

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&varmap)
		if err != nil {
			jsonhandler.SendJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		if metode == "regular" {
			ongkir.Hargatotal = int(varmap["berat"].(float64)) * 9000
		} else if metode == "express" {
			ongkir.Hargatotal = int(varmap["berat"].(float64)) * 9000 * 2
		} else if metode == "instant" {
			ongkir.Hargatotal = int(varmap["berat"].(float64)) * 9000 * 5
		} else {
			jsonhandler.SendJSON(w, "Bukan metode yang sesuai", http.StatusInternalServerError)
		}
		respBody, err := json.MarshalIndent(ongkir, "", "  ")
		if err != nil {
			// log.Fatal(err)
		}
		jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
