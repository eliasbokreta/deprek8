BIN			= deprek8
PROJECT		= deprek8
BINFOLDER 	= /usr/local/bin
GO			= go
VERSION		= $(shell git describe --tags --abbrev=0 2>/dev/null)

LDFLAGS	= -ldflags="-X github.com/eliasbokreta/$(PROJECT)/cmd.version=$(VERSION)"

SRC		= $(shell find . -name "*.go")

.PHONY: tidy fmt fmt-check lint test build doc doc-check install uninstall clean

default: all

all: tidy fmt lint test build doc

tidy:
	$(info ▶ cleaning dependencies...)
	$(GO) mod tidy

fmt:
	$(info ▶ formatting...)
	gofmt -s -w .

fmt-check:
	$(info ▶ checking go format...)
ifneq ($(shell gofmt -d $(SRC) 2>/dev/null),)
	$(error Error, please format your code using 'make fmt')
endif

lint:
	$(info ▶ running lint tools...)
	golangci-lint run -v

test:
	$(info ▶ running tests...)
	$(GO) test -v -coverprofile cover.out ./...

build:
	$(info ▶ compiling program...)
	$(GO) build $(LDFLAGS)

doc: build
	$(info ▶ generating Cobra documentation...)
	@mkdir -p docs
	@rm -f docs/*
	./$(BIN) doc
	@mv ./docs/$(PROJECT).md ./docs/README.md

doc-check: doc
	$(info ▶ checking if documentation is up to date...)
ifneq ($(shell git status --porcelain docs/ 2>/dev/null),)
	$(error Error, please re generate the documentation)
endif

install: build
	$(info ▶ installing $(PROJECT)...)
	@mkdir -p ~/.$(PROJECT)
	@cp -r config/ ~/.$(PROJECT)
	sudo cp $(BIN) $(BINFOLDER)

uninstall:
	$(info ▶ uninstalling $(PROJECT)...)
	@rm -rf ~/.$(PROJECT)
	sudo rm $(BINFOLDER)/$(BIN)

clean:
	$(info ▶ removing $(PROJECT) binary...)
	rm $(BIN)
