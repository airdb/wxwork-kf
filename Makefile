SHELL = /bin/bash

# load env file
ifeq ($(shell test -e .env && echo -n yes),yes)
include .env
export $(shell sed 's/=.*//' .env)
endif

# generate build version info
VERSION:=$(shell git describe --dirty --always)
#VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse HEAD)
REPO := github.com/airdb/wxwork-kf

LDFLAGS=-ldflags
LDFLAGS += "-X=github.com/airdb/sailor/version.Repo=$(REPO) \
            -X=github.com/airdb/sailor/version.Version=$(VERSION) \
            -X=github.com/airdb/sailor/version.Build=$(BUILD) \
            -X=github.com/airdb/sailor/version.BuildTime=$(shell date +%s)"

default: swag build deploy

build:
	GOOS=linux go build $(LDFLAGS) -o main main.go

swag:
	swag init --generalInfo main.go

dev: swag
	env=dev go run $(LDFLAGS) main.go

deploy:
	sls deploy --stage test

release:
	sls deploy --stage release

log:
	sls logs --tail --stage test

.PHONY:
print-env:
	@env
