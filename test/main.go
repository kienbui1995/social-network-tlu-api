package main

import "fmt"

func main() {
	var n, i, j int
	fmt.Scanf("%d", &n)
	for n >= 0 {
		fmt.Scanf("%d%d", &i, &j)
		if i-j != 2 && i != j {
			fmt.Printf("No Number\n")
		} else {
			fmt.Printf("%d\n", i+j)
		}
		n--
	}
}
