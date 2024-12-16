#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define KR_HASHSIZE 800
// #define KR_HASHSIZE 101

typedef struct kr_nlist {
  struct kr_nlist *next;
  char *sym, *val;
} kr_nlist;

static kr_nlist *kr_hashtab[KR_HASHSIZE];

// hashval = *s + hashval;
// Would produce the same hash for "ab" and "ba".
// 
unsigned int kr_hash(const char *s) {
  unsigned int hashval = 0;
  for (;*s; s++)
    hashval = *s + 2 * hashval;
    // hashval = *s + 31 * hashval;
  return hashval;
  // return hashval % KR_HASHSIZE;
}

int main(int argc, char **argv) {
  printf("%d\n", kr_hash("ab"));
  printf("%d\n", kr_hash("ba"));
  return EXIT_SUCCESS;
}

