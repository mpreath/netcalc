/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    network_tree.c - subnetting/supernetting related fuctions. It is better
    to handle these functions as a tree structure

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

#include <network_tree.h>
#include <utility.h>
#include <network.h>
#include <host.h>
#include <stdio.h>
#include <math.h>

int split_network(tnode *n1) {

	/* verify the network hasn't already been split */
	if(n1->left != NULL)
		return -1;

	host h1, h2;
	tnode *t1, *t2;

	/* first network of split */
	h1.ip_address = n1->n.address.ip_address;
	h1.mask = extend_mask(n1->n.address.mask, 1);

	/* second network of split */
	h2.ip_address = get_broadcast_address(&h1) + 1;
	h2.mask = extend_mask(n1->n.address.mask, 1);

	/* create first network */
	n1->left = (tnode *)malloc(sizeof(tnode));
	t1 = n1->left;
	initialize_network(&t1->n, &h1);
	t1->left = NULL;
	t1->right = NULL;
	t1->parent = n1;

	/* create second network */
	n1->right = (tnode *)malloc(sizeof(tnode));
	t2 = n1->right;
	initialize_network(&t2->n, &h2);
	t2->left = NULL;
	t2->right = NULL;
	t2->parent = n1;
	

	return 1;
}

void print_network_tree(tnode *n1) {

	char ip_address[16];

	//if(n1->left == NULL && n1->right == NULL) {

		inttodd(ip_address, n1->n.address.ip_address);

		printf("Network: %s/%u", ip_address, get_bits_in_mask(n1->n.address.mask));
		if(n1->in_use)
			printf("*");

		printf("\n");
	//}

	if(n1->left != NULL)
		print_network_tree(n1->left);
	if(n1->right != NULL)
		print_network_tree(n1->right);

}

void build_tree_net_count(tnode *t1, int nets) {

	int tdepth;
	double val;

	// caclulate depth to split to
	// log base 2 on number of nets to
	// get depth of tree
	val = log(nets) / log(2);
	val = ceil(val);

	tdepth = (int)val;
	split_to_depth(t1,0,tdepth);
}

void build_tree_host_count(tnode *t1, int hosts) {

	double val;
	int tdepth;
	int bdepth;

	val = log(hosts) / log(2);
	val = ceil(val);
	
	bdepth = get_bits_in_mask(t1->n.address.mask);
	
	tdepth = 32 - val;

	split_to_depth(t1, bdepth, tdepth);  
}

void build_tree_vlsm(tnode *t1, int hosts, int right) {

	if(t1->n.host_count/2 >= hosts) {
		if(t1->left == NULL && t1->right == NULL)
			split_network(t1);

		if(right) {
			build_tree_vlsm(t1->right, hosts, right);
		} else {
			build_tree_vlsm(t1->left, hosts, right);
		}
	} else {
		/* ok, we are at the end point for our host count */

		/* make sure we are a leaf , and not in use*/

		if(t1->left == NULL && t1->right == NULL && t1->in_use == 0) {
			t1->in_use = 1;
			return;
		} else {
			/* we need to traverse up the tree */
			if(right) {
		
				if(t1 == t1->parent->left)
					build_tree_vlsm(t1->parent->parent->left, hosts, right);
				else if(t1 != t1->parent->left)
					build_tree_vlsm(t1->parent->left, hosts, right);

			} else {

				if(t1 == t1->parent->right)
					build_tree_vlsm(t1->parent->parent->right, hosts, right);
				else if(t1 != t1->parent->right)
					build_tree_vlsm(t1->parent->right, hosts, right);
			}
		}

	} 

	
		

}

void split_to_depth(tnode *t1, int depth, int target) {

	if(depth >= target)
		return;

	split_network(t1);
	
	split_to_depth(t1->left, depth+1, target);
	split_to_depth(t1->right, depth+1, target);

}


tnode* combine_networks(tnode *s1, tnode *s2) {

	unsigned int new_net1, new_net2;;
	unsigned int new_mask;
	host h1, h2;

	tnode *t1;

	if(s1->n.address.mask == s2->n.address.mask) {

		//printf("Masks are the same: %u\n", get_bits_in_mask(s1->n.address.mask));
		
		new_mask = shorten_mask(s1->n.address.mask, 1);


	
		//printf("New mask: %u, %u[%u]\n", s1->n.address.mask, new_mask, get_bits_in_mask(new_mask));	

		h1.ip_address = s1->n.address.ip_address;
		h1.mask = new_mask;

		h2.ip_address = s2->n.address.ip_address;
		h2.mask = new_mask;

		new_net1 = get_network_address(&h1);
		new_net2 = get_network_address(&h2);
		/*
		char ip_address[16];
		inttodd(ip_address, new_net1);
		printf("New net 1: %s\n", ip_address);
		inttodd(ip_address, new_net2);
		printf("New net 2: %s\n", ip_address);
		*/
		if(new_net1 == new_net2) {
			// we can summarize these two
			//printf("Summarizing.\n");

			t1 = (tnode *)malloc(sizeof(tnode));

			initialize_network(&t1->n, &h1);
			/* set pointers */
			s1->parent = t1;
			s2->parent = t1;
			if(s1->n.address.ip_address < s2->n.address.ip_address) {
				t1->left = s1;
				t1->right = s2;
			} else {
				t1->right = s1;
				t1->left = s2;
			}

			return t1;

		} else {

			return NULL;
		}
	} else {
		return NULL;
	}

}

void free_network_tree(tnode* t1) {

	if(t1->left != NULL && t1->right != NULL) {
		free_network_tree(t1->left);
		free_network_tree(t1->right);
	}

	free(t1);
	
}
