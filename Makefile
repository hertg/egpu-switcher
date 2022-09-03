BINDIR := /usr/bin
SHAREDIR := /usr/share
MANDIR := /usr/share/man/man1
OUT_DIR := ./bin
BINARY_NAME := egpu-switcher

all: build

build: test
	go build -o ${OUT_DIR}/${BINARY_NAME}

gendocs: build
	${OUT_DIR}/${BINARY_NAME} gendocs

test:
	go test ./...

install: build manpage
	mkdir -p ${DESTDIR}${BINDIR}
	cp ${BINARY_NAME} ${DESTDIR}${BINDIR}/
	mkdir -p ${DESTDIR}${SHAREDIR}/egpu-switcher
	mkdir -p ${DESTDIR}${MANDIR}
	cp docs/man/egpu-switcher*.1 ${DESTDIR}${MANDIR}/
	rm -f ${DESTDIR}${MANDIR}/egpu-switcher*.1.gz
	gzip ${DESTDIR}${MANDIR}/egpu-switcher*.1

uninstall:
	rm -f ${DESTDIR}${BINDIR}/egpu-switcher
	rm -rfd ${DESTDIR}${SHAREDIR}/egpu-switcher
	rm -f ${DESTDIR}${MANDIR}/egpu-switcher.1*

