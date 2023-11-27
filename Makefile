BINARY_NAME=gross
INSTALL_DIR=/usr/local/bin

build:
	go build -buildvcs=false -o ${BINARY_NAME}

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}	
	rm ${BINARY_NAME}
	rm ${BINARY_NAME}

install: build
	cp ./${BINARY_NAME} ${INSTALL_DIR}

