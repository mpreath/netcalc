/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    host.h - defines host structures and functions

    Copyright (C) 2009  Matthew Reath 

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#ifndef __HOST_H__
#define __HOST_H__
#include <glib.h>

typedef struct host_struct {
	guint32 ip_address; /* 32-bit IPv4 address */
	guint64 ipv6_address[2]; /* 128-bit IPv6 address */
	guint32 mask;
} host;

void initialize_host(host* h1, char* ip_address, char* mask);

#endif
