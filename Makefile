.PHONY: all

commands = migrate update app backup rotate-backups

all: $(commands)

$(commands): %: cmd/%/main.go
	go build -o deposit-$@ $<

lint:
	golangci-lint run

load:
	./deposit-update

dump:
	./deposit-backup

clean:
	@rm -f deposit-*
