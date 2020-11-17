.PHONY: build
build: 
	go build internal/cmd/main.go

.PHONY: vet
vet:
	go vet ./...
