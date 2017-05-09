
# The ldflags here halve the size of the binary, but remove the debugging symbols.
# See: https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
magic: main.go
	go build -o magic -ldflags="-s -w" main.go

clean:
	- rm magic

.PHONY: clean
