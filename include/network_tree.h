#ifndef __NETWORK_TREE_H__
#define __NETWORK_TREE_H__

#include <utility.h>
#include <host.h>
#include <network.h>
#include <stdlib.h>

/*
typedef struct network_tree_node *tptr;

typedef struct network_tree_node {

	network n;
	tptr right;
	tptr left;

} tnode;

int split_network(tptr n1);
*/

typedef struct network_tree_node {
	
	network n;
	struct network_tree_node *left;
	struct network_tree_node *right;

} tnode;

int split_network(tnode *n1);

void print_network_tree(tnode *n1);

void build_tree_net_count(tnode *t1, int nets);

void split_to_depth(tnode *t1, int depth, int target);

void build_tree_host_count(tnode *t1, int hosts);
#endif
