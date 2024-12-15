warn := -Wall -Wextra -Wpedantic
nowarn := -Wno-unused-parameter
override CPPFLAGS := -std=c99 $(warn) $(nowarn) $(CPPFLAGS)

.PHONY: clean

awc: awc.c

clean:
	rm -f awc