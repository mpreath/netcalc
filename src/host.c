#include <host.h>
#include <utility.h>

void initialize_host(host* h1, char* ip_address, char* mask) {

	h1->ip_address = ddtoint(ip_address);
	h1->mask = ddtoint(mask);

}
