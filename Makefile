# Makefile for gitlab group manager
# vim: set ft=make ts=8 noet
# Copyright Yakshaving.art
# Licence MIT

# Variables
# UNAME		:= $(shell uname -s)

COMMIT_ID := `git log -1 --format=%H`
COMMIT_DATE := `git log -1 --format=%aI`
VERSION := $${CI_COMMIT_TAG:-SNAPSHOT-$(COMMIT_ID)}
SHELL := /bin/bash

GOOS ?= linux
GOARCH ?= amd64

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:	### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

# Targets
#
.PHONY: debug
debug:	### debug Makefile itself
	@echo $(UNAME)

.PHONY: check
check:	### sanity checks
	@find . -type f \( -name \*.yml -o -name \*yaml \) \! -path './vendor/*' \
		| xargs -r yq '.' # >/dev/null

.PHONY: lint
lint:	check
lint:	### run all the lints
	gometalinter

.PHONY: test
test:	### run all the unit tests
# test: lint
	@go test -v -coverprofile=coverage.out $$(go list ./... | grep -v '/vendor/') \
		&& go tool cover -func=coverage.out

.PHONY: integration
integration: ### run integration tests (requires a bootstrapped local environment)
	@go test -v ./... -tags "integration" -coverprofile=coverage.out $$(go list ./... | grep -v '/vendor/') \
		&& go tool cover -func=coverage.out

.PHONY: build
build: ### build the binary applying the correct version from git
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X \
		gitlab.com/yakshaving.art/alertsnitch/version.Version=$(VERSION) -X \
		gitlab.com/yakshaving.art/alertsnitch/version.Commit=$(COMMIT_ID) -X \
		gitlab.com/yakshaving.art/alertsnitch/version.Date=$(COMMIT_DATE)" \
		-o alertsnitch-$(GOARCH)

CURRENT_DIR:=$(shell pwd)

.PHONY: bootstrap_local_testing
bootstrap_local_testing: ### builds and bootstraps a local integration testing environment using docker-compose
	@if [[ -z "$(MYSQL_ROOT_PASSWORD)" ]]; then echo "MYSQL_ROOT_PASSWORD is not set" ; exit 1; fi
	@if [[ -z "$(MYSQL_DATABASE)" ]]; then echo "MYSQL_DATABASE is not set" ; exit 1; fi
	@echo "Launching alertsnitch-mysql integration container"
	@docker run --rm --name alertsnitch-mysql \
		-e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) \
		-e MYSQL_DATABASE=$(MYSQL_DATABASE) \
		-p 3306:3306 \
		-v $(CURRENT_DIR)/db.d/mysql:/db.scripts \
		-d \
		mysql:5.7
	@while ! docker exec alertsnitch-mysql mysql --database=$(MYSQL_DATABASE) --password=$(MYSQL_ROOT_PASSWORD) -e "SELECT 1" >/dev/null 2>&1 ; do \
    echo "Waiting for database connection..." ; \
    sleep 1 ; \
	done
	@echo "Bootstrapping model"
	@docker exec alertsnitch-mysql sh -c "exec mysql -uroot -p$(MYSQL_ROOT_PASSWORD) $(MYSQL_DATABASE) < /db.scripts/0.0.1-bootstrap.sql"
	@docker exec alertsnitch-mysql sh -c "exec mysql -uroot -p$(MYSQL_ROOT_PASSWORD) $(MYSQL_DATABASE) < /db.scripts/0.1.0-fingerprint.sql"
	@echo "Everything is ready to run 'make integration'; remember to teardown_local_testing when you are done"

.PHONY: teardown_local_testing
teardown_local_testing: ### Tears down the integration testing environment
	docker stop alertsnitch-mysql
