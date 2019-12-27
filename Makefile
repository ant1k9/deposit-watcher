.PHONY: all

commands = migrate update app backup rotate-backups

all: $(commands)

$(commands): %: cmd/%/main.go
	go build -o deposit-$@ $<

clean:
	@rm -f deposit-*
