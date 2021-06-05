VERSION_MAJOR = 1
VERSION_MINOR = 1
VERSION_PATCH = 0

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

.PHONY: publish
publish: 
	git checkout main
	git pull
	@git diff --exit-code || { git status ; echo You have uncommitted changes. ; false ; };
	$(eval dots = $(subst ., ,$(VERSION)))
	$(eval new_major = $(word 1, $(dots)))
	$(eval new_minor = $(word 2, $(dots)))
	$(eval new_patch = $(word 3, $(dots)))
	sed -i.bak -e 's/^\(VERSION_MAJOR = \).*/\1$(new_major)/g' Makefile
	sed -i.bak -e 's/^\(VERSION_MINOR = \).*/\1$(new_minor)/g' Makefile
	sed -i.bak -e 's/^\(VERSION_PATCH = \).*/\1$(new_patch)/g' Makefile
	rm Makefile.bak

	git commit -am 'v$(VERSION)'
	git tag v$(VERSION)
	git push --follow-tags
	git push origin v$(VERSION)

.PHONY: publish-major
publish-major:
	@make publish VERSION=$$(($(VERSION_MAJOR) + 1)).0.0

.PHONY: publish-minor
publish-minor:
	@make publish VERSION=$(VERSION_MAJOR).$$(($(VERSION_MINOR) + 1)).0

.PHONY: publish-patch
publish-patch:
	@make publish VERSION=$(VERSION_MAJOR).$(VERSION_MINOR).$$(($(VERSION_PATCH) + 1))
