
test-build:
	go build -tags testing

test: test-build
	bats t