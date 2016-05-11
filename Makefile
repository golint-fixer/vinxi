default: all

all: test vet lint

test:
	go test -v -race ./...

fmt:
	gofmt -s -d ./...

lint:
	golint ./...

vet:
	go vet ./...

sloc:
	wc -l */**.go

update:
	go get -u ./...

link:
	mkdir -p $(GOPATH)/src/gopkg.in/vinxi
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0 ]; then ln -s $(PWD) $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/manager ]; then ln -s $(PWD)/manager $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/manager; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/context ]; then ln -s $(PWD)/context $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/context; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/forward ]; then ln -s $(PWD)/forward $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/forward; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/layer ]; then ln -s $(PWD)/layer $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/layer; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/mux ]; then ln -s $(PWD)/mux $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/mux; fi;
	@if [ ! -d $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/router ]; then ln -s $(PWD)/router $(GOPATH)/src/gopkg.in/vinxi/vinxi.v0/router; fi;
