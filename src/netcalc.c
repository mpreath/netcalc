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

int main(int argc, char* argv[]) {

	fprintf(stderr, "netcalc  Copyright (C) 2009  Matthew Reath\n");
	fprintf(stderr, "This program comes with ABSOLUTELY NO WARRANTY;\n");
	fprintf(stderr, "This is free software, and you are welcome to redistribute it\n");
	fprintf(stderr, "under certain circumstances. See the included LICENSE file\n");
	fprintf(stderr, "for more information.\n");

	if(argc != 4) {
		fprintf(stderr, "Wrong number of arguments.\n");
		return 0;	
	}

	host h1;

	h1.ip_address = ddtoint(argv[1]);
	h1.mask = ddtoint(argv[2]);
	/*
	network n1;	

	initialize_network(&n1, &h1);
	
	print_network_info(&n1);
	*/

	/*
	unsigned int new_mask;
	int obits,nbits;

	new_mask = extend_mask(h1.mask, 1);
	obits = get_bits_in_mask(h1.mask);
	nbits = get_bits_in_mask(new_mask);
	printf("Old mask:\t/%u\tNew mask:\t/%u\n", obits, nbits);
	*/
	
	tnode *root;
	tnode *t1;

	root = (tnode *)malloc(sizeof(tnode));
	initialize_network(&root->n, &h1);

	root->left = NULL;
	root->right = NULL;

	build_tree_host_count(root, atoi(argv[3]));
	/*
	int i = 0;
	int j = atoi(argv[3]);

	t1 = &root;
	for(i = 0; i < j; i++, t1=t1->left) {
		split_network(t1);
	}

	*/
	/* print out tree */
	print_network_tree(root);

	/* release tree memory */	
	free_network_tree(root);	
	
	return 0;
}
