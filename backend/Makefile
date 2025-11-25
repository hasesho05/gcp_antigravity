.PHONY: run generate-sample test

run:
	/usr/local/go/bin/go run cmd/api/main.go

generate-types:
	cd internal/domain && /usr/local/go/bin/go run ../../cmd/gen_enum_methods/main.go
	PATH=/usr/local/go/bin:$(PATH) $(HOME)/go/bin/tygo generate
	/usr/local/go/bin/go run cmd/gen_ts_enums/main.go

test:
	/usr/local/go/bin/go test ./...
