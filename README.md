# mapfirst

mapfirst builds a structure like a map[int]int, then records the key `k` first
visited in `for k := range m`. It tests a few different types of map with
parameters controlling different map sizes and number of first keys to look at.
It is written with Go 1.0.0 compatiblity to make comparison easy.

By default, the results are written to a file `mapfirst-$GOOS-$GOARCH-$GOVERSION.tsv`.

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

If you're using Google Sheets, try pasting that into A2 and setting call A1's
formula to `=SPARKLINES(A1:A)` and replicating it all the way across row 1.
