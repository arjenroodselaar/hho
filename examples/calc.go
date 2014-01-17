package main

func mul(x, y int) int {
	return x * y
}

func div(x, y int) int {
	return x / y
}

func main() {
	println(div(mul(3, 3), 3))
}
