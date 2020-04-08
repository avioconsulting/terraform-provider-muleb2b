TEST?=./...
VERSION?="v0.0.0"
ENV?="DEV"

# Acceptance Testing
testacc:
	TEST_ENV_NAME=$(ENV) TF_ACC=1 go test -v $(TEST) $(TESTARGS)

# Build the provider
build: clean
	go build -o terraform-provider-muleb2b_$(VERSION)

package: clean build
	tar -czf "terraform-provider-muleb2b_${GOOS}_${GOARCH}.tar.gz" terraform-provider-muleb2b_$(VERSION)

default: build

clean:
	rm -f terraform-provider-muleb2b*