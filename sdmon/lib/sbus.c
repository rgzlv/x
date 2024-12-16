#include <sbus.h>
/*
 * Error checking
 */

void sbus_check_error(const char *func_name, DBusError *error) {
	if (dbus_error_is_set(error))
	{
		fprintf(stderr, "Error in %s->sbus_check_error, %s\n", func_name, error->message); dbus_error_free(error); exit(EXIT_FAILURE);
	}
}

void sbus_check_msg_is_null(const char *func_name, DBusMessage *check_msg) {
	if (check_msg == NULL) { fprintf(stderr, "Error in %s->sbus_check_msg_is_null\n", func_name); exit(EXIT_FAILURE); }
}

void sbus_check_bool_not_true(const char *func_name, dbus_bool_t check_bool) {
	if (check_bool != TRUE) { fprintf(stderr, "Error in %s->sbus_check_bool_not_true\n", func_name); exit(EXIT_FAILURE); }
}

void sbus_check_bool_is_false(const char *func_name, dbus_bool_t check_bool) {
	if (check_bool == FALSE) { fprintf(stderr, "Error in %s->sbus_check_bool_is_false\n", func_name); exit(EXIT_FAILURE); }
}

// Could be an error message either from dbus_connection_send_with_reply() if timeout expires or remote
void sbus_check_msg_is_reply(const char *func_name, DBusMessage *check_msg) {
	if (dbus_message_get_type(check_msg) != DBUS_MESSAGE_TYPE_METHOD_RETURN)
	{
		fprintf(stderr, "Error in %s->sbus_check_msg_is_reply\n", func_name); exit(EXIT_FAILURE);
	}
}

void sbus_check_arg_type_is_string(const char *func_name, int iter_arg_type) {
	// I think these are all the string-like types
	// They were described here at the bottom of this section
	// https://dbus.freedesktop.org/doc/dbus-specification.html#basic-types
	if (iter_arg_type != DBUS_TYPE_STRING && iter_arg_type != DBUS_TYPE_OBJECT_PATH && iter_arg_type != DBUS_TYPE_SIGNATURE)
	{
		fprintf(stderr, "Error in %s->sbus_check_arg_type_is_string\n", func_name); exit(EXIT_FAILURE);
	}
}

void sbus_check_arg_type_is_variant(const char *func_name, int iter_arg_type) {
	if (iter_arg_type != DBUS_TYPE_VARIANT) { fprintf(stderr, "Error in %s->sbus_check_arg_type_is_variant\n", func_name); exit(EXIT_FAILURE); }
}

/*
 * Creating messages
 */

// Create new method message with 1 string arg appended
void sbus_new_method_msg_s(DBusMessage **msg, const char *dest, const char *path, const char *iface, const char *method, char *arg1) {
	dbus_bool_t ret_bool;

	*msg = dbus_message_new_method_call(dest, path, iface, method);
	sbus_check_msg_is_null(__func__, *msg);

	ret_bool = dbus_message_append_args(*msg,
				DBUS_TYPE_STRING, &arg1,
				DBUS_TYPE_INVALID);
	sbus_check_bool_not_true(__func__, ret_bool);
}

// Create new method message with 2 string args appended
void sbus_new_method_msg_ss(DBusMessage **msg, const char *dest, const char *path, const char *iface, const char *method, char *arg1, char *arg2) {
	dbus_bool_t ret_bool;

	*msg = dbus_message_new_method_call(dest, path, iface, method);
	sbus_check_msg_is_null(__func__, *msg);

	ret_bool = dbus_message_append_args(*msg,
				DBUS_TYPE_STRING, &arg1,
				DBUS_TYPE_STRING, &arg2,
				DBUS_TYPE_INVALID);
	sbus_check_bool_not_true(__func__, ret_bool);
}

/*
 * Sending
 */

// Send a message and return reply (blocking)
void sbus_send_msg(DBusConnection *con, DBusMessage *msg, DBusMessage **reply) {
	dbus_bool_t ret_bool;
	DBusPendingCall *pending;
	DBusMessage *tmp_reply;

	ret_bool = dbus_connection_send_with_reply(con, msg, &pending, DBUS_TIMEOUT_USE_DEFAULT);
	sbus_check_bool_not_true(__func__, ret_bool);

	dbus_pending_call_block(pending);

	tmp_reply = dbus_pending_call_steal_reply(pending);
	sbus_check_msg_is_null(__func__, tmp_reply);
	sbus_check_msg_is_reply(__func__, tmp_reply);
	*reply = tmp_reply;

	dbus_pending_call_unref(pending);
}

/*
 * Getting arg values from messages/iterators
 */

// Return string argument from msg using the DBusBasicValue union
DBusBasicValue sbus_get_msg_arg_s(DBusMessage **msg) {
	dbus_bool_t ret_bool;
	DBusMessageIter iter;
	DBusBasicValue value;

	ret_bool = dbus_message_iter_init(*msg, &iter);
	sbus_check_bool_is_false(__func__, ret_bool);

	sbus_check_arg_type_is_string(__func__, dbus_message_iter_get_arg_type(&iter));
	dbus_message_iter_get_basic(&iter, &value);

	return value;
}

