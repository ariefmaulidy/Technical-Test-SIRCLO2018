package beratmodel

/*
Untuk disimpan menjadi collection pada basis data MongoDB. Tanggal menggunakan Unix timestamp
*/
type Berat struct {
	IDBerat string `json:"idberat"`
	Tanggal int64  `json:"tanggal"`
	Max     int    `json:"max"`
	Min     int    `json:"min"`
}

/*
Untuk JSON return pada tampilan index
*/
type DataIndex struct {
	DataBerat     []Berat   `json:"databerat"`
	TanggalString []string  `json:"tanggalstring"`
	DataPerbedaan []int     `json:"dataperbedaan"`
	DataRataan    []float64 `json:"datarataan"`
}

/*
Untuk JSON return pada tampilan show
*/
type DataShow struct {
	DataBerat     Berat  `json:"databerat"`
	DataPerbedaan int    `json:"dataperbedaan"`
	TanggalString string `json:"tanggalstring"`
}
