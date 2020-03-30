TEST?=./...

# Acceptance Testing
testacc:
	TF_ACC=1 go test -v $(TEST) $(TESTARGS)

# Build the provider
build:
	go build -o terraform-provider-muleb2b

default: build
