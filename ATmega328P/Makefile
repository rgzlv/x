CC_AVR := avr-gcc

CFLAGS_COMMON := -std=c99 -Wall -Wextra -Wpedantic

CFLAGS_AVR_COMMON := $(CFLAGS_AVR_COMMON) -mmcu=avr5 -O3
CFLAGS_AVR_REL := $(CFLAGS_AVR_COMMON)
CFLAGS_AVR_DEBUG := $(CFLAGS_AVR_COMMON) -g
AVR_BUILD ?= REL
CFLAGS_AVR := ${CFLAGS_AVR_${AVR_BUILD}}

AVRCONF := ./avrdude.conf
BIN ?= build/bin/t1
PORT ?= /dev/cu.usbmodem11401

FILES_MAIN := src/main.s
OBJS_MAIN := $(patsubst src/%.s,build/src/%.s.o,$(FILES_MAIN))
OBJS_MAIN += $(patsubst src/%.c,build/src/%.c.o,$(FILES_MAIN))
OBJS_MAIN := $(filter-out src/%,$(OBJS_MAIN))
DIRS_MAIN := $(sort $(dir $(OBJS_MAIN)))
CFLAGS_MAIN := $(CFLAGS_AVR)

FILES_UTIL := util/fsize/main.c
OBJS_UTIL := $(patsubst util/%.s,build/util/%.s.o,$(FILES_UTIL))
OBJS_UTIL += $(patsubst util/%.c,build/util/%.c.o,$(FILES_UTIL))
OBJS_UTIL := $(filter-out util/%,$(OBJS_UTIL))
DIRS_UTIL := $(sort $(dir $(OBJS_UTIL)))
CFLAGS_UTIL := $(CFLAGS_COMMON)

.PHONY: all main util send test dirs clean demo1
.DEFAULT_GOAL := all

all: main util | dirs


main: $(BIN)

$(BIN): build/src/main.s.o | dirs
	$(CC_AVR) $(CFLAGS_AVR) -o $@ $^

build/src/%.s.o: src/%.s | dirs
	$(CC_AVR) $(CFLAGS_AVR) -c -o $@ $<

build/src/%.c.o: src/%.c | dirs
	$(CC_AVR) $(CFLAGS_AVR) -c -o $@ $<


util: build/bin/fsize

build/bin/fsize: $(OBJS_UTIL) | dirs
	$(CC) $(CFLAGS_COMMON) -o $@ $^

build/util/%.c.o: util/%.c | dirs
	$(CC) $(CFLAGS_COMMON) -c -o $@ $<

build/util/%.s.o: util/%.s | dirs
	$(CC) $(CFLAGS_COMMON) -c -o $@ $<


demo1: build/bin/demo1

build/bin/demo1: build/demo/1/main.s.o
	$(CC_AVR) $(CFLAGS_AVR) -o $@ $^

build/demo/%.c.o: demo/%.c | dirs
	$(CC_AVR) $(CFLAGS_AVR) -c -o $@ $<

build/demo/%.s.o: demo/%.s | dirs
	$(CC_AVR) $(CFLAGS_AVR) -c -o $@ $<


dirs:
	mkdir -p build/{bin,src,util,demo}
	mkdir -p build/util/fsize
	mkdir -p build/demo/1

send: main
	avrdude -U flash:w:$(BIN):e -p m328p -c arduino -P $(PORT) -C $(AVRCONF)

clean:
	rm -rf build
