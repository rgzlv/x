CFLAGS=-std=c99 -Wall -Wextra -Wpedantic -Wno-unused-parameter -Wno-unused-variable -Wno-unused-but-set-variable -I./include/
CFLAGS_DBG=-g -O0
CFLAGS_DBUS=$(shell pkg-config --cflags dbus-1)
CFLAGS_NCURSES=$(shell pkg-config --cflags ncurses)
LDFLAGS=-L./lib/
LDFLAGS_DBUS=$(shell pkg-config --libs dbus-1)
LDFLAGS_NCURSES=$(shell pkg-config --libs ncurses)

.PHONY := clean dbg

sdmon: sdmon.c sbus.o sdmon-curses.o
	$(CC) $(CFLAGS) $(CFLAGS_DBUS) $(CFLAGS_NCURSES) -c -o sdmon.o sdmon.c
	$(CC) $(CFLAGS) $(CFLAGS_DBUS) $(CFLAGS_NCURSES) -o sdmon sdmon.o sbus.o sdmon-curses.o $(LDFLAGS) $(LDFLAGS_DBUS) $(LDFLAGS_NCURSES)

sbus.o: ./include/sbus.h ./lib/sbus.c
	$(CC) $(CFLAGS) $(CFLAGS_DBUS) -c -o sbus.o ./lib/sbus.c

sdmon-curses.o: ./include/sdmon-curses.h ./lib/sdmon-curses.c
	$(CC) $(CFLAGS) $(CFLAGS_NCURSES) -c -o sdmon-curses.o ./lib/sdmon-curses.c

clean:
	rm -f sdmon sdmon.o sbus.o sdmon-curses.o
