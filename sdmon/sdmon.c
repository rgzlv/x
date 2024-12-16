#include <stdio.h>

#include <dbus/dbus.h>
#include <curses.h>

#include <sbus.h>
#include <sdmon-curses.h>

int main(int argc, char *argv[]) {
	if (argc < 2) {
		fprintf(stderr, "No args\n"); exit(EXIT_FAILURE); }

	int ch;
	int pad_rows = 201;
	int beg_row = 0;

	struct unit
	{
	    char *name;
	    char *serviceName;
	};
	
	struct unit units[argc-1];

	DBusConnection *con;
	DBusError err;

	WINDOW *scr, *pad;

	dbus_error_init(&err);
	sdmon_handle_null( (con = dbus_bus_get(DBUS_BUS_SYSTEM, &err)) );
	sbus_check_error(__func__, &err);


	for (int i = 0; i < argc-1; i++)
	{
		units[i].name = argv[i+1];
		if ( (units[i].serviceName = malloc(strlen(units[i].name) + strlen(".service") + 1)) == NULL ) {
			fprintf(stderr, "ERROR: units[i].serviceName malloc()\n"); exit(EXIT_FAILURE); }
		strcpy(units[i].serviceName, units[i].name);
		strcat(units[i].serviceName, ".service");
	}

	sdmon_init_curses(&scr, &pad, pad_rows);

	wtimeout(pad, 0);
	while ( (ch = wgetch(pad)) != 'q' )
	{
		if (sdmon_resize_needed)
			sdmon_resize(&pad, pad_rows);

		int display_row = 0;
		int display_col = 0;
		for (int i = 0; i < argc-1; i++)
		{
			if ( (i + 1) % 2 == 0 ) // Every 2nd from units[]
				display_col = COLS/2;
			else
				display_col = 0;
			if ( i % 2 == 0 && i != 0 ) // Every 3rd from units[]
				display_row += 6;

			wbkgdset(pad, COLOR_PAIR(2));
			sdmon_handle_int_err(wattron(pad, A_BOLD));
			sdmon_handle_int_err(mvwprintw(pad, display_row, display_col, "%s" ,units[i].name));
			sdmon_handle_int_err(wattroff(pad, A_BOLD));
			wbkgdset(pad, COLOR_PAIR(1));

			sdmon_handle_int_err(mvwprintw(pad, display_row+1, display_col, "LoadState: "));
			sdmon_display_LoadState(&pad, sbus_LoadState(&con, units[i].serviceName));

			sdmon_handle_int_err(mvwprintw(pad, display_row+2, display_col, "ActiveState: "));
			sdmon_display_ActiveState(&pad, sbus_ActiveState(&con, units[i].serviceName));

			sdmon_handle_int_err(mvwprintw(pad, display_row+3, display_col, "SubState: "));
			sdmon_display_SubState(&pad, sbus_SubState(&con, units[i].serviceName));

			sdmon_handle_int_err(mvwprintw(pad, display_row+4, display_col, "UnitFileState: "));
			sdmon_display_UnitFileState(&pad, sbus_UnitFileState(&con, units[i].serviceName));
		}
		sdmon_handle_int_err(prefresh(pad, beg_row, 0, 0, 0, LINES-1, COLS-1));

		if (ch == KEY_DOWN)
			if (beg_row < pad_rows-1)
				beg_row++;

		if (ch == KEY_UP)
			if (beg_row > 0)
				beg_row--;

		sdmon_handle_int_err(prefresh(pad, beg_row, 0, 0, 0, LINES-1, COLS-1));

		wtimeout(pad, 10);
	}

	sdmon_end_curses();

	dbus_connection_flush(con);
	dbus_connection_unref(con);
	
	dbus_shutdown();

	for (int i = 0; i < argc-1; i++) {
		free(units[i].serviceName);
	}

	return EXIT_SUCCESS;
}

