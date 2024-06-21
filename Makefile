GOOS_LINUX = linux
GOOS_WINDOWS = windows
GOOS_MACOS = darwin

GOARCH_LINUX ?= amd64
GOARCH_WINDOWS ?= amd64
GOARCH_MACOS = arm64
GOARCH_MACOS_INTEL = amd64

CURR_TIME = $(shell date +%s)

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[36m  _             _   _ \n | |           | | | |\n | | _____  ___| |_| |\n | |/ / __|/ __| __| |\n |   <\\__ \\ (__| |_| |\n |_|\\_\\___/\\___|\\__|_| \033[0m\n\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Generate
.PHONY: gen-docs
gen-docs: ## Generates docs
	go run gen_docs.go

##@ Install (Dev)


.PHONY: install_linux_mock
install_linux_mock:  ## Install ksctl
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_LINUX} GOARCH=${GOARCH_LINUX} ./builder-mock.sh

.PHONY: install_macos_mock
install_macos_mock: ## Install ksctl on macos m1,m2,..
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_MACOS} GOARCH=${GOARCH_MACOS} ./builder-mock.sh

.PHONY: install_macos_intel_mock
install_macos_intel_mock: ## Install ksctl on macos intel
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_MACOS} GOARCH=${GOARCH_MACOS_INTEL} ./builder-mock.sh

.PHONY: install_linux
install_linux:  ## Install ksctl
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_LINUX} GOARCH=${GOARCH_LINUX} ./builder.sh

.PHONY: install_macos
install_macos: ## Install ksctl on macos m1,m2,..
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_MACOS} GOARCH=${GOARCH_MACOS} ./builder.sh

.PHONY: install_macos_intel
install_macos_intel: ## Install ksctl on macos intel
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_MACOS} GOARCH=${GOARCH_MACOS_INTEL} ./builder.sh

.PHONY: uninstall
uninstall:  ## Uninstall ksctl
	@echo "Started to Uninstall ksctl"
	cd scripts && \
		./uninstall.sh

##@ Linters
.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
		$(GOLANGCI_LINT) run && echo -e "\n=========\n\033[91m✔ PASSED\033[0m\n=========\n" || echo -e "\n=========\n\033[91m✖ FAILED\033[0m\n=========\n"


##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= /tmp/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)

## Tool Versions
GOLANGCI_LINT_VERSION ?= v1.54.2

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
