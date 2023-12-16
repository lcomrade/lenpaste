module github.com/lcomrade/lenpaste

go 1.16

replace github.com/lcomrade/lenpaste/internal => ./internal

require (
	github.com/alecthomas/chroma/v2 v2.4.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
)
