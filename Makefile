all: gen_comp

gen_comp:
	go run cmd/antha/antha.go -outdir=antha/component/lib antha/component/an
	gofmt -w -s antha/component/lib

test:
	go test -v `go list ./... | grep -v internal | grep -v bvendor`

fmt_json:
	for i in `find antha/examples -name '*.json' -o -name '*.yml'`; do \
	  python -mjson.tool "$$i" > "$$i.bak" && mv "$$i.bak" "$$i"; \
	done

compile:
	go install github.com/antha-lang/antha/cmd/...

test_workflows: compile
	for d in `find antha/examples -type d -o -name '*.yml'`; do \
	  if [[ -f "$$d/workflow.json" && -f "$$d/parameters.yml" ]]; then \
	    (cd "$$d" && antharun --workflow workflow.json --parameters parameters.yml > /dev/null); \
	    if [[ $$? == 0 ]]; then \
	      echo "OK $$d"; \
	    else \
	      echo "FAIL $$d"; \
	    fi; \
	  fi; \
	done

.PHONY: all gen_comp fmt_json test test_workflows compile
