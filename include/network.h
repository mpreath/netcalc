/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    network.h - network structure and function definitions

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


#ifndef __NETWORK_H__
#define __NETWORK_H__

#include <host.h>
#include <utility.h>
#include <glib.h>

typedef struct network_struct {

	host address;
	guint32 host_count;
	
} network;

/* initialize network */
void initiaize_network(network* n, host* h);

int print_network_info(network* n);

guint32 get_network_address(host* h);

guint32 get_broadcast_address(host* h);

int get_bits_in_mask(guint32 mask);

guint32 get_mask_from_bits(int bits);

guint32 extend_mask(guint32 mask, int bits);

guint32 shorten_mask(guint32 mask, int bits);

#endif
