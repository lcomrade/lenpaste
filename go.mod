module git.lcomrade.su/root/lenpaste

go 1.16

replace git.lcomrade.su/root/lenpaste/internal => ./internal

require (
	git.lcomrade.su/root/lineend v1.0.0
	github.com/alecthomas/chroma/v2 v2.7.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
)
