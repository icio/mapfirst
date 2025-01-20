package main

import (
	"fmt"
	"runtime"
	"strconv"
)

func main() {
	const (
		mapElements = 1000
		iterations  = 10_000
	)

	m := make(map[int]int, mapElements)
	for i := range mapElements {
		m[i] = 0
	}

	for range iterations {
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
