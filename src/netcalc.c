/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    netcalc.c - main entry point into the netcalc application

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
#include <config.h>
#include <utility.h>
#include <stdio.h>
#include <stdlib.h>
#include <host.h>
#include <network.h>
#include <network_tree.h>
#include <string.h>
#include <glib.h>
#include <glib/gprintf.h>

#define MAX_NUM_TREES 256

void print_info();
void print_usage();
void net_info(char *ip_address, char *mask);
void host_tree(char *ip_address, char *mask, int hosts);
void net_tree(char *ip_address, char *mask, int nets);
void vlsm_tree(char *ip_address, char *mask, char *nets);
void net_summary();

static char *vlsm_counts = NULL;
static char *ip_address = NULL;
static char *netmask = NULL;
static gint do_host_count = 0;
static gint do_net_count = 0;
static gboolean do_summary = FALSE;
static gboolean verbose = FALSE;
static gboolean version = FALSE;
//static gboolean do_help = FALSE;
static GOptionEntry entries[] =
	{
		{"verbose", 'v', 0, G_OPTION_ARG_NONE, &verbose,
		 "Use verbose output", NULL},
		 {"version", 'V', 0, G_OPTION_ARG_NONE, &version,
		 "Display version information", NULL},
		{"hosts", 'o', 0, G_OPTION_ARG_INT, &do_host_count,
		 "Specify the host count to use to subnet the network",
		 "HOSTS"},
		{"nets", 'n', 0, G_OPTION_ARG_INT, &do_net_count,
		 "Specify the net count to use to subnet the network",
		 "NETS"},
		{"vlsm", 'l', 0, G_OPTION_ARG_STRING, &vlsm_counts,
		 "Comma seperated list of host counts for VLSM network",
		 "VLSM_LIST"},
		{"summary", 's', 0, G_OPTION_ARG_NONE, &do_summary,
		 "Summarize a list of subnets into a one or more supernets",
		 NULL},
		{NULL}};

int main(int argc, char *argv[])
{

	GError *error = NULL;
	GOptionContext *context;
	context = g_option_context_new("[ip_address] [mask] - calculate network information");
	g_option_context_add_main_entries(context, entries, NULL);
	//g_option_context_add_group (context, glib_get_option_group (TRUE));
	if (!g_option_context_parse(context, &argc, &argv, &error))
	{
		g_print("option parsing failed: %s\n", error->message);
		exit(1);
	}

	if (argc == 3)
	{
		ip_address = argv[argc - 2];
		netmask = argv[argc - 1];
	} else if (argc == 2)
	{
		gchar** split_values = g_strsplit(argv[argc - 1], "/", 2);
		if(split_values[0] != NULL)
		{
			ip_address = (gchar *) malloc(1 + sizeof(gchar) * strlen(split_values[0]));
			g_stpcpy(ip_address, split_values[0]);
		}
		if(split_values[1] != NULL)
		{
			if(is_number(split_values[1]))
			{
				netmask = (gchar *) malloc (sizeof(gchar) * 16); 
				inttodd(netmask, cidrtoint(split_values[1]));
			}
		}
	}

	
	if (vlsm_counts)
	{
		if (ip_address && netmask)
		{
			if (verbose)
				print_info();
			vlsm_tree(ip_address, netmask, vlsm_counts);
		}
		else
			g_print("no ip address and/or netmask specified, use --help for usage information\n");
	}
	else if (do_host_count)
	{
		if (ip_address && netmask)
		{
			if (verbose)
				print_info();
			host_tree(ip_address, netmask, do_host_count);
		}
		else
			g_print("no ip address and/or netmask specified, use --help for usage information\n");
	}
	else if (do_net_count)
	{
		if (ip_address && netmask)
		{
			if (verbose)
				print_info();
			net_tree(ip_address, netmask, do_net_count);
		}
		else
			g_print("no ip address and/or netmask specified, use --help for usage information\n");
	}
	else if (do_summary)
	{
		if (verbose)
			print_info();
		net_summary();
	}
	else if (version)
	{
		if (verbose)
			print_info();
		else
			printf("%s\n", PACKAGE_STRING);
	}
	else
	{

		if (ip_address && netmask)
		{
			if (verbose)
				print_info();
			net_info(ip_address, netmask);
		}
		else
			g_print("%s", g_option_context_get_help(context, TRUE, NULL));
	}

	return 0;
}

void print_info()
{
	fprintf(stderr, "%s Copyright (c) 2021 Matt Reath\n", PACKAGE_STRING);
	fprintf(stderr, "This program comes with ABSOLUTELY NO WARRANTY;\n");
	fprintf(stderr, "This is free software, and you are welcome to redistribute it\n");
	fprintf(stderr, "under certain circumstances. See the included COPYING file\n");
	fprintf(stderr, "for more information.\n");
}

void net_info(char *ip_address, char *mask)
{

	host h1;
	network n1;

	initialize_host(&h1, ip_address, mask);
	initialize_network(&n1, &h1);
	print_network_info(&n1);
}

