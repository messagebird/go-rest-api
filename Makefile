.PHONY: examples format

GORUN	= go run
GOFMT	= gofmt -d=true -s=true -w=true

all:
	@echo "make examples - Run all the out-of-the-box workable examples"
	@echo "make format   - Reformat the .go files in this project"

examples:
	$(GORUN) examples/balance.go
	$(GORUN) examples/hlr_create.go
	$(GORUN) examples/message_create.go
	$(GORUN) examples/voice_message_create.go

format:
	$(GOFMT) messagebird/messagebird.go 
	$(GOFMT) examples/*.go
