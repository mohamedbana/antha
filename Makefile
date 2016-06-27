all: gen_comp

gen_comp:
	go run cmd/antha/antha.go -outdir=antha/component/lib antha/component/an
	gofmt -w -s antha/component/lib

test:
	go test -v `go list ./... | grep -v vendor | grep -v bvendor`

gen_pb:
	go generate github.com/antha-lang/antha/driver
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|proto "([^"]*)"|proto "github.com/antha-lang/antha/bvendor/\1"|'
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|context "([^"]*)"|context "github.com/antha-lang/antha/bvendor/\1"|'
	find driver/pb -name '*.pb.go' | xargs perl -p -i -e 's|grpc "([^"]*)"|grpc "github.com/antha-lang/antha/bvendor/\1"|'

fmt_json:
	for i in `find antha/examples -name '*.json' -o -name '*.yml'`; do \
	  python -mjson.tool "$$i" > "$$i.bak" && mv "$$i.bak" "$$i"; \
	done

compile:
	go install github.com/antha-lang/antha/cmd/...

test_workflows: compile
	for d in `find antha/examples -mindepth 2 -maxdepth 2 -type d`; do \
	  abs=`cd $$d; pwd` ;\
	  (go test github.com/antha-lang/antha/antha/component/lib -args $$abs) ;\
	  if [[ $$? == 0 ]]; then \
	    echo "PASS $$d"; \
	  else \
	    echo "FAIL $$d"; \
	  fi; \
	done

.PHONY: all gen_comp fmt_json test test_workflows compile
