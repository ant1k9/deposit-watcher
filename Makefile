.PHONY: all

commands = migrate update app backup rotate-backups

all: $(commands)

$(commands): %: cmd/%/main.go
	go build -o deposit-$@ $<

lint:
	golangci-lint run

clean:
	@rm -f deposit-*
