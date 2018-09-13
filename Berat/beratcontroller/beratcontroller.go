package beratcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"../beratmodel"
	"../jsonhandler"
	"goji.io/pat"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
Fungsi untuk membuat indeks pada basis data MongoDB dengan menggunakan library mgo v2.
*/
func EnsureBerat(s *mgo.Session) {
	session := s.Copy()
	defer session.Close() // Defer berfungsi untuk mengeksekusi fungsi ketika fungsi utama sudah mengembalikan sesuatu

	c := session.DB("testing").C("berat") // mengakses db bernama "testing" dengan collection "berat"

	index := mgo.Index{
		Key:        []string{"idberat"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

/*
Fungsi untuk mengambil keseluruhan data pada basis data berat yang ada pada halaman index
*/
func AllBerat(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("testing").C("berat")

		var berats []beratmodel.Berat
		var rataanmax float64
		var rataanmin float64
		var rataanbeda float64
		var rataan []float64
		var datasend beratmodel.DataIndex

		err := c.Find(bson.M{}).Sort("-tanggal").All(&berats) // Mengambil semua data berat yang diurutkan berdasarkan tanggalnya secara menurun dan menaruhnya pada variabel berats
		if err != nil {
			jsonhandler.SendJSON(w, "Database error", http.StatusInternalServerError) // Dapat dilihat pada jsonhandler.go
			log.Println("Failed get all berat: ", err)
			return
		}

		for _, data := range berats { // Loop untuk menghitung rata2 setiap data
			rataanmax += float64(data.Max)
			rataanmin += float64(data.Min)
			rataanbeda += float64(data.Max - data.Min)
			datasend.DataPerbedaan = append(datasend.DataPerbedaan, data.Max-data.Min)
			tgl := time.Unix(data.Tanggal, 0)
			tglstring := tgl.Format("2006-01-02 15:04:05")
			datasend.TanggalString = append(datasend.TanggalString, tglstring)
		}
		rataanmax /= float64(len(berats))
		rataan = append(rataan, rataanmax)
		rataanmin /= float64(len(berats))
		rataan = append(rataan, rataanmin)
		rataanbeda /= float64(len(berats))
		rataan = append(rataan, rataanbeda)

		datasend.DataBerat = berats
		datasend.DataRataan = rataan

		respBody, err := json.MarshalIndent(datasend, "", "  ")
		if err != nil {
			// log.Fatal(err)
		}
		jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

/*
Fungsi untuk mengambil sebuah data pada basis data berat yang ada pada halaman show
*/
func OneBerat(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("testing").C("berat")

		IDBerat := pat.Param(r, "idberat")
		var berat beratmodel.Berat
		var datasend beratmodel.DataShow

		err := c.Find(bson.M{"idberat": IDBerat}).One(&berat) // Mengambil sebuah data berat berdasarkan Idnya dan menaruhnya pada variabel berat
		if err != nil {
			jsonhandler.SendJSON(w, "Database error", http.StatusInternalServerError) // Dapat dilihat pada jsonhandler.go
			log.Println("Failed get one berat: ", err)
			return
		}

		datasend.DataBerat = berat
		datasend.DataPerbedaan = berat.Max - berat.Min
		tgl := time.Unix(berat.Tanggal, 0)
		datasend.TanggalString = tgl.Format("2006-01-02 15:04:05")

		respBody, err := json.MarshalIndent(datasend, "", "  ")
		if err != nil {
			// log.Fatal(err)
		}

		jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

/*
Fungsi untuk menambah sebuah data berat
*/
func AddBerat(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var berat beratmodel.Berat
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&berat)
		if err != nil {
			jsonhandler.SendJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB("testing").C("berat")

		//untuk auto increment
		var lastBerat beratmodel.Berat
		var lastId int

		err = c.Find(nil).Sort("-$natural").Limit(1).One(&lastBerat)
		if err != nil {
			lastId = 0
		} else {
			lastId, err = strconv.Atoi(lastBerat.IDBerat)
		}
		currentId := lastId + 1
		berat.IDBerat = strconv.Itoa(currentId)

		err = c.Insert(berat) // Menambahkan data berat
		if err != nil {
			if mgo.IsDup(err) {
				jsonhandler.SendJSON(w, "duplicate", http.StatusOK)
				return
			}

			jsonhandler.SendJSON(w, "Database error", http.StatusNotFound)
			log.Println("Failed insert berat: ", err)
			return
		}
		respBody, err := json.MarshalIndent(berat, "", "  ")
		if err != nil {
			// log.Fatal(err)
		}

		jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

/*
Fungsi untuk mengedit sebuah data berat
*/
func UpdateBerat(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		IDBerat := pat.Param(r, "idberat")

		var berat beratmodel.Berat
		var varmap map[string]interface{}
		varmap = make(map[string]interface{})
		in := []byte(`{}`)
		json.Unmarshal(in, &varmap)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&berat)
		if err != nil {
			jsonhandler.SendJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}
		varmap["max"] = berat.Max
		varmap["min"] = berat.Min
		c := session.DB("testing").C("berat")
		err = c.Update(bson.M{"idberat": IDBerat}, bson.M{"$set": varmap})
		if err != nil {
			switch err {
			default:
				jsonhandler.SendJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed update lapak: ", err)
				jsonhandler.SendJSON(w, "Gagal mengupdate lapak", http.StatusOK)
				return
			case mgo.ErrNotFound:
				jsonhandler.SendJSON(w, "lapak not found", http.StatusNotFound)
				return
			}
		}
		err = c.Find(bson.M{"idberat": IDBerat}).One(&berat)
		respBody, err := json.MarshalIndent(berat, "", "  ")
		if err != nil {
			// log.Fatal(err)
		}

		jsonhandler.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

/*
Fungsi untuk menghapus sebuah data berat
*/
func DeleteBerat(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		IDBerat := pat.Param(r, "idberat")

		c := session.DB("testing").C("berat")

		err := c.Remove(bson.M{"idberat": IDBerat}) // Menghapus data berat berdasarkan Id nya
		if err != nil {
			switch err {
			default:
				jsonhandler.SendJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed delete lapak: ", err)
				return
			case mgo.ErrNotFound:
				jsonhandler.SendJSON(w, "lapak not found", http.StatusNotFound)
				return
			}
		}
		jsonhandler.SendJSON(w, "Data berhasil dihapus", http.StatusOK)
	}
}
