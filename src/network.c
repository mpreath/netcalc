/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    network.c - network related functions

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

#include <network.h>
#include <host.h>
#include <utility.h>
#include <stdio.h>

int initialize_network(network* n, host* h) {

	unsigned int s, e;

	n->address.ip_address = get_network_address(h);
	n->address.mask = h->mask;

	s = n->address.ip_address;
	e = get_broadcast_address(&n->address);

	n->host_count = (e - s) - 1;

	return 1;
}

int print_network_info(network* n) {

	unsigned int i, s, e;
	int bc = 0;	

	s = n->address.ip_address;
	e = get_broadcast_address(&n->address);	

	bc = get_bits_in_mask(n->address.mask);
	
	char ip_address[16];
	char mask[16];
	inttodd(ip_address, s);
	printf("Network:\t%s\n", ip_address);
	inttodd(ip_address, e);
	printf("Broadcast:\t%s\n", ip_address);
	inttodd(ip_address, n->address.mask);
	printf("Mask:\t\t%s[/%i]\n", ip_address, bc);
	printf("Hosts:\t\t%i\n", n->host_count);
	
	inttodd(mask, n->address.mask);

	for(i = s+1; i < e; i++) {
		inttodd(ip_address,i);
		printf("%s\t%s\n", ip_address, mask);
	}
	

	return 1;
}

unsigned int get_network_address(host* h) {

	return (h->ip_address & h->mask);
}

unsigned int get_broadcast_address(host* h) {

	return (h->ip_address | ~(h->mask));
}

int get_bits_in_mask(unsigned int mask) {

	int bc;
	
	for(bc = 0; mask != 0; mask = mask << 1, bc++)
		;
	return bc;
}

unsigned int get_mask_from_bits(int bits) {

	unsigned int mask;
	int i;
	int bit_count;

	bit_count = 32 - bits;
	mask = 0;
	mask = ~mask;

	for(i = 0; i < bit_count; mask = mask << 1, i++)
		;
	

	return mask;

}

unsigned int extend_mask(unsigned int mask, int bits) {

	int i;
	unsigned int new_mask, m;
	new_mask = mask;
	m = 1;
	m = m << 31;
	// shit to the right one bit, then set the high order bit to 1
	for(i = 0; i < bits; i++) {
		new_mask = new_mask >> 1;
		new_mask = new_mask | m;
	}

	return new_mask;

}

unsigned int shorten_mask(unsigned int mask, int bits) {

	int i;

	for(i = 0; i < bits; i++, mask = mask << 1)
		;
	
	return mask;

}
