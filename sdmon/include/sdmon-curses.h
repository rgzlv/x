#ifndef SDMON_CURSES_H
#define SDMON_CURSES_H
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <locale.h>

#include <curses.h>
#include <sys/ioctl.h>

#include <curses.h>

extern _Bool sdmon_resize_needed;

extern void sdmon_end_curses();
extern void sdmon_update_lc(void);
extern void sdmon_handle_int_err(int ret);
extern void sdmon_handle_null(void *p);
extern void sdmon_handle_null_end(void *p);
extern void sdmon_resize(WINDOW **pad, int pad_rows);
extern void sdmon_handle_sigwinch(int sig);
extern void sdmon_init_curses(WINDOW **scr, WINDOW **pad, int pad_rows);
extern void sdmon_wprintw_red(WINDOW **pad, char *str);
extern void sdmon_wprintw_green(WINDOW **pad, char *str);
extern void sdmon_wprintw_yellow(WINDOW **pad, char *str);
extern void sdmon_display_LoadState(WINDOW **pad, char *str);
extern void sdmon_display_ActiveState(WINDOW **pad, char *str);
extern void sdmon_display_SubState(WINDOW **pad, char *str);
extern void sdmon_display_UnitFileState(WINDOW **pad, char *str);

#endif // SDMON_CURSES_H

