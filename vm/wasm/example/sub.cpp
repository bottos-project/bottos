#include "types.h"
#include "string.hpp"

extern "C" {

    
   #define MSG "I 'm sub-contract !!!"
   #define DEFAULT "default method"
   #define REGUSER "reguser method"

    void prints(char *str, uint32_t len);
    void parse_param(unsigned char *param , int len);
    uint32_t get_param(unsigned char *param, uint32_t buf_len);    
    
    static inline void myprints(char *s)
    {
        prints(s, strlen(s));
    }

    
    void init()
    {
    }
	

    int start(int method)
    //int start(int method , char *param , int len) 
    {
        prints(MSG, strlen(MSG));
	unsigned char parameter[256];
	uint32_t len = get_param(parameter, 256);
	
	switch (method){
	case 104:
		prints(REGUSER , strlen(REGUSER));
		parse_param(parameter , len);
		break;
	default:
		prints(DEFAULT , strlen(DEFAULT));
	}
	

	return 0;
    }
}
