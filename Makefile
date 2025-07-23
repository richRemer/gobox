bin = gobox
prefix = /usr/local

default: build

build: $(bin)

clean:
	rm -rf $(bin)

install: $(bin)
	cp $< $(prefix)/bin/$(name)

$(bin): $(wildcard *.go)
	@mkdir -p $(@D)
	go build -o $@

.PHONY: default build clean install
