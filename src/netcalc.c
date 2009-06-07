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
#include <config.h>
#include <utility.h>
#include <stdio.h>
#include <stdlib.h>
#include <host.h>
#include <network.h>
#include <network_tree.h>
#include <string.h>
#include <glib.h>

#define MAX_NUM_TREES 256

void print_info();
void print_usage();
void net_info(char* ip_address, char* mask);
void host_tree(char* ip_address, char* mask, int hosts);
void net_tree(char* ip_address, char* mask, int nets);
void vlsm_tree(char* ip_address, char* mask, char* nets);
void net_summary();

int main(int argc, char* argv[]) {

	if(argc == 3 && argv[1][0] != '-') {
		print_info();
		net_info(argv[1], argv[2]);
	} else if(argc > 1 && argv[1][0] == '-') {
		switch(argv[1][1]) {
		
			case 'h':
				print_info();
				host_tree(argv[3], argv[4], atoi(argv[2]));
				break;
			case 'n':
				print_info();
				net_tree(argv[3], argv[4], atoi(argv[2]));
				break;
			case 'v':
				print_info();
				vlsm_tree(argv[3], argv[4], argv[2]);
				break;
			case 's':
				net_summary();
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
/*
	fprintf(stderr, "netcalc %s  Copyright (C) 2009  Matthew Reath\n", VERSION);
        fprintf(stderr, "This program comes with ABSOLUTELY NO WARRANTY;\n");
        fprintf(stderr, "This is free software, and you are welcome to redistribute it\n");
        fprintf(stderr, "under certain circumstances. See the included LICENSE file\n");
        fprintf(stderr, "for more information.\n");
*/
}

void print_usage() {
	fprintf(stderr, "netcalc %s usage:\n\n", VERSION);
	fprintf(stderr, "netcalc <address> <mask>\n");
	fprintf(stderr, "netcalc -h <host count> <address> <mask>\n");
	fprintf(stderr, "netcalc -n <network count> <address> <mask>\n");
	fprintf(stderr, "netcalc -v <host list> <address> <mask>\n\n");

	fprintf(stderr, "Examples:\n\n");
	fprintf(stderr, "netcalc 192.168.2.10 255.255.255.252\n");
	fprintf(stderr, "netcalc -h 50 192.168.2.0 255.255.255.0\n");
	fprintf(stderr, "netcalc -n 8 192.168.2.0 255.255.255.0\n");
	fprintf(stderr, "netcalc -v 2,2,2,2,50,50 192.168.2.0 255.255.255.0\n");
	

}

void net_info(char* ip_address, char* mask) {

	host h1;

	initialize_host(&h1, ip_address, mask);

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
	

	tnode* networks[MAX_NUM_TREES]; 
	host h1,h2;
	char ip[16];
	char mask[16];
	int i;

	/* initialize the nodes to NULL */
	for(i = 0; i < MAX_NUM_TREES; i++)
		networks[i] = NULL;

	/* populate the array with the networks from the input */
	for(i = 0; scanf("%s %s", ip, mask) != EOF; i++) {
		initialize_host(&h1, ip, mask);
		networks[i] = malloc(sizeof(tnode));
		initialize_network(&networks[i]->n, &h1);		
		networks[i]->left = NULL;
		networks[i]->right = NULL;
		networks[i]->parent = NULL;
	}


	int j;
	int k;
	int l;
	tnode *t1;
	int made_changes = 1;

	/* loop until there are no summarization left to do */
	while(made_changes) {

		made_changes = 0;
	
		/* loop through each member */
		for(j = 0; j < i; j++) {

			/* compare each member for summarization*/
			for(k = 0; k < i; k++) {

				if((t1 = combine_networks(networks[j], networks[k])) != NULL) {
					networks[j] = NULL;
					networks[k] = NULL;

					for(l = 0; networks[l] != NULL; l++)
						;

					networks[l] = t1;
				
				
					made_changes = 1;
				
				}

			}
		}

	}

	for(i = 0; i < MAX_NUM_TREES; i++) {
		if(networks[i] != NULL) {
			print_network_tree(networks[i]);
			free_network_tree(networks[i]);
		}
	}

}


