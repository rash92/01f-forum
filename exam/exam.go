package main

import (
	"strconv"

	"github.com/01-edu/z01"
)

func main() {
	arr := []int{1, 2, 3, 4, 5}
	ReduceInt(arr, Add)
}

func ReduceInt(a []int, f func(int, int) int) {
	acc := a[0]
	for i := 1; i < len(a); i++ {
		acc = f(acc, a[i])
	}
	result := strconv.Itoa(acc)

	for _, l := range result {
		z01.PrintRune(l)
	}
	// for i := 0; i < len(result); i++ {
	// 	z01.PrintRune(result[i])
	// }
	z01.PrintRune('\n')
}

func Add(x, y int) int {
	return x + y
}

func Multiply(x, y int) int {
	return x * y
}
