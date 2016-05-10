default: all

all: test vet lint

test:
	go test ./...

fmt:
	gofmt -s -d ./...

lint:
	golint ./...

vet:
	go vet ./...

sloc:
	wc -l */**.go

link:
	mkdir -p $(GOPATH)/src/gopkg.in/vinxi
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/context ]; then ln -s $(PWD)/context $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/context; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/forward ]; then ln -s $(PWD)/forward $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/forward; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/intercept ]; then ln -s $(PWD)/intercept $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/intercept; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/layer ]; then ln -s $(PWD)/layer $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/layer; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/mux ]; then ln -s $(PWD)/mux $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/mux; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/router ]; then ln -s $(PWD)/router $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/router; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/utils ]; then ln -s $(PWD)/utils $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/utils; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/vinxi ]; then ln -s $(PWD)/vinxi $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/vinxi; fi;
