#include <sdmon-curses.h>

static struct winsize sdmon_ws;
static struct sigaction sdmon_sa;
_Bool sdmon_resize_needed = false;

void sdmon_end_curses() {
	endwin();
}

// Updates LINES and COLS
void sdmon_update_lc(void) {
	if (ioctl(STDOUT_FILENO, TIOCGWINSZ, &sdmon_ws) == -1) {
		sdmon_end_curses(); fprintf(stderr, "ERROR: sdmon_update_lc()\n"); exit(EXIT_FAILURE); }
	LINES=sdmon_ws.ws_row;
	COLS=sdmon_ws.ws_col;
}

void sdmon_handle_int_err(int ret) {
	if(ret == ERR) {
	sdmon_end_curses(); fprintf(stderr, "ERROR: sdmon_handle_int()\n"); exit(EXIT_FAILURE); }
}

void sdmon_handle_null(void *p) {
	if (p == NULL) {
		fprintf(stderr, "ERROR: sdmon_handle_null()\n"); exit(EXIT_FAILURE); }
}

void sdmon_handle_null_end(void *p) {
	if (p == NULL) {
		sdmon_end_curses(); fprintf(stderr, "ERROR: sdmon_handle_null()\n"); exit(EXIT_FAILURE); }
}

void sdmon_resize(WINDOW **pad, int pad_rows) {
	sdmon_resize_needed = false;
	sdmon_update_lc();
	resizeterm(LINES, COLS);
	wresize(*pad, pad_rows, COLS);
	prefresh(*pad, 0, 0, 0, 0, LINES-1, COLS-1);
}

void sdmon_handle_sigwinch(int sig) {
	sdmon_resize_needed = true;
}

void sdmon_init_curses(WINDOW **scr, WINDOW **pad, int pad_rows) {
	// if "", use from environment
	sdmon_handle_null(setlocale(LC_ALL, ""));

	sdmon_handle_null( *scr = initscr() );
	sdmon_update_lc();
	sdmon_handle_int_err(start_color());
	sdmon_handle_int_err(init_pair(1, COLOR_WHITE, COLOR_BLACK));
	sdmon_handle_int_err(init_pair(2, COLOR_BLUE, COLOR_BLACK));
	sdmon_handle_int_err(init_pair(3, COLOR_RED, COLOR_BLACK));
	sdmon_handle_int_err(init_pair(4, COLOR_GREEN, COLOR_BLACK));
	sdmon_handle_int_err(init_pair(5, COLOR_YELLOW, COLOR_BLACK));
	sdmon_handle_int_err(cbreak());
	sdmon_handle_int_err(noecho());

	sdmon_handle_null_end( *pad = newpad(pad_rows, COLS) );
	sdmon_handle_int_err(keypad(*pad, TRUE));
	sdmon_handle_int_err(wbkgd(*pad, COLOR_PAIR(1)));
	sdmon_handle_int_err(refresh());
	/* pminrow - line/row of the top left edge (in pad)
	 * pmincol - column of the top left edge (in pad)
	 * sminrow - line/row of the top left edge (in terminal)
	 * smincol - column of the top left edge (in terminal)
	 * smaxrow - line/row of the bottom right edge (in terminal)
	 * smaxcol - column of the bottom right edge (in terminal) */
	sdmon_handle_int_err(prefresh(*pad, 0, 0, 0, 0, LINES-1, COLS-1));

	// Handle resize
	sdmon_sa.sa_handler = sdmon_handle_sigwinch;
	if ( sigaction(SIGWINCH, &sdmon_sa, NULL) == -1 ) {
		fprintf(stderr, "ERROR: sigaction()\n"); exit(EXIT_FAILURE); }
}


void sdmon_wprintw_red(WINDOW **pad, char *str) {
	wbkgdset(*pad, COLOR_PAIR(3));
	sdmon_handle_int_err(wattron(*pad, A_BOLD));
	sdmon_handle_int_err(wprintw(*pad, str));
	sdmon_handle_int_err(wattroff(*pad, A_BOLD));
	wbkgdset(*pad, COLOR_PAIR(1));
}

void sdmon_wprintw_green(WINDOW **pad, char *str) {
	wbkgdset(*pad, COLOR_PAIR(4));
	sdmon_handle_int_err(wattron(*pad, A_BOLD));
	sdmon_handle_int_err(wprintw(*pad, str));
	sdmon_handle_int_err(wattroff(*pad, A_BOLD));
	wbkgdset(*pad, COLOR_PAIR(1));
}

void sdmon_wprintw_yellow(WINDOW **pad, char *str) {
	wbkgdset(*pad, COLOR_PAIR(5));
	sdmon_handle_int_err(wattron(*pad, A_BOLD));
	sdmon_handle_int_err(wprintw(*pad, str));
	sdmon_handle_int_err(wattroff(*pad, A_BOLD));
	wbkgdset(*pad, COLOR_PAIR(1));
}

void sdmon_display_LoadState(WINDOW **pad, char *str) {
	if (!strcmp(str, "loaded"))
		sdmon_wprintw_green(pad, str);
	else
		sdmon_wprintw_red(pad, str);
}

void sdmon_display_ActiveState(WINDOW **pad, char *str) {
	if (!strcmp(str, "active"))
		sdmon_wprintw_green(pad, str);
	else if (!strcmp(str, "reloading") || !strcmp(str, "activating") || !strcmp(str, "deactivating"))
		sdmon_wprintw_yellow(pad, str);
	else
		sdmon_wprintw_red(pad, str);
}

void sdmon_display_SubState(WINDOW **pad, char *str) {
	if (!strcmp(str, "running"))
		sdmon_wprintw_green(pad, str);
	else if (!strcmp(str, "exited"))
		sdmon_wprintw_red(pad, str);
	else
		wprintw(*pad, str);
}

void sdmon_display_UnitFileState(WINDOW **pad, char *str) {
	if (!strcmp(str, "enabled") || !strcmp(str, "static"))
		sdmon_wprintw_green(pad, str);
	else if (!strcmp(str, "enabled-runtime") || !strcmp(str, "linked") || !strcmp(str, "linked-runtime") || !strcmp(str, "masked-runtime"))
		sdmon_wprintw_yellow(pad, str);
	else
		sdmon_wprintw_red(pad, str);
}

