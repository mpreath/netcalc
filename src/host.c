#include <host.h>
#include <utility.h>
#include <glib.h>

void initialize_host(host* h1, char* ip_address, char* mask) {

	h1->ip_address = ddtoint(ip_address);
	h1->mask = ddtoint(mask);
	if(!is_valid_mask(h1->mask)) {
		g_error("invalid network mask");
	}

}
