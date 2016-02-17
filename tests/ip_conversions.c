/*
    This program is part of the test suite. It exercises the preprocessor
    macros in the header file to ensure Go produces correct IP address encodings.
*/

#include <lwipv6.h>
#include <stdio.h>

int main(int argc, char** argv) {
    struct ip_addr addr;
    struct ip_addr mask;

    /* Test 1 */
    IP64_ADDR(&addr,0,0,0,0);
    IP64_MASKADDR(&mask,0,0,0,0);
        
    printf("0.0.0.0/0.0.0.0 %u %u %u %u %u %u %u %u\n", 
        addr.addr[0], addr.addr[1], addr.addr[2], addr.addr[3], 
        mask.addr[0], mask.addr[1], mask.addr[2], mask.addr[3]);

    /* Test 2 */
    IP64_ADDR(&addr,192,168,250,20);
    IP64_MASKADDR(&mask,192, 168, 250, 20);    

    printf("192.168.250.20/192.168.250.20 %u %u %u %u %u %u %u %u\n", 
        addr.addr[0], addr.addr[1], addr.addr[2], addr.addr[3], 
        mask.addr[0], mask.addr[1], mask.addr[2], mask.addr[3]);    
    /* Test 3 */
    IP64_ADDR(&addr,255,255,255,255);
    IP64_MASKADDR(&mask,255, 255, 255, 255);
    printf("255.255.255.255/255.255.255.255 %u %u %u %u %u %u %u %u\n", 
        addr.addr[0], addr.addr[1], addr.addr[2], addr.addr[3], 
        mask.addr[0], mask.addr[1], mask.addr[2], mask.addr[3]);        
        
    /* TODO: IPV6 tests */
    
    return 0;
}
