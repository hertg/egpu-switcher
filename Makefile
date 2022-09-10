BINDIR := /usr/bin
SHAREDIR := /usr/share
MANDIR := /usr/share/man/man1

DOCS_DIR := ./docs
OUT_DIR := ./bin
BINARY_NAME := egpu-switcher
OUT_BIN := ${OUT_DIR}/${BINARY_NAME}

all: build

build:
	go build -o ${OUT_BIN}
	@echo "binary compiled => ${OUT_BIN}"
	${OUT_BIN} gendocs -o ${DOCS_DIR}
	@echo "docs generated => ${DOCS_DIR}"

clean:
	rm -f ${OUT_BIN}
	rm -fd ${OUT_DIR}
	rm -rfd ${DOCS_DIR}
	@echo "cleanup successful"

test:
	go test ./...

install:
	@if [ ! -f ${OUT_BIN} ]; then\
		echo "run 'build' command first";\
		exit 1;\
	fi
	mkdir -p ${DESTDIR}${BINDIR}
	cp ${OUT_BIN} ${DESTDIR}${BINDIR}/
	@echo "binary installed at ${DESTDIR}${BINDIR}/${BINARY_NAME}"
	mkdir -p ${DESTDIR}${MANDIR}
	cp docs/man/egpu-switcher*.1 ${DESTDIR}${MANDIR}/
	rm -f ${DESTDIR}${MANDIR}/egpu-switcher*.1.gz
	gzip ${DESTDIR}${MANDIR}/egpu-switcher*.1
	@echo "manpages installed in ${DESTDIR}${MANDIR}"

uninstall:
	rm -f ${DESTDIR}${BINDIR}/egpu-switcher
	@echo "removed binary at ${DESTDIR}${BINDIR}/egpu-switcher"
	rm -f ${DESTDIR}${MANDIR}/egpu-switcher*.1.gz
	@echo "removed manpages at ${DESTDIR}${MANDIR}/egpu-switcher*.1.gz"

