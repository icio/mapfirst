package main

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
)

func main() {
	n := flag.Int("n", 100, "number of elements in the map")
	w := flag.Int("w", 100_000, "number of times to iterate over the first item")
	flag.Parse()

	m := make(map[int]int, *n)
	for i := range *n {
		m[i] = 0
	}

	for range *w {
		for i := range m {
			m[i] = m[i] + 1
			break
		}
	}

	fmt.Print("i,", runtime.Version(), "\n")
	for i := range len(m) {
		fmt.Print(strconv.Itoa(i), ",", strconv.Itoa(m[i]), "\n")
	}
}
