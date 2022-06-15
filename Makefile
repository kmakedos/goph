## Before we start test that we have the mandatory executables available
RUNTIME = check_api
BUILD_ITEMS = check_api.go
EXECS = go docker
API_URL ?= http://localhost:8080
K := $(foreach exec,$(EXECS),\
$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH, consider installing $(exec)")))
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
export CGO_ENABLED=0


.PHONY: clean

.ONESHELL:
all:  docker

compile:
	go build $(BUILD_ITEMS)

docker: compile
	docker build -t $(RUNTIME) .

run: docker
	docker run -it -p 8000:8000  -e API_URL=$(API_URL) $(RUNTIME)

clean:
	docker rmi -f $(RUNTIME)
	rm -vf $(RUNTIME)