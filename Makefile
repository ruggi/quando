GOLANGCI_VERSION = 1.32.0
GOLANGCI = .bin/golangci/$(GOLANGCI_VERSION)/golangci-lint

$(GOLANGCI):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $(GOLANGCI)) v$(GOLANGCI_VERSION)

.PHONY: lint
lint: $(GOLANGCI)
	$(GOLANGCI) run ./...

.PHONY: test
test:
	go test ./...
