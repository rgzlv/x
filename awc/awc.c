#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>

typedef struct bin_tree {
	const char *word;
	int len;
	size_t count;
	struct bin_tree *left, *right;
} bin_tree;

void err(const char *s) {
	if (s) fputs(s, stderr);
	if (errno) {
		if (s)
			fprintf(stderr, ": %s\n", strerror(errno));
		else
			fputs(strerror(errno), stderr);
	}
	if (s || errno)
		putc('\n', stderr);
	else
		fprintf(stderr, "unknown error\n");
	exit(EXIT_FAILURE);
}

void print_word_counts(bin_tree *t, size_t *total) {
	if (t->word) {
		printf("word: '%.*s', count: %zu\n", t->len, t->word, t->count);
		*total += t->count;
	}
	if (t->left) print_word_counts(t->left, total);
	if (t->right) print_word_counts(t->right, total);
}

void print_tree(bin_tree *t, int off) {
	char *prefix = malloc(off + 1);
	memset(prefix, ' ', off);
	prefix[off] = '\0';
	if (t->word) printf("%sword: %.*s (len: %d, count: %zu)\n", prefix, t->len, t->word, t->len, t->count);
	if (t->left) {
		printf("%sleft:\n", prefix);
		print_tree(t->left, off + 2);
	}
	if (t->right) {
		printf("%sright:\n", prefix);
		print_tree(t->right, off + 2);
	}
	free(prefix);
}

void insert_word(bin_tree *t, const char *word, size_t len) {
	if (!t->word) {
		t->word = word;
		t->count++;
		t->len = len;
		return;
	}
	int cmp = strncmp(t->word, word, len);
	if (!cmp) {
		t->count++;
	} else if (cmp < 0) {
		if (!t->right) {
			t->right = calloc(1, sizeof(*t->right));
		}
		insert_word(t->right, word, len);
	} else if (cmp > 0) {
		if (!t->left) {
			t->left = calloc(1, sizeof(*t->left));
		}
		insert_word(t->left, word, len);
	}
}

int main(int argc, char **argv) {
	if (argc != 2) err("expected 1 filename argument");
	int fd = open(argv[1], O_RDONLY);
	if (fd == -1) err("couldn't open file");
	struct stat st;
	if (fstat(fd, &st)) err(NULL);
	char *buf = malloc(st.st_size + 1);
	if (read(fd, buf, st.st_size) == -1) err(NULL);
	buf[st.st_size] = '\0';

	bin_tree t = {0};
	char *beg = buf, *end;
	for (;;) {
		end = strchr(beg, ' ');
		size_t len;
		if (!end) {
			len = strlen(beg);
			if (!len) break;
			insert_word(&t, beg, len);
			break;
		}
		len = end - beg;
		if (!len) break;
		insert_word(&t, beg, len);
		beg = end + 1;
	}

	size_t total = 0;
	print_word_counts(&t, &total);
	printf("total: %zu\n", total);

	free(buf);
	close(fd);
	return EXIT_SUCCESS;
}
