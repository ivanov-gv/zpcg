APP_NAME?=matcher
BUILD_DIR=build

lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -count=1 -race ./...

clean:
	rm -f ${BUILD_DIR}/${APP_NAME}

build: clean
	go build -o ${BUILD_DIR}/${APP_NAME} ./cmd/

run: build
	${BUILD_DIR}/${APP_NAME}