# Definitions for values
.PHONY: build clean tool lint help

PACKAGE_NAME := deepin-app-analyze
PACKAGE_VERSION := 0.0.1
PACKAGE_MAINTAINER := Zhou Fei <zhoufei@uniontech.com>
GOPKG_PREFIX := github.com/spf13/cobra
TEMP_FILE := /tmp/gocache
PREFIX := /usr/bin
PWD := $(shell pwd)
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
export GO111MODULE=auto


# Debian directory
DEBIAN_DIR := debian
DEBIAN_DIR_BASELINE := ${DEBIAN_DIR}/etc/${PACKAGE_NAME}/baseline
DEBIAN_DIR_SCRIPT := ${DEBIAN_DIR}${PREFIX}

# build directory
BUILD_DIR := ${PWD}/buildpkgtest/${PACKAGE_NAME}-${PACKAGE_VERSION}


# before build ready
# prepare: 
# 	@if [ ! -d ${GOPATH}/src/${GOPKG_PREFIX} ]; then \
# 		ln -sf ${PWD}/vendor/* ${GOPATH}/src/; \
# 	fi 


# Build your project binary main file 
build:  
	@echo "[INFO]: compile packages with vwndor"
	go build -ldflags "-s -w" -v -mod vendor -o ${PACKAGE_NAME} ${PWD}/main.go


# Build deepin-app-analyze application package
package: build
	cp -f ${PWD}/${PACKAGE_NAME} ${DEBIAN_DIR_SCRIPT}
	cp -r ${PWD}/config/* ${DEBIAN_DIR_BASELINE}
	mkdir -p ${BUILD_DIR}
	cp -r ${PWD}/${DEBIAN_DIR} ${BUILD_DIR}/


# Debian install functions
install: 
	mkdir -p ${DESTDIR}/etc/${PACKAGE_NAME}/baseline
	cp -f datas/${PACKAGE_NAME}/baseline/* ${DESTDIR}/etc/${PACKAGE_NAME}/baseline/
	cp -f datas/${PACKAGE_NAME}/usr/bin/* ${DESTDIR}${PREFIX}/


# Clean histroy build files
clean: 
	rm -fr ${BUILD_DIR}
	rm -rf ${DEBIAN_DIR_BASELINE}/*
	rm -rf ${DEBIAN_DIR_SCRIPT}/*
	rm -fr ${PACKAGE_NAME}

# Uninstall the package
uninstall: 
	dpkg -r ${PACKAGE_NAME}
