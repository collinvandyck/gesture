NAME = gesture
BINARY = bin/$(NAME)
GOPATH=~/go

.PHONY: all clean distclean

all:
	@GOPATH=$(GOPATH) go build -o $(BINARY)

clean:
	@GOPATH=$(GOPATH) go clean

distclean: clean
	rm -rf $(BINARY)
