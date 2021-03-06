/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    network.c - network related functions

    Copyright (c) 2021  Matthew Reath 

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

#include <network.h>
#include <host.h>
#include <utility.h>
#include <stdio.h>
#include <glib.h>

void initialize_network(network *n, host *h)
{

	guint32 t1, t2;

	t1 = 0;
	t2 = 0;

	guint32 s, e;

	n->address.ip_address = get_network_address(h);
	n->address.mask = h->mask;

	is_valid_mask(n->address.mask);

	s = n->address.ip_address;
	e = get_broadcast_address(&n->address);

	n->host_count = (e - s) - 1;
}

int print_network_info(network *n, gboolean verbose)
{

	guint32 i, s, e;
	int bc = 0;

	s = n->address.ip_address;
	e = get_broadcast_address(&n->address);

	bc = get_bits_in_mask(n->address.mask);

	char ip_address[16];
	char mask[16];
	if(verbose) 
	{
		inttodd(ip_address, s);
		printf("Network:\t%s\n", ip_address);
		inttodd(ip_address, e);
		printf("Broadcast:\t%s\n", ip_address);
		inttodd(ip_address, n->address.mask);
		printf("Mask:\t\t%s[/%i]\n", ip_address, bc);
		printf("Hosts:\t\t%i\n", n->host_count);
	}

	inttodd(mask, n->address.mask);

	for (i = s + 1; i < e; i++)
	{
		inttodd(ip_address, i);
		if(verbose)
			printf("%s\t%s\n", ip_address, mask);
		else
			printf("%s %s\n", ip_address, mask);
	}

	return 1;
}

guint32 get_network_address(host *h)
{

	return (h->ip_address & h->mask);
}

guint32 get_broadcast_address(host *h)
{

	return (h->ip_address | ~(h->mask));
}

int get_bits_in_mask(guint32 mask)
{

	int bc;

	for (bc = 0; mask != 0; mask = mask << 1, bc++)
		;
	return bc;
}

guint32 get_mask_from_bits(int bits)
{

	guint32 mask;
	int i;
	int bit_count;

	bit_count = 32 - bits;
	mask = 0;
	mask = ~mask;

	for (i = 0; i < bit_count; mask = mask << 1, i++)
		;

	return mask;
}

guint32 extend_mask(guint32 mask, int bits)
{

	int i;
	guint32 new_mask, m;
	new_mask = mask;
	m = 1;
	m = m << 31;
	// shift to the right one bit, then set the high order bit to 1
	for (i = 0; i < bits; i++)
	{
		new_mask = new_mask >> 1;
		new_mask = new_mask | m;
	}

	return new_mask;
}

guint32 shorten_mask(guint32 mask, int bits)
{

	int i;

	for (i = 0; i < bits; i++, mask = mask << 1)
		;

	return mask;
}
