/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    network_tree.h - network tree data structure and functions

    Copyright (C) 2021  Matthew Reath

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

#ifndef __NETWORK_TREE_H__
#define __NETWORK_TREE_H__

#include <utility.h>
#include <host.h>
#include <network.h>
#include <stdlib.h>
#include <glib.h>

/*
typedef struct network_tree_node *tptr;

typedef struct network_tree_node {

	network n;
	tptr right;
	tptr left;

} tnode;

int split_network(tptr n1);
*/

typedef struct network_tree_node
{

	network n;
	int in_use;
	struct network_tree_node *left;
	struct network_tree_node *right;
	struct network_tree_node *parent;

} tnode;

int split_network(tnode *n1);

void print_network_tree(tnode *n1, int depth, gboolean verbose);

void build_tree_net_count(tnode *t1, int nets);

void split_to_depth_hc(tnode *t1, int usable_hosts);

void split_to_depth(tnode *t1, int depth, int target);

void build_tree_host_count(tnode *t1, int hosts);

void build_tree_vlsm(tnode *t1, int hosts, int right);

tnode *combine_networks(tnode *s1, tnode *s2);

void free_network_tree(tnode *t1);

#endif
