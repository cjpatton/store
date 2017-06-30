/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "openssl/rand.h"
#include "stdio.h"
#include "string.h"

#include "bloom.h"
#include "const.h"
#include "tiny.h"

bloom_t *bloom_new(int filter_bits, int salt_bytes, int hash_ct) {
  bloom_t *bloom = malloc(sizeof(bloom_t));
  bloom->filter_bits = filter_bits;
  bloom->filter_bytes = filter_bits / 8;
  if ((filter_bits % 8) > 0) {
    bloom->filter_bytes ++;
  }
  bloom->filter = malloc(bloom->filter_bytes * sizeof(char));
  bloom->salt_bytes = salt_bytes;
  bloom->salt = malloc((salt_bytes+1)* sizeof(char));
  bloom->hash_ct = hash_ct;
  return bloom;
}

void bloom_free(bloom_t *bloom) {
  free(bloom->filter);
  free(bloom->salt);
  free(bloom);
}

int bloom_init(bloom_t *bloom, const char *salt) {
  if (salt) {
    memcpy(bloom->salt, salt, bloom->salt_bytes * sizeof(char));
  }
  memset(bloom->filter, 0, bloom->filter_bytes * sizeof(char));
  return OK;
}

int bloom_init_generate_salt(bloom_t *bloom) {
  RAND_bytes((unsigned char *)bloom->salt, bloom->salt_bytes);
  return bloom_init(bloom, NULL);
}

int compute_hash(tiny_ctx *tiny, const char *in, int in_bytes, char *salt,
    int salt_bytes, int *h1, int *h2) {

  if (tiny->_use_prf) {
    salt[salt_bytes] = 1;
    int res = tinyprf(tiny, in, in_bytes, salt, salt_bytes+1);
    if (res < 0) {
      return res;
    }
    *h1 = res;

    salt[salt_bytes] = 2;
    res = tinyprf(tiny, in, in_bytes, salt, salt_bytes+1);
    if (res < 0) {
      return res;
    }
    *h2 = res;
  } else {
    salt[salt_bytes] = 1;
    int res = tinyhash(tiny, in, in_bytes, salt, salt_bytes+1);
    if (res < 0) {
      return res;
    }
    *h1 = res;

    salt[salt_bytes] = 2;
    res = tinyhash(tiny, in, in_bytes, salt, salt_bytes+1);
    if (res < 0) {
      return res;
    }
    *h2 = res;
  }
  return OK;
}

int bloom_insert(
    bloom_t *bloom,
    tiny_ctx *tiny,
    const char *in, int in_bytes) {

  // TODO Write test
  if (bloom->filter_bits != tiny->radix) {
    return ERR_BLOOM_PARAMS_MISMATCH;
  }

  int h1, h2;
  int err = compute_hash(
      tiny, in, in_bytes, bloom->salt, bloom->salt_bytes, &h1, &h2);
  if (err != OK) {
    return err;
  }

  for (int j = 0; j < bloom->hash_ct; j++) {
    int i = (h1 + (j * h2)) % bloom->filter_bits;
    set_bit(bloom->filter, i);
  }

  return OK;
}

int bloom_get(
    bloom_t *bloom,
    tiny_ctx *tiny,
    const char *in, int in_bytes) {

  // TODO Write test
  if (bloom->filter_bits != tiny->radix) {
    return ERR_BLOOM_PARAMS_MISMATCH;
  }

  int h1, h2;
  int err = compute_hash(
      tiny, in, in_bytes, bloom->salt, bloom->salt_bytes, &h1, &h2);
  if (err != OK) {
    return err;
  }

  int res = 1;
  for (int j = 0; j < bloom->hash_ct; j++) {
    int i = (h1 + (j * h2)) % bloom->filter_bits;
    res = res && get_bit(bloom->filter, i);
  }

  return res;
}
