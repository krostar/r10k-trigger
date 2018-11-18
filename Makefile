# ex: /Users/alice/go/src/github.com/alice/project/
DIR_ABS    := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
DIR_BIN    := $(DIR_ABS)/build/bin

# use this rule as the default make rule
.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@# absolutely awesome -> http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# create directory in case they don't exists
$(DIR_BIN):
	@mkdir -p $@
$(DIR_COVER):
	@mkdir -p $@

.PHONY: clean
clean: ## Remove any rebuildable files
	-@$(RM) **/mock_*.go
	-@$(RM) -r $(DIR_BIN)

build-trigger-api: $(DIR_BIN)/trigger-api ## Build the trigger-api
.SECONDEXPANSION: # rebuild the project when any go files or if Makefile or go.sum is modified
$(DIR_BIN)/%: GO_FILES   := $(shell go list -f '{{ range $$f := .GoFiles }}{{ printf "%s/%s\n" $$.Dir $$f }}{{end}}' $(shell head -n1 go.mod | cut -d' ' -f2)/...)
$(DIR_BIN)/%: SH_FILES   := $(shell ls -1 $(DIR_ABS)/scripts/*.sh )
$(DIR_BIN)/%: $$(GO_FILES) $$(SH_FILES) Makefile go.sum | $(DIR_BIN)
	@DOCKER_RUN_ARGS=$(*) $(MAKE) ci-build-go

.PHONY: lint-all test-all
lint-all: ci-lint-go ci-lint-markdown ci-lint-yaml	## Run all possible linters
test-all: ci-test-go ci-test-go-deps 				## Run all possible tests

.PHONY: lint-go lint-markdown lint-yaml
lint-go: ci-lint-go				## Lint go files
lint-markdown: ci-lint-markdown	## Lint markdown files
lint-yaml: ci-lint-yaml			## Lint yaml files

.PHONY: test-go test-go-deps test-go-fast
test-go: ci-test-go				## Test go code
test-go-deps: ci-test-go-deps	## Test go dependencies
test-go-fast: ci-test-go-fast	## Fast test co code

.PHONY: ci-lint-% ci-test-%
ci-build-go: DOCKER_RUN_OPTS += --env BUILD_FOR_OS=$(shell go env GOOS)
ci-build-go: DOCKER_RUN_OPTS += --env BUILD_FOR_ARCH=$(shell go env GOARCH)
# ci-build-go: DOCKER_RUN_OPTS += --env BUILD_COMPRESS=1
# specify a special reusable volume for go-related docker builds
ci-%-go: DOCKER_RUN_OPTS += --mount type=volume,source='gomodcache',target='/go/pkg/mod/'
ci-%:
	@docker --log-level warn run \
		--rm \
		--tty \
		--mount type=bind,source="$(DIR_ABS)",target=/app \
		$(DOCKER_RUN_OPTS) \
		"krostar/ci:$(*)" \
		$(DOCKER_RUN_ARGS)
