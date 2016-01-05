all:  gen_comp

gen_comp:
	go run cmd/antha/antha.go -outdir=antha/component/lib antha/component/an
	gofmt -w -s antha/component/lib

fmt_json:
	for i in `find antha/examples -name '*.json' -o -name '*.yml'`; do python -mjson.tool "$$i" > "$$i.bak" && mv "$$i.bak" "$$i"; done

test:
	go test -v `go list ./... | grep -v internal | grep -v bvendor`

.PHONY: all test gen_comp fmt_json
