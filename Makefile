DIST_DIR := dist
APP_NAME := ai
.PHONY: build

BUILD_OPT := -v
build: 
	mkdir -p ${DIST_DIR}
	go build ${BUILD_OPT} -o ${DIST_DIR}/${APP_NAME}
