/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 *
 * const.h - Constants and error codes.
 */
#ifndef CONST_H
#define CONST_H

#include "stdio.h"

#define SHA512_BITS 512
#define HASH_BITS SHA512_BITS
#define HASH_BYTES HASH_BITS / 8
#define HMAC_KEY_BYTES 16

// Error codes
#define OK 0
#define ERR -1
#define BAD -2

#define ERR_HMAC -3
#define ERR_SHA512 -4
#define ERR_TINY_BAD_PARAMS -5

#define ERR_GRAPH -6
#define ERR_GRAPH_EXCEEDED_MAX_OUT_DEGREE -7
#define ERR_GRAPH_CYCLE -8
#define ERR_GRAPH_MULTI_EDGE -9

#define ERR_BLOOM_PARAMS_MISMATCH -10

#define ERR_DICT_BAD_KEY -11
#define ERR_DICT_BAD_PADDING -12
#define ERR_DICT_TOO_MANY_ITEMS -13
#define ERR_DICT_LONG_VALUE -14

// For tests
//
// TODO Refactor these.
#define ASSERT_ERROR(msg, exp) if (!(exp)) {\
  fprintf(stderr, "error: %s: %s\n", test, msg); ret = ERR;\
}

#define ASSERT_WARNING(msg, exp) if (!(exp)) {\
  fprintf(stderr, "warning: %s: %s\n", test, msg);\
}

#define ASSERT_FATAL(msg, exp) if (!(exp)) {\
  fprintf(stderr, "fatal: %s: %s\n", test, msg); return ERR;\
}

#define ASSERT_EQ_ERROR(var, got, exp) if ((got) != (exp)) {\
  fprintf(stderr, "error: %s: %s = %d, expected %d\n", test, var, got, exp); ret = ERR;\
}

#define ASSERT_EQ_FATAL(var, got, exp) if ((got) != (exp)) {\
  fprintf(stderr, "error: %s: %s = %d, expected %d\n", test, var, got, exp); return ERR;\
}

#define ASSERT_OK_ERROR(var, got) if ((got) != OK)  {\
  fprintf(stderr, "error: %s: %s = %d, expected OK\n", test, var, got); ret = ERR;\
}

#define ASSERT_OK_FATAL(var, got) if ((got) != OK)  {\
  fprintf(stderr, "error: %s: %s = %d, expected OK\n", test, var, got); return ERR;\
}

#endif //CONST_H
