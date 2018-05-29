#include "types.h"
#include "string.hpp"
#include "msgpack.hpp"

extern "C" {

    
    #define REG_USER_METHOD "reguser"
    #define USER_LOGIN_METHOD "userlogin"
    #define PARAM_MAX_LEN 256
    #define USER_NAME_MAX_LEN 128
    #define REG_INFO_MAX_LEN 256
    #define END_STR "it's end"
   
    struct transform_param {
		char to[USER_NAME_MAX_LEN];
		uint32_t amount;
    };

    void prints(char *str, uint32_t len);
    void printi(uint32_t u);
    int  call_trx(char *contract , char *method , char *buf , uint32_t buf_len );
    void parse_param(unsigned char *param , int len);
    uint32_t get_param(unsigned char *param, uint32_t buf_len);

    int pack_transform_param(char *buf, uint32_t buflen, transform_param *trans)
    {
		MsgPackCtx ctx;
		msgpack_init(&ctx, buf, buflen);

		pack_array16(&ctx, 2);
		pack_str16(&ctx, trans->to, strlen(trans->to));
		pack_u32(&ctx, trans->amount);

		return ctx.pos;
    }

    
    static inline void myprints(char *s)
    {
        prints(s, strlen(s));
    }

    static inline void printbuf(char *buf  , int len)
    {
	for (int i = 0 ; i < len ; i++)
	{
	     printi(buf[i]);
	}
    }

    
    void init()
    {
    }

    int start(int m ) 
    {
        prints(USER_LOGIN_METHOD, strlen(USER_LOGIN_METHOD));
		unsigned char parameter[256];
		uint32_t len = get_param(parameter , 256);        

		parse_param(parameter , len);


		char *contract = "sub";
		char *method   = "func1";

		//set parameter for sub-trx
		struct transform_param t1;
		strcpy(t1.to , "Tom");
		t1.amount = 1123;

		char bf1[256];

		//package parameter for sub-trx1
		len = pack_transform_param(bf1 , 256 , &t1);

		int res = call_trx(contract , method , bf1 , len);
		if (res != 0){
			return res;
		}
	
		//package parameter for sub-trx2
		struct transform_param t2;
		strcpy(t2.to , "Jack");
		t2.amount = 2433;
	
		len = pack_transform_param(bf1 , 256 , &t2);
	

		char *method2 = "func2";
		res = call_trx(contract , method2 , bf1 , len);
		if (res != 0)
		{
			return res;
		}

		
		prints(END_STR , strlen(END_STR));
		return 0;
    }
}

