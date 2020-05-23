BIN = gohttpd

default: build

build:
	go fmt
	go build

FBSD_PREFIX = /usr/local
FBSD_BIN_DEST = $(FBSD_PREFIX)/bin/$(BIN)
FBSD_RC_DEST = $(FBSD_PREFIX)/etc/rc.d/gohttpd

install-freebsd: build _install-freebsd

_install-freebsd:
	sudo cp $(BIN) $(FBSD_BIN_DEST)
	sudo cp ./contrib/gohttpd.rc.d $(FBSD_RC_DEST)
	sudo sysrc gohttpd_enable=YES
	@echo 'Set the domains config with:'
	@echo '    sudo sysrc gohttpd_domains="example.com,www.example.com"'

pre-freebsd:
	@sudo service gohttpd stop || echo -n

reinstall-freebsd: build pre-freebsd _install-freebsd
	sudo service gohttpd start
