module git.lcomrade.su/root/lenpaste

go 1.11

replace git.lcomrade.su/root/lenpaste/internal => ./internal

require (
	git.lcomrade.su/root/lineend v1.0.0
	github.com/alecthomas/chroma v0.10.0
	github.com/mattn/go-sqlite3 v1.14.15
	github.com/lib/pq v1.10.6
)
