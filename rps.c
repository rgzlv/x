/* rock paper scissors */

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <time.h>

char *shapes[] = {"rock", "paper", "scissors"};

void
help(bool bad)
{
	size_t i;

	if (bad)
		fputs("bad choice.\n", stderr);
	fputs("pick a shape: ", bad ? stderr : stdout);
	for (i = 0; i < sizeof(shapes) / sizeof(shapes[0]); i++)
		fprintf(bad ? stderr : stdout, "%s ", shapes[i]);
	fputc('\n', bad ? stderr : stdout);
}

bool
beats(int a, int b)
{
	return a == (b + 1) % 3;
}

int
main(void)
{
	int shape_input;
	int shape_random;

	help(false);
	for (;;) {
		char choice[64];
		size_t sz;
		int i;

		if (!fgets(choice, sizeof(choice), stdin))
			return 0;

		sz = strlen(choice);
		if (choice[sz - 1] == '\n')
			choice[sz - 1] = 0;

		for (i = 0; i < sizeof(shapes) / sizeof(shapes[0]); i++)
			if (!strcmp(choice, shapes[i])) {
				shape_input = i;
				goto choice_ok;
			}

		help(true);
	}
choice_ok:

	srand(time(NULL));
	shape_random = rand() % 3;

	printf("computer: %s\n", shapes[shape_random]);
	if (shape_input == shape_random)
		puts("it's a draw");
	else if (beats(shape_input, shape_random))
		printf("%s beats %s, you win\n", shapes[shape_input], shapes[shape_random]);
	else
		printf("%s beats %s, you lose\n", shapes[shape_random], shapes[shape_input]);

	return 0;
}
