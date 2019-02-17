BINDIR := /usr/bin
SHAREDIR := /usr/share

all:

install:
	mkdir -p ${DESTDIR}${BINDIR}
	cp egpu-switcher ${DESTDIR}${BINDIR}
	mkdir -p ${DESTDIR}${SHAREDIR}/egpu-switcher
	cp xorg.conf.template ${DESTDIR}${SHAREDIR}/egpu-switcher
	cp egpu-switcher-startup /etc/init.d/
	ln -sf /etc/init.d/egpu-switcher-startup /etc/rc5.d/S15egpu-switcher

uninstall:
	egpu-switcher cleanup
	rm -rfd ${DESTDIR}${SHAREDIR}/egpu-switcher
	rm /etc/rc5.d/S15egpu-switcher
	rm /etc/init.d/egpu-switcher-startup
