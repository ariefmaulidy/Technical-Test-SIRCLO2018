package main

func fivaa(x int) {

	for i := 0; x > 0; i++ {
		for j := 0; j <= x+1; j++ {
			if j <= 1 {
				print(x - 1)
			} else if j != x+1 {
				print(x + 1)
			} else {
				print(x+1, "\n")
			}
		}
		x--
	}
}

func main() {
	fivaa(5)
}
