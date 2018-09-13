package cartmodel

type Cart struct {
	M map[string]int
}

func (c Cart) New() Cart {
	c.M = make(map[string]int)
	return c
}

func (c Cart) TambahProduk(kodeProduk string, kuantitas int) {
	if value, exits := c.M[kodeProduk]; exits {
		c.M[kodeProduk] = value + kuantitas
	} else {
		c.M[kodeProduk] = kuantitas
	}
}

func (c Cart) HapusProduk(kodeProduk string) {
	delete(c.M, kodeProduk)
}

func (c Cart) TampilkanCart() {
	for k, v := range c.M {
		println(k, "(", v, ")")
	}
}
