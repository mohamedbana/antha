CWD=$(shell pwd)
LIBDIR=antha/component/lib
ANDIR=antha/component/an
LIBPACKAGE=$(shell basename "$(LIBDIR)")
ANTHA_GO=$(GOPATH)/src/github.com/antha-lang/antha/cmd/antha/antha.go

all: go_from_antha

go_from_antha:
	cd "$(LIBDIR)" && \
		find "$(CWD)/$(ANDIR)" -name '*.an' -print0 | \
		xargs -0 go run "$(ANTHA_GO)" -componentLib=true -componentLibPackage="$(LIBPACKAGE)" && \
		gofmt -w -s .

clean:
	find "$(LIBDIR)" -type d -depth 1 -print0 | xargs -0 rm -r
	rm -f "$(LIBDIR)/$(LIBPACKAGE).go"

.PHONY: all go_from_antha clean
