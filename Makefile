BINDIR := /usr/bin

all:

install:
	mkdir -p ${DESTDIR}${BINDIR}
	cp egpu-switcher ${DESTDIR}${BINDIR}
