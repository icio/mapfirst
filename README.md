# mapfirst

mapfirst builds a structure like a map[int]int, then records the key `k` first
visited in `for k := range m`. It is written to be runnable with many different
versions of Go to make comparison easy.

Modify the definitions of `m` and `k` in mapfirst.go if you want to try test
different types.

```sh
# Clone the source and enter the directory.
git clone https://github.com/icio/mapfirst
cd mapfirst

# Run mapfirst with each desired Go version.
for R in go1.23.5 go1.24rc2; do
    go install golang.org/dl/$R@latest
    $R download
    $R run .
done

# Copy the combined TSV into your graphing software.
paste mapfirst*.tsv | pbcopy
```
