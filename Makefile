all: 
	go run cmd/antha/antha.go -componentlib=true -outdir=antha/component/lib antha/component/an
	gofmt -w -s antha/component/lib

.PHONY: all
