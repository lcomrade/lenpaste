# Frequently Asked Questions
## Build from source code
### Errors
#### cannot find module for path _/x/lenpaste/internal/api
This error appears on GO 1.16+. Solution:
1. Run: `export GO111MODULE=off`
2. Try to build again: `go build ./cmd/lenpaste.go`
