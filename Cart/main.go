package main

import (
	"./cartmodel"
)

func main() {
	keranjang := new(cartmodel.Cart)
	keranjang.M = make(map[string]int)

	keranjang.TambahProduk("Topi Putih", 2)

	keranjang.TambahProduk("Kemeja Hitam", 3)

	keranjang.TambahProduk("Sepatu Merah", 1)
	keranjang.TambahProduk("Sepatu Merah", 4)
	keranjang.TambahProduk("Sepatu Merah", 2)

	keranjang.HapusProduk("Kemeja Hitam")

	keranjang.HapusProduk("Baju Hijau")

	keranjang.TampilkanCart()
}
