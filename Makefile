PROTOC_C ?= protoc-c
PKG_CONFIG ?= pkg-config

CFLAGS	?= -O2 -g -ggdb -pipe -Wall
LDFLAGS ?= -Wl,-z,relro

CFLAGS  += -std=c99 -DGOWHATSAPP_PLUGIN_VERSION='"$(PLUGIN_VERSION)"' -DMARKDOWN_PIDGIN

GIT ?= git
CC ?= gcc
GO ?= go
export GOPATH?=$(CURDIR)/go
export GO111MODULE?=auto

# Note: Use "-C .git" to avoid ascending to parent dirs if .git not present
GIT_REVISION_ID = $(shell $(GIT) -C .git rev-parse --short HEAD 2>/dev/null)
PLUGIN_VERSION ?= $(shell cat VERSION)~git$(GIT_REVISION_ID)

ifeq ($(shell $(PKG_CONFIG) --exists purple 2>/dev/null && echo "true"),)
  TARGET = FAILNOPURPLE
  DEST =
else
  TARGET = libgowhatsapp.so
  DEST = $(DESTDIR)`$(PKG_CONFIG) --variable=plugindir purple`
  ICONS_DEST = $(DESTDIR)`$(PKG_CONFIG) --variable=datadir purple`/pixmaps/pidgin/protocols
  LOCALEDIR = $(DESTDIR)$(shell $(PKG_CONFIG) --variable=datadir purple)/locale
endif

CFLAGS += -DLOCALEDIR=\"$(LOCALEDIR)\"

PURPLE_COMPAT_FILES := purple2compat/http.c purple2compat/purple-socket.c

.PHONY:	all FAILNOPURPLE clean update-dep gdb install

LOCALES = $(patsubst %.po, %.mo, $(wildcard po/*.po))

all: $(TARGET)

GO_WHATSAPP_A = $(GOPATH)/pkg/$(shell $(GO) env GOOS)_$(shell $(GO) env GOARCH)/github.com/Rhymen/go-whatsapp.a
GO_WHATSAPP_GIT = $(GOPATH)/src/github.com/Rhymen/go-whatsapp/.git

update-dep:
	$(GO) get -u github.com/skip2/go-qrcode
	$(GO) get -u github.com/gabriel-vasile/mimetype
	$(GO) get -u github.com/Rhymen/go-whatsapp
	touch -d "$(shell $(GIT) --git-dir="$(GO_WHATSAPP_GIT)" log -1 --date=rfc --format=%cd)" $(GO_WHATSAPP_A)

# TODO: add targets for all go dependencies
$(GO_WHATSAPP_A):
	$(GO) get github.com/Rhymen/go-whatsapp

purplegwa.a: purplegwa.go purplegwa-media.go $(GO_WHATSAPP_A) gwa-to-purple.o
	$(GO) get github.com/skip2/go-qrcode
	$(GO) get github.com/gabriel-vasile/mimetype
	$(GO) build -buildmode=c-archive -o purplegwa.a purplegwa.go purplegwa-media.go

gwa-to-purple.o: gwa-to-purple.c constants.h
	$(CC) $(CFLAGS) -c -o $@ gwa-to-purple.c

$(TARGET): libgowhatsapp.c $(PURPLE_COMPAT_FILES) purplegwa.a constants.h
	$(CC) $(CFLAGS) -fPIC -shared -o $@ libgowhatsapp.c purplegwa.a $(PURPLE_COMPAT_FILES) $(LDFLAGS) `$(PKG_CONFIG) purple glib-2.0 --libs --cflags` $(INCLUDES) -Ipurple2compat

FAILNOPURPLE:
	echo "You need libpurple development headers installed to be able to compile this plugin"

clean:
	rm -f $(TARGET) purplegwa.a gwa-to-purple.o

install: $(TARGET) install-icons
	mkdir -m 0755 -p $(DEST)
	install -m 0755 -p $(TARGET) $(DEST)

install-icons: icons/16/whatsapp.png icons/22/whatsapp.png icons/48/whatsapp.png
	mkdir -m 0755 -p $(ICONS_DEST)/16
	mkdir -m 0755 -p $(ICONS_DEST)/22
	mkdir -m 0755 -p $(ICONS_DEST)/48
	install -m 0644 -p icons/16/whatsapp.png $(ICONS_DEST)/16/whatsapp.png
	install -m 0644 -p icons/22/whatsapp.png $(ICONS_DEST)/22/whatsapp.png
	install -m 0644 -p icons/48/whatsapp.png $(ICONS_DEST)/48/whatsapp.png

gdb:
	gdb --args pidgin -d -c .pidgin
