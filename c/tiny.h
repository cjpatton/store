/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 *
 * tiny.h -- Hashing over a small range.
 *
 * Implementation of a hash and pseudorandom function (from SHA512 and
 * HMAC-SHA512 resp.) suitable for construction Bloom filters and the like.
 */
#ifndef TINY_H
#define TINY_H

#include "openssl/hmac.h"
#include "openssl/sha.h"

#include "const.h"

// Context for tinyhash and tinyprf.
//
// WARNING Not thread safe!
typedef struct Tiny {

  // Parameters of the algorithm.
  int radix, chunk_bits, chunks;

  // Context for SHA512 (hash) and HMAC-SHA512 (prf).
  union {
    HMAC_CTX hmac;
    SHA512_CTX sha;
  } _st;

  // Indicates whether to use (tiny)prf or (tiny)hash.
  int _use_prf;

  // Stores the output of SHA512 (hash) or HMAC-SHA512 (prf).
  char digest [HASH_BYTES];

} tiny_ctx;

// Creates a tiny_ctx* and returns it.
//
// Computes parameters from `radix. Initializes the HMAC context. Must be freed
// via tinyprf_ctx_free(). Returns NULL if ceil(log2(radix)) not in range [1,
// HASH_BITS].
tiny_ctx *tinyprf_new(int radix);

// Frees memory associated to 'ctx'.
void tinyprf_free(tiny_ctx *ctx);

// Initializes tiny_ctx* ctx with a secret key 'key' of length HMAC_KEY_BYTES.
//
// Returns OK if success and ERR_HMAC otherwise.
int tinyprf_init(tiny_ctx *ctx, const char *key);

// Initializes tiny_ctx* ctx with a randomly-generated key of length HMAC_KEY_BYTES.
int tinyprf_init_generate_key(tiny_ctx *ctx);

// Computes HMAC-SHA512 of 'tweak | in' and outputs the result to 'out'.
//
// TODO Perhaps assert that out_bytes == HASH_BYTES and remove this parameter.
// Here and for hash().
int prf(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes, char *out, int out_bytes);

// Maps inputs to an integer in range [0, ctx->radix).
//
// The string 'tweak' is prepended to 'in'. Returns BAD if no suitable output is
// available and ERR if something went wrong with HMAC.
int tinyprf(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes);

// Creates a tiny_ctx* and returns it.
//
// Must be freed via tiny_ctx_free(). Returns NULL if ceil(log2(radix)) not in
// range [1, HASH_BITS].
tiny_ctx *tinyhash_new(int radix);

// Frees memory associated to 'ctx'.
void tinyhash_free(tiny_ctx *ctx);

// Initializes tiny_ctx* ctx.
//
// Returns OK if success and ERR otherwise.
int tinyhash_init(tiny_ctx *ctx);

// Computes SHA512 of 'tweak | in' and outputs the result to 'out'..
int hash(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes, char *out, int out_bytes);

// Maps inputs to an integer in range [0, ctx->radix).
//
// The string 'tweak' is prepended to 'in'.  Returns BAD if no suitable value
// was found.
int tinyhash(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes);

#endif //TINY_H
