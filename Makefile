BINDIR := /usr/bin
SHAREDIR := /usr/share

all:

install:
	mkdir -p ${DESTDIR}${BINDIR}
	cp egpu-switcher ${DESTDIR}${BINDIR}
	mkdir -p ${DESTDIR}${SHAREDIR}/egpu-switcher
	cp xorg.conf.template ${DESTDIR}${SHAREDIR}/egpu-switcher
	cp egpu.service ${DESTDIR}${SHAREDIR}/egpu-switcher

uninstall:
	rm -rfd ${DESTDIR}${BINDIR}
	rm -rfd ${DESTDIR}${SHAREDIR}/egpu-switcher
	