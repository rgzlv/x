#ifndef SBUS_H
#define SBUS_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <dbus/dbus.h>

void sbus_check_error(const char *func_name, DBusError *error);
void sbus_check_msg_is_null(const char *func_name, DBusMessage *check_msg);
void sbus_check_bool_not_true(const char *func_name, dbus_bool_t check_bool);
void sbus_check_bool_is_false(const char *func_name, dbus_bool_t check_bool);
void sbus_check_msg_is_reply(const char *func_name, DBusMessage *check_msg);
void sbus_check_arg_type_is_string(const char *func_name, int iter_arg_type);
void sbus_check_arg_type_is_variant(const char *func_name, int iter_arg_type);
void sbus_new_method_msg_s(DBusMessage **msg, const char *dest, const char *path, const char *iface, const char *method, char *arg1);
void sbus_new_method_msg_ss(DBusMessage **msg, const char *dest, const char *path, const char *iface, const char *method, char *arg1, char *arg2);
void sbus_send_msg(DBusConnection *con, DBusMessage *msg, DBusMessage **reply);
DBusBasicValue sbus_get_msg_arg_s(DBusMessage **msg);
DBusBasicValue sbus_get_msg_arg_s_in_v(DBusMessage **msg);
char *sbus_GetUnit(DBusConnection **con, char *unit);
char *sbus_LoadState(DBusConnection **con, char *unit);
char *sbus_ActiveState(DBusConnection **con, char *unit);
char *sbus_SubState(DBusConnection **con, char *unit);
char *sbus_UnitFileState(DBusConnection **con, char *unit);
#endif // SBUS_H
