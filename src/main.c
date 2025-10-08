#include <glib.h>

static gboolean verbose = FALSE;

int 
main (int argc, char *argv[])
{
    static GOptionEntry options[] = {
        { "verbose", 'v', 0, G_OPTION_ARG_NONE, &verbose, "Be verbose", NULL },
        { NULL }
    };

    GError *error = NULL;
    GOptionContext *context;

    context = g_option_context_new ("- autobase");
    g_option_context_add_main_entries (context, options, NULL);
    if (!g_option_context_parse (context, &argc, &argv, &error)) {
        g_print ("option parsing failed: %s\n", error->message);
        exit (1);
    }

    if (verbose) {
        g_print ("Verbose mode\n");
    }

    g_print ("Welcome to Autobase\n");

    return 0;
}
