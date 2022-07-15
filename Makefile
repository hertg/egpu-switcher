BINDIR := /usr/bin
SHAREDIR := /usr/share
MANDIR := /usr/share/man/man1
BINARY_NAME := egpu-switcher

all:

build:
	go build -o ${BINARY_NAME}

install:
	mkdir -p ${DESTDIR}${BINDIR}
	cp egpu-switcher ${DESTDIR}${BINDIR}/
	mkdir -p ${DESTDIR}${SHAREDIR}/egpu-switcher
	cp xorg.conf.template ${DESTDIR}${SHAREDIR}/egpu-switcher/
	cp egpu.service ${DESTDIR}${SHAREDIR}/egpu-switcher/
	mkdir -p ${DESTDIR}${MANDIR}
	cp docs/egpu-switcher.1 ${DESTDIR}${MANDIR}/
	rm -f ${DESTDIR}${MANDIR}/egpu-switcher.1.gz
	gzip ${DESTDIR}${MANDIR}/egpu-switcher.1

uninstall:
	rm -f ${DESTDIR}${BINDIR}/egpu-switcher
	rm -rfd ${DESTDIR}${SHAREDIR}/egpu-switcher
	rm -f ${DESTDIR}${MANDIR}/egpu-switcher.1*

