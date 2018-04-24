/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "openssl/hmac.h"
#include "openssl/rand.h"
#include "openssl/sha.h"
#include "stdio.h"
#include "string.h"

#include "bits.h"
#include "const.h"
#include "tiny.h"

// Computes paraemters from radix
//
// Called by tinyprf_ctx_new() and tinyhash_ctx_new().  Returns ERR if radix < 2 or
// ceil(log2(radix)) > HASH_BITS.
int compute_tiny_params(tiny_ctx *ctx, int radix) {
  ctx->radix = radix;
  int k = ceillog2(radix);
  if (radix < 1 || k > HASH_BITS || k < 1) {
    return ERR_TINY_BAD_PARAMS;
  }
  ctx->chunk_bits = k;
  ctx->chunks = HASH_BITS / k;
  return OK;
}

tiny_ctx *tinyprf_new(int radix) {
  tiny_ctx *ctx = malloc(sizeof(tiny_ctx));
  ctx->_use_prf = 1;

  int err = compute_tiny_params(ctx, radix);
  if (err != OK) {
    free(ctx);
    return NULL;
  }

  // FIXME This constructor is deprecated in OpenSSL>=1.0.1. It was changed to:
  //
  //   HMAC_CTX *ctx = HMAC_CTX_new();
  //
  // My system (Ubuntu 16.04) has 1.0.0 installed.
  HMAC_CTX_init(&(ctx->_st.hmac));

  return ctx;
}

void tinyprf_free(tiny_ctx *ctx) {
  // FIXME New destructor is:
  //
  //   HMAC_CTX_free(ctx);
  //
  // See above.
  //
  // This erases the key:
  // https://www.openssl.org/docs/man1.0.2/crypto/hmac.html
  HMAC_CTX_cleanup(&(ctx->_st.hmac));

  // Clear the saved digest.
  memset(ctx->digest, 0, HASH_BYTES*sizeof(char));

  free(ctx);
}

int tinyprf_init(tiny_ctx *ctx, const char *key) {
  if (ctx == NULL) {
    return ERR;
  }

  // Initialize HMAC-SHA512 context with 'key' and return the result.
  int res = HMAC_Init_ex(&(ctx->_st.hmac), key, HMAC_KEY_BYTES, EVP_sha512(), NULL);
  if (!res) {
    return ERR_HMAC;
  }

  return OK;
}

int tinyprf_init_generate_key(tiny_ctx *ctx) {
  char key[HMAC_KEY_BYTES];
  RAND_bytes((unsigned char *)key, HMAC_KEY_BYTES);
  return tinyprf_init(ctx, key);
}

int prf(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes, char *out, int out_bytes) {

  // Compute HMAC-SHA512 digest of 'tweak || in'.
  int res = HMAC_Update(&(ctx->_st.hmac), (unsigned char *)tweak, tweak_bytes);
  if (!res) {
    return ERR_HMAC;
  }

  res = HMAC_Update(&(ctx->_st.hmac), (unsigned char *)in, in_bytes);
  if (!res) {
    return ERR_HMAC;
  }

  res = HMAC_Final(&(ctx->_st.hmac), (unsigned char *)ctx->digest, NULL);
  if (!res) {
    return ERR_HMAC;
  }

  if (out != NULL) {
    memcpy(out, ctx->digest, out_bytes);
  }

  // Reset the state of the HMAC context. key == NULL and evp_md == NULL, so
  // the same key is used.
  res = HMAC_Init_ex(&(ctx->_st.hmac), NULL, 0, NULL, NULL);
  if (!res) {
    return ERR_HMAC;
  }

  return OK;
}

int tinyprf(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes) {

  int res = prf(ctx, in, in_bytes, tweak, tweak_bytes, NULL, 0);
  if (res != OK) {
    return res;
  }

  int z = BAD;
  for (int i = 0; i < ctx->chunks; i++) {
    int y = get_chunk(ctx->digest, ctx->chunk_bits, i);
    if (y < ctx->radix) {
      z = y;
    }
  }

  return z;
}

tiny_ctx *tinyhash_new(int radix) {
  tiny_ctx *ctx = malloc(sizeof(tiny_ctx));
  ctx->_use_prf = 0;

  int err = compute_tiny_params(ctx, radix);
  if (err != OK) {
    free(ctx);
    return NULL;
  }

  return ctx;
}

void tinyhash_free(tiny_ctx *ctx) {
  // Clear the saved digest.
  memset(ctx->digest, 0, HASH_BYTES*sizeof(char));
  free(ctx);
}

int tinyhash_init(tiny_ctx *ctx){
  int res = SHA512_Init(&(ctx->_st.sha));
  if (!res) {
    return ERR_SHA512;
  }

  return OK;
}

int hash(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes, char *out, int out_bytes) {

  // Compute SHA512 hash of 'tweak || in'.
  int res = SHA512_Init(&(ctx->_st.sha));
  if (!res) {
    return ERR_SHA512;
  }

  res = SHA512_Update(&(ctx->_st.sha), (unsigned char *)tweak, tweak_bytes);
  if (!res) {
    return ERR_SHA512;
  }

  res = SHA512_Update(&(ctx->_st.sha), (const unsigned char *)in, in_bytes);
  if (!res) {
    return ERR_SHA512;
  }

  res = SHA512_Final((unsigned char *)ctx->digest, &(ctx->_st.sha));
  if (!res) {
    return ERR_SHA512;
  }

  if (out != NULL) {
    memcpy(out, ctx->digest, out_bytes);
  }

  return OK;
}

int tinyhash(tiny_ctx *ctx, const char *in, int in_bytes, const char *tweak,
    int tweak_bytes) {

  int res = hash(ctx, in, in_bytes, tweak, tweak_bytes, NULL, 0);
  if (res != OK) {
    return res;
  }

  // Find uniform-random value in [0, ctx->radix).
  int z = BAD;
  for (int i = 0; i < ctx->chunks; i++) {
    int y = get_chunk(ctx->digest, ctx->chunk_bits, i);
    if (y < ctx->radix) {
      z = y;
    }
  }

  return z;
}
