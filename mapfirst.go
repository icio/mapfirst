package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type mapfuncs struct {
	get      func(i int) int64
	zero     func(i int)
	rangeInc func(done func() bool)
}

func (mf mapfuncs) Zero(i int)                { mf.zero(i) }
func (mf mapfuncs) Get(i int) int64           { return mf.get(i) }
func (mf mapfuncs) RangeInc(done func() bool) { mf.rangeInc(done) }

type intslice struct {
	ints []int
}

var _ flag.Value = (*intslice)(nil)

func (s *intslice) Set(vs string) error {
	s.ints = s.ints[:0]
	for _, v := range strings.Split(vs, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return err
		}
		s.ints = append(s.ints, n)
	}
	return nil
}

func (s *intslice) String() string {
	if len(s.ints) == 0 {
		return ""
	}

	b := bytes.NewBuffer(make([]byte, 3*len(s.ints)))
	for i, n := range s.ints {
		if i > 1 {
			b.WriteRune(',')
		}
		b.WriteString(strconv.Itoa(n))
	}
	return b.String()
}

func main() {
	expN := intslice{[]int{100, 1000, 10000}}
	expJ := intslice{[]int{1, 2, 3}}

	// Config options.
	flag.Var(&expN, "n", "map sizes to test (comma-separated)")
	flag.Var(&expJ, "j", "number of range iterations to test (comma-separated)")
	p := flag.String("p", runtime.GOOS+"-"+runtime.GOARCH+"-"+runtime.Version(), "experiment name prefix")
	d := flag.String("d", "\t", "delimited string (default tab)")
	o := flag.String("o", "", "file to write results to, or - for stdout (defaults to mapfirst-<goos>-<goarch>-<n>-<j>.tsv)")
	flag.Parse()

	// Different experiments, covering different key types and insert methods.
	type maplike interface {
		Zero(i int)
		Get(i int) int64
		RangeInc(done func() bool)
	}
	type experiment struct {
		name string              // name of the experiment (we later add "-n<n>-j<j>")
		make func(n int) maplike // map/key structure.
		n    int                 // map size.
		j    int                 // number of keys to visit from start of range.
		w    int                 // number of ranges to sample.
	}
	// tpls define the experiments without n, j, w which get filled in below.
	tpls := []experiment{
		{
			name: "int",
			make: func(n int) maplike {
				m := make(map[int]int64, n)
				return mapfuncs{
					zero: func(i int) { m[i] = 0 },
					get:  func(i int) int64 { return m[i] },
					rangeInc: func(done func() bool) {
						for k := range m {
							m[k] = m[k] + 1
							if done() {
								return
							}
						}
					},
				}
			},
		},
		{
			name: "intrev",
			make: func(n int) maplike {
				m := make(map[int]int64)
				return mapfuncs{
					zero: func(i int) { m[n-i] = 0 },
					get:  func(i int) int64 { return m[n-i] },
					rangeInc: func(done func() bool) {
						for k := range m {
							m[k] = m[k] + 1
							if done() {
								return
							}
						}
					},
				}
			},
		},
		{
			name: "int64",
			make: func(n int) maplike {
				m := make(map[int64]int64, n)
				return mapfuncs{
					zero: func(i int) { m[int64(i)] = 0 },
					get:  func(i int) int64 { return m[int64(i)] },
					rangeInc: func(done func() bool) {
						for k := range m {
							m[k] = m[k] + 1
							if done() {
								return
							}
						}
					},
				}
			},
		},
		{
			name: "strhex",
			make: func(n int) maplike {
				m := make(map[string]int64, n)
				k := func(i int) string { return fmt.Sprintf("%x", i) }
				return mapfuncs{
					zero: func(i int) { m[k(i)] = 0 },
					get:  func(i int) int64 { return m[k(i)] },
					rangeInc: func(done func() bool) {
						for k := range m {
							m[k] = m[k] + 1
							if done() {
								return
							}
						}
					},
				}
			},
		},
	}

	// Replicate these experiments with extra variables.
	var exps []experiment
	var N int
	for _, n := range expN.ints {
		if n > N {
			N = n
		}
		for _, j := range expJ.ints {
			if j >= n {
				continue
			}
			for _, t := range tpls {
				exps = append(exps, experiment{
					name: fmt.Sprintf("%s-n%d-j%d", t.name, n, j),
					n:    n,
					j:    j,
					w:    n * 100,
					make: t.make,
				})
			}
		}
	}

	// Run all experiments in exps, copying results to dists.
	dists := make([][]int64, len(exps))
	for ei, exp := range exps {
		println(*p + "-" + exp.name)
		m := exp.make(exp.n)

		// Generate the results.
		dists[ei] = make([]int64, exp.n)
		for i := 0; i < exp.n; i++ {
			m.Zero(i)
		}
		for w := exp.w; w > 0; w-- {
			j := exp.j
			m.RangeInc(func() bool {
				j--
				return j == 0
			})
		}

		// Store the results.
		for i := 0; i < exp.n; i++ {
			dists[ei][i] = m.Get(i)
		}
	}

	// Open the output file.
	var f io.Writer
	if *o == "-" {
		f = os.Stdout
	} else {
		if *o == "" {
			*o = "mapfirst-" + *p + ".tsv"
		}
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

	// Print the result header.
	for ei, exp := range exps {
		if ei > 0 {
			fmt.Fprint(f, *d)
		}
		fmt.Fprint(f, *p+"-"+exp.name)
	}
	fmt.Fprintln(f)

	for i := 0; i < N; i++ {
		for ei := range exps {
			if ei > 0 {
				fmt.Fprint(f, *d)
			}
			if i < exps[ei].n {
				fmt.Fprint(f, dists[ei][i])
			}
		}
		fmt.Fprintln(f)
	}
}
