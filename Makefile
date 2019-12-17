.PHONY: all

commands = migrate update app backup

all: $(commands)

$(commands): %: cmd/%/main.go
	go build -o deposit-$@ $<
