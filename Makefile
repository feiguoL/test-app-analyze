.PHONY: build

# Definitions for values
PREFIX := /usr
PACKAGE_NAME := deepin-app-analyze
PACKAGE_BINARY := app-analyze
PACKAGE_VERSION := $(shell scripts/version.sh)
PACKAGE_MAINTAINER := Zhou Fei <zhoufei@uniontech.com>
GOPKG_PREFIX := github.com/spf13/cobra
TEMP_FILE := /tmp/gocache

PWD := $(shell pwd)
GOBASE := $(shell pwd)
GoPath := GOPATH=${GOBASE}:$(GOBASE)/vendor:${GOPATH}
export GO111MODULE=auto

# build directory
BUILD_ROOT := ${PWD}/buildpkgtest
BUILD_DIR := ${BUILD_ROOT}/${PACKAGE_NAME}-${PACKAGE_VERSION}

# build pre
prepare:
	mkdir -pv ${BUILD_DIR}${PREFIX}/bin/
	mkdir -pv ${BUILD_DIR}/etc/${PACKAGE_NAME}/baseline/
	cp -r ${PWD}/config/* ${BUILD_DIR}/etc/${PACKAGE_NAME}/baseline/
	cp -f ${PWD}/scripts/${PACKAGE_NAME} ${BUILD_DIR}${PREFIX}/bin/

# Build your project binary main file 
build: prepare
	go build -ldflags "-s -w" -v -mod vendor -o ${BUILD_DIR}${PREFIX}/bin/${PACKAGE_BINARY}

# Debian install functions
install: 
	mkdir -p ${DESTDIR}${PREFIX}/bin/
	install -Dm755 ${BUILD_DIR}${PREFIX}/bin/${PACKAGE_BINARY} ${DESTDIR}${PREFIX}/bin/
	install -Dm755 ${BUILD_DIR}${PREFIX}/bin/${PACKAGE_NAME} ${DESTDIR}${PREFIX}/bin/
	mkdir -p ${DESTDIR}/etc/${PACKAGE_NAME}/baseline
	install -Dm644 ${BUILD_DIR}/etc/${PACKAGE_NAME}/baseline/* ${DESTDIR}/etc/${PACKAGE_NAME}/baseline/


# Debian uninstall package
uninstall: 
	rm -f ${PREFIX}/bin/${PACKAGE_NAME}
	rm -f ${PREFIX}/bin/${PACKAGE_BINARY}
	rm -fr ${DESTDIR}/etc/${PACKAGE_NAME}

# Clean histroy build files
clean:
	rm -fr ${BUILD_ROOT}
