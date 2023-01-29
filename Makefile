DHT2PROMETHEUS_BIN := bin/dht2prometheus

GO ?= go

VERSION := $(shell cat VERSION)
USE_VENDOR =
LOCAL_LDFLAGS = -buildmode=pie -ldflags "-X=github.com/thkukuk/dht2prometheus/pkg/dht2prometheus.Version=$(VERSION)"

.PHONY: all api build vendor
all: dep build

dep: ## Get the dependencies
	@$(GO) get -v -d ./...

update: ## Get and update the dependencies
	@$(GO) get -v -d -u ./...

tidy: ## Clean up dependencies
	@$(GO) mod tidy

vendor: dep ## Create vendor directory
	@$(GO) mod vendor

build: ## Build the binary files
	$(GO) build -v -o $(DHT2PROMETHEUS_BIN) $(USE_VENDOR) $(LOCAL_LDFLAGS) ./cmd/dht2prometheus

clean: ## Remove previous builds
	@rm -f $(DHT2PROMETHEUS_BIN)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: release
release: ## create release package from git
	git clone https://github.com/thkukuk/dht2prometheus
	mv dht2prometheus dht2prometheus-$(VERSION)
	sed -i -e 's|USE_VENDOR =|USE_VENDOR = -mod vendor|g' dht2prometheus-$(VERSION)/Makefile
	make -C dht2prometheus-$(VERSION) vendor
	cp VERSION dht2prometheus-$(VERSION)
	tar --exclude .git -cJf dht2prometheus-$(VERSION).tar.xz dht2prometheus-$(VERSION)
	rm -rf dht2prometheus-$(VERSION)
