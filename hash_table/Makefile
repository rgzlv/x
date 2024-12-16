warn := -Wall -Wextra -Wpedantic
nowarn := -Wno-unused-parameter -Wno-unused-variable
override CPPFLAGS := -std=c99 $(warn) $(nowarn) $(CPPFLAGS)

.PHONY: clean
ht: ht.c

clean:
	rm -f ht

