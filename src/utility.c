/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    utility.c - misc functions used by netcalc 

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

#include <utility.h>
#include <network.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <glib.h>

guint32 ddtoint(char *dd)
{

	guint32 ip;
	ip = 0;

	guint32 octets[4];

	int i = 0;
	int octet = 0;
	char *tok;
	char *sep = ".";

	// need to verify input size
	// max string size = 255.255.255.255 = 15 characters

	if (strlen(dd) > 15)
		g_error("address length is too long");

	/* replaced old code with a loop using strtok, much simpler */
	for (tok = strtok(dd, sep), i = 1; tok; tok = strtok(NULL, sep), octet++, i++)
	{
		//printf("%i\n", strlen(tok));
		// verify that the input is numerical

		if (!is_number(tok))
			g_error("address contains non-numeric values");

		// verify input is the correct value, put in array
		if ((octets[octet] = atoi(tok)) > 255)
			g_error("octet in address is too large (255 max)");
	}

	//printf("i=%i\n", i);

	// verify our loop count, it tells us if we have an address
	// with less than 4 octets or more than 4 octets
	if (i != 5)
		g_error("address does not have the correct number of octets");

	//printf("Initial:\t %i\n", ip);

	ip = ip | octets[3];

	//printf("First octet:\t %i\n", ip);

	ip = ip | (octets[2] << 8);

	//printf("Second octet:\t %i\n", ip);

	ip = ip | (octets[1] << 16);

	//printf("Third octet:\t %i\n", ip);

	ip = ip | (octets[0] << 24);

	//printf("Fourth octet:\t %i\n", ip);

	return ip;
}

guint32 cidrtoint(char* cidr_mask) {

	guint32 mask_length;
	guint32 mask = 4294967295;
	

	if (is_number(cidr_mask)) 
	{
		mask_length = atoi(cidr_mask);
		for(int i = 0; i < 32-mask_length; i++)
		{
			mask = mask << 1;			
		}

		return mask;

	}
	
	return -1;
}

int inttodd(char *dd, guint32 ip)
{

	guint32 octets[4];
	octets[0] = 0;
	octets[1] = 0;
	octets[2] = 0;
	octets[3] = 0;

	// parse out the octets int an array
	octets[3] = ip << 24 >> 24;

	//printf("[3] %i\n", octets[3]);

	octets[2] = ip << 16 >> 24;

	//printf("[2] %i\n", octets[2]);

	octets[1] = ip << 8 >> 24;

	//printf("[1] %i\n", octets[1]);

	octets[0] = ip >> 24;

	//printf("[0] %i\n", octets[0]);

	sprintf(dd, "%i.%i.%i.%i", octets[0], octets[1], octets[2], octets[3]);

	return 1;
}

int is_number(char *s)
{

	int i;
	int b = 1;

	/* ONLY WORKS FOR ASCII codes */

	for (i = 0; i < strlen(s); i++)
	{
		if (s[i] < 48 || s[i] > 57)
			b = 0;
	}

	return b;
}

int is_valid_mask(guint32 mask)
{

	/* how do we test this ? */

	/* 255.255.192.0 = 11111111.11111111.11000000.00000000 */
	/* 255.224.0.0 = 11111111.11100000.00000000.00000000 */
	//printf("mask value: %u\n", mask);

	/* a subnet mask must have consective 1-bits starting high-order to lower-order */
	/* can iterate through all possible subnet masks (32 of them) and verify that it is valid */

	int i;
	int is_valid = 0;
	for (i = 8; i <= 32; i++)
	{

		if (mask == get_mask_from_bits(i))
		{
			is_valid = 1;
			continue;
		}
	}

	//printf("is valid?: %i\n", is_valid);

	return is_valid;
}
