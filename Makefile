IMAGE = opolis/config:dev
GOPATH = /go/src/github.com/opolis/config
BIN=config

RUN = docker run -it --rm \
	  -v $(HOME)/.aws:/root/.aws \
	  -v $(PWD):$(GOPATH) \
	  -v $(PWD)/.cache:/root/.cache/go-build \
	  -w $(GOPATH) \
	  $(IMAGE)

COMPILE = env GOOS=linux go build -ldflags="-s -w"

.PHONY: image
image:
	@docker build -t $(IMAGE) .

.PHONY: deps
deps:
	@$(RUN) dep ensure

.PHONY: build
build:
	@$(RUN) $(COMPILE) -o $(BIN) main.go

.PHONY: shell
shell:
	@$(RUN) sh
