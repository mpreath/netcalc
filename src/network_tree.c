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
	
	/* create second network */
	n1->right = (tnode *)malloc(sizeof(tnode));
	t2 = n1->right;
	initialize_network(&t2->n, &h2);
	t2->left = NULL;
	t2->right = NULL;

	return 1;
}

void print_network_tree(tnode *n1) {

	char ip_address[16];

	if(n1->left == NULL && n1->right == NULL) {

		inttodd(ip_address, n1->n.address.ip_address);

		printf("Network: %s/%u\n", ip_address, get_bits_in_mask(n1->n.address.mask));
	}

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

build_tree_vlsm(tnode *t1, int hosts) {

	/* can the current network support the host requirements? */
	

	/* if not, and children are null split network */

	/* if not, and children are not null, traverse to children */

	

}

void split_to_depth(tnode *t1, int depth, int target) {

	if(depth >= target)
		return;

	split_network(t1);
	
	split_to_depth(t1->left, depth+1, target);
	split_to_depth(t1->right, depth+1, target);

}

void free_network_tree(tnode* t1) {

	if(t1->left != NULL && t1->right != NULL) {
		free_network_tree(t1->left);
		free_network_tree(t1->right);
	}

	free(t1);
	
}
