module git.lcomrade.su/root/lenpaste

go 1.18

replace git.lcomrade.su/root/lenpaste/internal => ./internal

require (
	git.lcomrade.su/root/lineend v1.0.0
	github.com/alecthomas/chroma/v2 v2.7.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/urfave/cli/v2 v2.25.1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
)
