package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	// The map we're testing. Update k to try different types of key.
	m := make(map[int]int)
	k := func(i int) int { return i }

	// Config options.
	n := flag.Int("n", 100, "number of elements in the map")
	d := flag.String("d", "\t", "delimited string (default tab)")
	w := flag.Int("w", 100000, "number of times to iterate over the first item")
	o := flag.String("o", fmt.Sprintf("mapfirst-%s-%T.tsv", runtime.Version(), m), "file to write results to, or - for stdout")
	flag.Parse()

	// Generate the results.
	for i := 0; i < *n; i++ {
		m[k(i)] = 0
	}
	for w := *w; w > 0; w-- {
		for k := range m {
			m[k] = m[k] + 1
			break
		}
	}

	// Open the output file.
	var f io.Writer
	if *o == "-" || *o == "" {
		f = os.Stdout
	} else {
		fc, err := os.Create(*o)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		defer func() {
			if err := fc.Close(); err != nil {
				println(err.Error())
				os.Exit(1)
			}
		}()
		f = fc
	}

	// Print the results.
	fmt.Fprintf(f, "i%s%s:%T\n", *d, runtime.Version(), m)
	for i := 0; i < *n; i++ {
		fmt.Fprint(f, i, *d, m[k(i)], "\n")
	}
}
