/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    netcalc.c - main entry point into the netcalc application

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

#include <utility.h>
#include <stdio.h>
#include <stdlib.h>
#include <host.h>
#include <network.h>
#include <network_tree.h>

void print_info();
void print_usage();
void net_info(char* ip_address, char* mask);
void host_tree(char* ip_address, char* mask, int hosts);
void net_tree(char* ip_address, char* mask, int nets);
void vlsm_tree(char* ip_address, char* mask, char* nets);
void net_summary();

int main(int argc, char* argv[]) {

	net_info(argv[1], argv[2]);	
	
	return 0;
}

void print_info() {

	fprintf(stderr, "netcalc  Copyright (C) 2009  Matthew Reath\n");
        fprintf(stderr, "This program comes with ABSOLUTELY NO WARRANTY;\n");
        fprintf(stderr, "This is free software, and you are welcome to redistribute it\n");
        fprintf(stderr, "under certain circumstances. See the included LICENSE file\n");
        fprintf(stderr, "for more information.\n");
}

void print_usage() {
	

}

void net_info(char* ip_address, char* mask) {

	host h1;

	h1.ip_address = ddtoint(ip_address);
	h1.mask = ddtoint(mask);

	network n1;

	initialize_network(&n1, &h1);

	print_network_info(&n1);

}

void host_tree(char* ip_address, char* mask, int hosts) {

	
}

void net_tree(char* ip_address, char* mask, int nets) {

}

void vlsm_tree(char* ip_address, char* mask, char* nets) {

}

void net_summary() {

}

