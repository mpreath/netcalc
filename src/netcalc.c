/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    netcalc.c - main entry point into the netcalc application

    Copyright (C) 2011  Matthew Reath 

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

static char* vlsm_counts = NULL;
static char* ip_address = NULL;
static char* netmask = NULL;
static gint do_host_count = 0;
static gint do_net_count = 0;
static gboolean do_summary = FALSE;
static gboolean verbose = FALSE;
//static gboolean do_help = FALSE;
static GOptionEntry entries[]  =
{
	{ "verbose", 'v', 0, G_OPTION_ARG_NONE, &verbose,
	"Use verbose output", NULL },
	{ "hosts", 'h', 0, G_OPTION_ARG_INT, &do_host_count,
	"Specify the host count to use to subnet the network [IPv4 only]",
	"HOSTS"},
	{ "nets", 'n', 0, G_OPTION_ARG_INT, &do_net_count,
	"Specify the net count to use to subnet the network",
	"NETS"},
	{ "vlsm", 'l', 0, G_OPTION_ARG_STRING, &vlsm_counts,
	"Comma seperated list of host count requirements for VLSM network",
	"VLSM_LIST"},
	{ "summary", 's', 0, G_OPTION_ARG_NONE, &do_summary,
	"Summarize a list of subnets into a one or more supernets",
	NULL},
	{ "ip-address", 'i', 0, G_OPTION_ARG_STRING, &ip_address,
	"Specific the IPv4 Address to use for calculation", 
	NULL },
	{ "netmask", 'm', 0, G_OPTION_ARG_STRING, &netmask,
	"Specific the IPv4 Network Mask to use for calculation", 
	NULL },
	{ NULL }
};



int main(int argc, char* argv[]) {

	GError *error = NULL;
	GOptionContext *context;
	context = g_option_context_new ("- calculate network information");
	g_option_context_add_main_entries (context, entries, NULL);
	//g_option_context_add_group (context, glib_get_option_group (TRUE));
	if(!g_option_context_parse(context, &argc, &argv, &error))
	{
		g_print("option parsing failed: %s\n", error->message);
		exit(1);
	}	
	

	if(verbose) {
		print_info();
	}

	if(vlsm_counts) {
		if(ip_address && netmask)
			vlsm_tree(ip_address,netmask,vlsm_counts);
		else
			g_print("no ip address or netmask specified, use -? for usage information\n");
	}
	else if(do_host_count) {
		if(ip_address && netmask)
			host_tree(ip_address,netmask,do_host_count);
		else
			g_print("no ip address or netmask specified, use -? for usage information\n");
	}
	else if(do_net_count) {
		if(ip_address && netmask)
			net_tree(ip_address,netmask,do_net_count);
		else
			g_print("no ip address or netmask specified, use -? for usage information\n");
	}
	else if(do_summary) {
		net_summary();
	} 
	else 
	{

		if(ip_address && netmask)
			net_info(ip_address,netmask);
		else
			g_print("no ip address or netmask specified, use -? for usage information\n");
			
	}

	return 0;
}

void print_info() {
	fprintf(stderr, "netcalc %s  Copyright (C) 2017 Matt Reath\n", VERSION);
        fprintf(stderr, "This program comes with ABSOLUTELY NO WARRANTY;\n");
        fprintf(stderr, "This is free software, and you are welcome to redistribute it\n");
        fprintf(stderr, "under certain circumstances. See the included COPYING file\n");
        fprintf(stderr, "for more information.\n\n");
}

void net_info(char* ip_address, char* mask) {

	host h1;
	network n1;

	initialize_host(&h1, ip_address, mask);
	initialize_network(&n1, &h1);
	print_network_info(&n1);

}

void host_tree(char* ip_address, char* mask, int hosts) {

	host h1;

	initialize_host(&h1, ip_address, mask);

	tnode* t1;

	t1 = (tnode *) malloc(sizeof(tnode));

	t1->right = NULL;
	t1->left = NULL;
	t1->parent = NULL;

	initialize_network(&t1->n, &h1);

	build_tree_host_count(t1, hosts);

	print_network_tree(t1, 0, verbose);

	free_network_tree(t1);
	
}

void net_tree(char* ip_address, char* mask, int nets) {

	host h1;

	initialize_host(&h1, ip_address, mask);

	tnode* t1;

	t1 = (tnode *) malloc(sizeof(tnode));

	t1->right = NULL;
	t1->left = NULL;
	t1->parent = NULL;

	initialize_network(&t1->n, &h1);

	build_tree_net_count(t1, nets);

	print_network_tree(t1, 0, verbose);

	free_network_tree(t1);


}

void vlsm_tree(char* ip_address, char* mask, char* nets) {

	host h1;


	initialize_host(&h1, ip_address, mask);

	tnode *t1;

	t1 = (tnode *) malloc(sizeof(tnode));

	t1->right = NULL;
	t1->left = NULL;
	t1->parent = NULL;

	initialize_network(&t1->n, &h1);

	//int j = 0;
	char* tok;
	char* sep = ",";
	
	
	for(tok = strtok(nets, sep); tok; tok = strtok(NULL,sep)) {
		//verify host counts are legit
		if(!is_number(tok))
			g_error("vlsm string contains non-numeric values");
	
		// we should have left and right as an option here, set by CLI flag
		// default is left
		build_tree_vlsm(t1, atoi(tok), 0);
	}

	print_network_tree(t1, 0, verbose);

	free_network_tree(t1);
}

void net_summary() {
	

	tnode* networks[MAX_NUM_TREES]; 
	host h1;
	char ip[16];
	char mask[16];
	int i;

	/* initialize the nodes to NULL */
	for(i = 0; i < MAX_NUM_TREES; i++)
		networks[i] = NULL;

	/* populate the array with the networks from the input */
	for(i = 0; scanf("%s %s", ip, mask) != EOF; i++) {
		initialize_host(&h1, ip, mask);
		networks[i] = g_malloc(sizeof(tnode));
		initialize_network(&networks[i]->n, &h1);		
		networks[i]->left = NULL;
		networks[i]->right = NULL;
		networks[i]->parent = NULL;
	//	printf("Added network %s %s\n", ip, mask);
	}

	//printf("%i\n", i);

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

				if((j != k) && (t1 = combine_networks(networks[j], networks[k])) != NULL) {
	//				printf("Combined networks\n[%i:%i]\n", j,k);
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

	printf("Networks were summarized as follows:\n\n");
	for(i = 0; i < MAX_NUM_TREES; i++) {
		if(networks[i] != NULL) {
			if(verbose) 
			{
				print_network_tree(networks[i], 0, TRUE);
				printf("\n");
			}
			else
			{
				char ip_address[16];
				inttodd(ip_address, networks[i]->n.address.ip_address);
				printf("%s/%i\n", ip_address, get_bits_in_mask(networks[i]->n.address.mask));
			}
			free_network_tree(networks[i]);
		}
	}

}


