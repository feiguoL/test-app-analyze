.PHONY: build

# Definitions for values
PREFIX := /usr
PACKAGE_NAME := deepin-app-analyze
PACKAGE_VERSION := 0.0.1
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
	mkdir -pv ${BUILD_DIR}${PREFIX}/local/bin
	mkdir -pv ${BUILD_DIR}/etc/${PACKAGE_NAME}/baseline/

# Build your project binary main file 
build: prepare
	go build -ldflags "-s -w" -v -mod vendor -o ${BUILD_DIR}${PREFIX}/local/bin/${PACKAGE_NAME}

# Debian install functions
install: 
	cp -r ${PWD}/config/* ${BUILD_DIR}/etc/${PACKAGE_NAME}/baseline/
	mkdir -pv ${DESTDIR}/${PREFIX}/bin/
	install -v ${BUILD_DIR}${PREFIX}/local/bin/${PACKAGE_NAME} ${DESTDIR}/${PREFIX}/bin/${PACKAGE_NAME}
	mkdir -pv ${DESTDIR}/etc/${PACKAGE_NAME}/baseline
	install -Dm644 ${PWD}/config/* ${DESTDIR}/etc/${PACKAGE_NAME}/baseline/


# Debian uninstall package
uninstall: 
	rm -f ${PREFIX}/bin/${PACKAGE_NAME}
	rm -fr ${DESTDIR}/etc/${PACKAGE_NAME}

# Clean histroy build files
clean:
	rm -fr ${BUILD_ROOT}
