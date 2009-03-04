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
#include <string.h>

void print_info();
void print_usage();
void net_info(char* ip_address, char* mask);
void host_tree(char* ip_address, char* mask, int hosts);
void net_tree(char* ip_address, char* mask, int nets);
void vlsm_tree(char* ip_address, char* mask, char* nets);
void net_summary();

int main(int argc, char* argv[]) {

	if(argc == 3 && argv[1][0] != '-') {
		net_info(argv[1], argv[2]);
	} else if(argc > 1 && argv[1][0] == '-') {
		switch(argv[1][1]) {
		
			case 'h':
				host_tree(argv[3], argv[4], atoi(argv[2]));
				break;
			case 'n':
				net_tree(argv[3], argv[4], atoi(argv[2]));
				break;
			case 'v':
				vlsm_tree(argv[3], argv[4], argv[2]);
				break;
			case 's':
				printf("Summarization\n");
				break;
			default:
				print_usage();
		}	
	} else {
		print_usage();
	}
	
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

	host h1;

	initialize_host(&h1, ip_address, mask);

	tnode* t1;

	t1 = (tnode *) malloc(sizeof(tnode));

	initialize_network(&t1->n, &h1);

	build_tree_host_count(t1, hosts);

	print_network_tree(t1);

	free_network_tree(t1);
	
}

void net_tree(char* ip_address, char* mask, int nets) {

	host h1;

	initialize_host(&h1, ip_address, mask);

	tnode* t1;

	t1 = (tnode *) malloc(sizeof(tnode));

	initialize_network(&t1->n, &h1);

	build_tree_net_count(t1, nets);

	print_network_tree(t1);

	free_network_tree(t1);


}

void vlsm_tree(char* ip_address, char* mask, char* nets) {

	host h1;

	initialize_host(&h1, ip_address, mask);

	tnode *t1;

	t1 = (tnode *) malloc(sizeof(tnode));

	initialize_network(&t1->n, &h1);

	char* tok;
	char* sep = ",";

	for(tok = strtok(nets, sep); tok; tok = strtok(NULL,sep)) {
		build_tree_vlsm(t1, atoi(tok), 0);
	}

	print_network_tree(t1);

	free_network_tree(t1);
}

void net_summary() {

}

