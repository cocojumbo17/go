package main

import (
	"fmt"
)

const (
	fir  = iota
	sec  = 5 * iota
	tree = 3 * iota
	qu
	pe
)

func zamicanie() {
	type prefPrint func(string)
	printWithPrefix := func(pref string) prefPrint {
		return func(text string) {
			fmt.Printf("[%s] %s", pref, text)
		}
	}

	printer := printWithPrefix("SUPHIX")
	printer("some text")
}

type myArr []int

func (ar *myArr) Add(val int) {
	*ar = append(*ar, val)
}

func (ar *myArr) Count() int {
	return len(*ar)
}

func methods() {
	ar := myArr([]int{1, 2})
	fmt.Println(ar.Count(), ar)
	ar.Add(5)
	fmt.Println(ar.Count(), ar)
}
func main() {
	fmt.Printf("Hello GO %d %d %d %d %d", fir, sec, tree, qu, pe)
	fmt.Println("-------------------")
	zamicanie()
	fmt.Println("-------------------")
	methods()
}
