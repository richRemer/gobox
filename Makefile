name = gobox
bin = bin
prefix = /usr/local

default: build

build: $(bin)/$(name)

clean:
	rm -rf $(bin)

install: $(bin)/$(name)
	cp $< $(prefix)/bin/$(name)

$(bin)/$(name): $(wildcard *.go)
	@mkdir -p $(@D)
	go build -o $@

.PHONY: default build clean install
