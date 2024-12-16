CC := gcc
# -Wconversion for implicit conversions that may alter a value
# -Wc++-compat for implicit cast of void * to non-void *
CFLAGS := -std=c99 -Wall -Wextra -Wpedantic -Wno-unused-variable -Wno-unused-parameter -Wconversion -Wc++-compat -I/usr/include
LDFLAGS := 

.PHONY: clean
.DEFAULT_TARGET := wlc

wlc: main.o xdg-shell.o
	$(CC) -o $@ $(CFLAGS) $^ $(LDFLAGS) -lwayland-client

main.o: main.c xdg-shell.h
	$(CC) -c -o $@ $(CFLAGS) $<

xdg-shell.o: xdg-shell.c xdg-shell.h
	$(CC) -c -o $@ $(CFLAGS) $<

xdg-shell.h: xdg-shell.xml
	wayland-scanner client-header < $< > $@

xdg-shell.c: xdg-shell.xml
	wayland-scanner private-code < $< > $@

clean:
	rm -rf wlc main.o xdg-shell.o xdg-shell.h xdg-shell.c