// Return string from msg by recursing into a variant
DBusBasicValue sbus_get_msg_arg_s_in_v(DBusMessage **msg) {
	dbus_bool_t ret_bool;
	DBusMessageIter iter, iter_sub;
	DBusBasicValue value;
	
	ret_bool = dbus_message_iter_init(*msg, &iter);
	sbus_check_bool_is_false(__func__, ret_bool);

	sbus_check_arg_type_is_variant(__func__, dbus_message_iter_get_arg_type(&iter));
	dbus_message_iter_recurse(&iter, &iter_sub);
	sbus_check_arg_type_is_string(__func__, dbus_message_iter_get_arg_type(&iter_sub));
	dbus_message_iter_get_basic(&iter_sub, &value);

	return value;
}

/*
 * High level, systemd specific probably
 */

// Return object path (string) for unit
char *sbus_GetUnit(DBusConnection **con, char *unit) {
	DBusMessage *msg, *reply;
	DBusBasicValue reply_value;
	char *reply_value_str;

	sbus_new_method_msg_s(&msg,
				"org.freedesktop.systemd1",
				"/org/freedesktop/systemd1",
				"org.freedesktop.systemd1.Manager",
				"GetUnit",
				unit);
	sbus_send_msg(*con, msg, &reply);
	reply_value = sbus_get_msg_arg_s(&reply);
	reply_value_str = reply_value.str;

	dbus_message_unref(msg);
	dbus_message_unref(reply);

	return reply_value_str;
}

// https://www.freedesktop.org/wiki/Software/systemd/dbus/
char *sbus_LoadState(DBusConnection **con, char *unit) {
	DBusMessage *msg, *reply;
	DBusBasicValue reply_value;
	char *unit_path;
	char *reply_value_str;

	unit_path = sbus_GetUnit(con, unit);
	sbus_new_method_msg_ss(&msg,
				"org.freedesktop.systemd1",
				unit_path,
				"org.freedesktop.DBus.Properties",
				"Get",
				"org.freedesktop.systemd1.Unit", "LoadState");
	sbus_send_msg(*con, msg, &reply);
	reply_value = sbus_get_msg_arg_s_in_v(&reply);
	reply_value_str = reply_value.str;

	dbus_message_unref(msg);
	dbus_message_unref(reply);

	return reply_value.str;
}

// https://www.freedesktop.org/wiki/Software/systemd/dbus/
char *sbus_ActiveState(DBusConnection **con, char *unit) {
	DBusMessage *msg, *reply;
	DBusBasicValue reply_value;
	char *unit_path;
	char *reply_value_str;

	unit_path = sbus_GetUnit(con, unit);
	sbus_new_method_msg_ss(&msg,
				"org.freedesktop.systemd1",
				unit_path,
				"org.freedesktop.DBus.Properties",
				"Get",
				"org.freedesktop.systemd1.Unit", "ActiveState");
	sbus_send_msg(*con, msg, &reply);
	reply_value = sbus_get_msg_arg_s_in_v(&reply);
	reply_value_str = reply_value.str;

	dbus_message_unref(msg);
	dbus_message_unref(reply);

	return reply_value_str;
}

// https://www.freedesktop.org/wiki/Software/systemd/dbus/
char *sbus_SubState(DBusConnection **con, char *unit) {
	DBusMessage *msg, *reply;
	DBusBasicValue reply_value;
	char *unit_path;
	char *reply_value_str;

	unit_path = sbus_GetUnit(con, unit);
	sbus_new_method_msg_ss(&msg,
				"org.freedesktop.systemd1",
				unit_path,
				"org.freedesktop.DBus.Properties",
				"Get",
				"org.freedesktop.systemd1.Unit", "SubState");
	sbus_send_msg(*con, msg, &reply);
	reply_value = sbus_get_msg_arg_s_in_v(&reply);
	reply_value_str = reply_value.str;

	dbus_message_unref(msg);
	dbus_message_unref(reply);

	return reply_value_str;
}

// https://www.freedesktop.org/wiki/Software/systemd/dbus/
char *sbus_UnitFileState(DBusConnection **con, char *unit) {
	DBusMessage *msg, *reply;
	DBusBasicValue reply_value;
	char *unit_path;
	char *reply_value_str;

	unit_path = sbus_GetUnit(con, unit);
	sbus_new_method_msg_ss(&msg,
				"org.freedesktop.systemd1",
				unit_path,
				"org.freedesktop.DBus.Properties",
				"Get",
				"org.freedesktop.systemd1.Unit", "UnitFileState");
	sbus_send_msg(*con, msg, &reply);
	reply_value = sbus_get_msg_arg_s_in_v(&reply);
	reply_value_str = reply_value.str;

	dbus_message_unref(msg);
	dbus_message_unref(reply);

	return reply_value_str;
}

