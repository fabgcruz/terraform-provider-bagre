BINARY=terraform-provider-bagre

default: build

build:
	go build -o $(BINARY) .

vet:
	go vet ./...

fmt:
	gofmt -s -w .

# Acceptance tests run real plan/apply against a live Bagre instance.
# Requires BAGRE_ENDPOINT and BAGRE_TOKEN to be set.
testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

.PHONY: default build vet fmt testacc