void host_tree(char *ip_address, char *mask, int hosts)
{

	host h1;
	tnode *t1;

	initialize_host(&h1, ip_address, mask);

	t1 = (tnode *)malloc(sizeof(tnode));

	t1->right = NULL;
	t1->left = NULL;
	t1->parent = NULL;

	initialize_network(&t1->n, &h1);
	build_tree_host_count(t1, hosts);
	print_network_tree(t1, 0, verbose);
	free_network_tree(t1);
}

void net_tree(char *ip_address, char *mask, int nets)
{

	host h1;
	tnode *t1;

	initialize_host(&h1, ip_address, mask);

	t1 = (tnode *)malloc(sizeof(tnode));

	t1->right = NULL;
	t1->left = NULL;
	t1->parent = NULL;

	initialize_network(&t1->n, &h1);
	build_tree_net_count(t1, nets);
	print_network_tree(t1, 0, verbose);
	free_network_tree(t1);
}

void vlsm_tree(char *ip_address, char *mask, char *nets)
{

	host h1;
	tnode *t1;

	initialize_host(&h1, ip_address, mask);

	t1 = (tnode *)malloc(sizeof(tnode));

	t1->right = NULL;
	t1->left = NULL;
	t1->parent = NULL;

	initialize_network(&t1->n, &h1);

	//int j = 0;
	char *tok;
	char *sep = ",";

	for (tok = strtok(nets, sep); tok; tok = strtok(NULL, sep))
	{
		//verify host counts are legit
		if (!is_number(tok))
			g_error("vlsm string contains non-numeric values");

		// we should have left and right as an option here, set by CLI flag
		// default is left
		build_tree_vlsm(t1, atoi(tok), 0);
	}

	print_network_tree(t1, 0, verbose);
	free_network_tree(t1);
}

void net_summary()
{

	tnode *networks[MAX_NUM_TREES];
	host h1;
	gchar* buffer;
	gchar* cidr_tok = "/";
	gchar* space_tok = " ";
	size_t bufsize = 32;
	int i;
	int j;
	int k;
	int l;
	tnode *t1;
	int made_changes = 1;

	/* initialize the nodes to NULL */
	for (i = 0; i < MAX_NUM_TREES; i++)
		networks[i] = NULL;

	/* 	read each line from stdin (using readline)
		check for "/" denoting a CIDR format
			if true, then split on "/" and initialize host
				need to develop an initialize host variant for CIDR
				need to develop cidrtoint utility function (convert CIDR to binary representation)
				check that mask is and integer and between 8 and 32 inside CIDRtoint
		if no "/" detected, check for space
			if true, split on space and initialize host as is done below */
			

	buffer = (gchar *)malloc(bufsize * sizeof(gchar));
	/* populate the array with the networks from the input */
	for (i = 0; getline(&buffer, &bufsize, stdin) > 0; i++)
	{
		if (buffer != NULL)
		{
			// remove newline characters from input
			int length = strlen(buffer);
			if (buffer[length-1] == '\n')
				buffer[length-1]  = '\0';
			//printf("%s\n", buffer);

			if(g_strrstr(buffer,cidr_tok)) 
			{
				gchar** split_values = g_strsplit(buffer, cidr_tok, 2);
				initialize_cidr_host(&h1, split_values[0], split_values[1]);
				networks[i] = g_malloc(sizeof(tnode));
				initialize_network(&networks[i]->n, &h1);
				networks[i]->left = NULL;
				networks[i]->right = NULL;
				networks[i]->parent = NULL;
				g_strfreev(split_values);
				
			}
			else if(strstr(buffer,space_tok))
			{
				gchar** split_values = g_strsplit(buffer, space_tok, 2);
				initialize_host(&h1, split_values[0], split_values[1]);
				networks[i] = g_malloc(sizeof(tnode));
				initialize_network(&networks[i]->n, &h1);
				networks[i]->left = NULL;
				networks[i]->right = NULL;
				networks[i]->parent = NULL;
				//	printf("Added network %s %s\n", ip, mask);
				g_strfreev(split_values);
			}
		}
	}

	g_free(buffer);

	//printf("%i\n", i);

	/* loop until there are no summarization left to do */
	while (made_changes)
	{

		made_changes = 0;
		/* loop through each member */
		for (j = 0; j < i; j++)
		{
			/* compare each member for summarization*/
			for (k = 0; k < i; k++)
			{

				if ((j != k) && (t1 = combine_networks(networks[j], networks[k])) != NULL)
				{
					//				printf("Combined networks\n[%i:%i]\n", j,k);
					networks[j] = NULL;
					networks[k] = NULL;

					for (l = 0; networks[l] != NULL; l++)
						;

					networks[l] = t1;
					made_changes = 1;
				}
			}
		}
	}

	if(verbose)
		printf("Networks were summarized as follows:\n\n");

	for (i = 0; i < MAX_NUM_TREES; i++)
	{
		if (networks[i] != NULL)
		{
			if (verbose)
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
