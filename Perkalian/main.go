package main

import ()

func perkalian(x int,y int) int{
	result:=0
	if x<=y{
		for i:=1;i<=x;i++{
			result+=y
		}
	}else{
		for i:=1;i<=y;i++{
			result+=x
		}
	}
	return result
}

func main() {
	println(perkalian(6,12))
}