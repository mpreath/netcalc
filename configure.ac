#                                               -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ(2.60)
AC_INIT(netcalc, 0.7, mpreath at gmail dot com)
AC_CONFIG_SRCDIR([src/netcalc.c])
AM_INIT_AUTOMAKE([-Wall -Werror])
AM_CONFIG_HEADER([config.h])

# Checks for programs.
AC_PROG_CC

# Checks for libraries.
PKG_CHECK_MODULES([GLIB], [glib-2.0 >= 2.0.0])
AC_SUBST([GLIB_CFLAGS])
AC_SUBST([GLIB_LIBS])

# Checks for header files.
AC_HEADER_STDC
AC_CHECK_HEADERS([stdlib.h string.h])
AC_CHECK_HEADERS([math.h])
AC_SEARCH_LIBS(cos,m)

#AC_CHECK_HEADERS([glib.h])
#AC_SEARCH_LIBS(g_free,glib-2.0)

# Checks for typedefs, structures, and compiler characteristics.

# Checks for library functions.
AC_FUNC_MALLOC

AC_CONFIG_FILES([Makefile
                 src/Makefile])
AC_OUTPUT
