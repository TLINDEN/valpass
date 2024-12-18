# Copyright © 2023 Thomas von Dein

# This  module  is published  under  the  terms  of the  BSD  3-Clause
# License. Please read the file LICENSE for details.

#
# no need to modify anything below

VERSION   = $(shell grep VERSION handler.go | head -1 | cut -d '"' -f2)

all: buildlocal

buildlocal:
	go build  -o example/example example/example.go

clean:
	rm -rf $(tool) coverage.out testdata t/out example/example

test: clean
	go test $(ARGS)

singletest:
	@echo "Call like this: make singletest TEST=TestName ARGS=-v"
	go test -run $(TEST) $(ARGS)

cover-report:
	go test -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

goupdate:
	go get -t -u=patch ./...

lint:
	golangci-lint run -p bugs -p unused

release: buildlocal test
	gh release create v$(VERSION) --generate-notes


bench:
	go test -bench=. -count 5
