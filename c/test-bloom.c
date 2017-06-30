/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "math.h"
#include "stdio.h"
#include "string.h"

#include "bloom.h"
#include "const.h"
#include "tiny.h"

int test_new_free() {
  int ret = OK;
  const char test[] = "new_free";

  int salt_bytes = 16;
  int filter_bits = 1000;
  int filter_bytes = filter_bits / 8;
  if ((filter_bits % 8) > 0) {
    filter_bytes ++;
  }
  int hash_ct = 4;

  bloom_t *bloom = bloom_new(filter_bits, salt_bytes, hash_ct);
  ASSERT_EQ_ERROR("salt_bytes", bloom->salt_bytes, salt_bytes);
  ASSERT_EQ_ERROR("filter_bits", bloom->filter_bits, filter_bits);
  ASSERT_EQ_ERROR("filter_bytes", bloom->filter_bytes, filter_bytes);
  ASSERT_EQ_ERROR("hash_ct", bloom->hash_ct, hash_ct);

  bloom_free(bloom);
  return ret;
}

int test_init() {
  int ret = OK;
  const char test[] = "test_init";

  const char salt[] = "fella";
  bloom_t *bloom = bloom_new(1000, strlen(salt), 4);
  int err = bloom_init(bloom, salt);
  ASSERT_OK_ERROR("bloom_init()", err);
  ASSERT_EQ_ERROR("bloom_init(): strncmp(salt, bloom->salt, bloom->salt_bytes)",
    strncmp(salt, bloom->salt, bloom->salt_bytes), 0);

  char sum = 0x00;
  for (int i = 0; i < bloom->filter_bytes; i++) {
    sum |= bloom->filter[i];
  }
  ASSERT_ERROR("bloom_init(): filter uninitialized", sum == 0x00);

  bloom_free(bloom);
  return ret;
}

int test_init_generate_salt() {
  int ret = OK;
  const char test[] = "test_init";

  bloom_t *bloom = bloom_new(1000, 16, 4);
  int err = bloom_init_generate_salt(bloom);
  ASSERT_OK_ERROR("bloom_init_generate_salt()", err);

  char sum = 0x00;
  for (int i = 0; i < bloom->filter_bytes; i++) {
    sum |= bloom->filter[i];
  }
  ASSERT_ERROR("bloom_init(): filter uninitialized", sum == 0x00);

  bloom_free(bloom);
  return ret;
}

// NOTE Test flakes with a very low probability. (It tests for an element not in
// the set.)
int test_insert_get() {
  int ret = OK;
  const char test[] = "test_insert_get";

  int filter_bits = 31;
  int salt_bytes = 8;
  int hash_ct = 4;

  tiny_ctx *tiny = tinyhash_new(filter_bits);
  int err = tinyhash_init(tiny);
  ASSERT_OK_FATAL("tinyhash_init()", err);
  // FIXME callback
  //
  // tinyahsh_free(tiny);

  bloom_t *bloom = bloom_new(filter_bits, salt_bytes, hash_ct);
  err = bloom_init_generate_salt(bloom);
  ASSERT_OK_FATAL("bloom_init_generate_salt()", err);
  // FIXME callaback
  //
  // tinyhash_free(tiny);
  // bloom_free(bloom);

  const char in [] = "This is a great input.";
  err = bloom_insert(bloom, tiny, in, strlen(in));
  ASSERT_OK_FATAL("bloom_insert()", err);
  if (err != OK) {
    fprintf(stderr, "bloom_insert(): fiilter = ");
    for (int i = 0; i < filter_bits; i++) {
      fprintf(stderr, "%d", get_bit(bloom->filter, i));
    }
    fprintf(stderr, "\n");
  }

  int res = bloom_get(bloom, tiny, in, strlen(in));
  ASSERT_EQ_ERROR("bloom_get(in)", res, 1);

  res = bloom_get(bloom, tiny, "cool", strlen("cool"));
  ASSERT_EQ_ERROR("bloom_get(\"cool\")", res, 0);

  bloom_free(bloom);
  tinyhash_free(tiny);
  return ret;
}

int test_many() {
  int ret = OK;
  const char test[] = "test_many";

  int hash_ct = 4;
  int filter_bits = 1000;
  int size = 100;
  int salt_bytes = 16;

  int trials = 10000;
  double fp_rate = 0.01;
  double err_tolerance = 0.003;

  tiny_ctx *tiny = tinyprf_new(filter_bits);
  int err = tinyprf_init_generate_key(tiny);
  ASSERT_OK_FATAL("tinyprf_init_generate_key()", err);
  // FIXME callback
  //
  // tinyprf_free(tiny);

  bloom_t *bloom = bloom_new(filter_bits, salt_bytes, hash_ct);
  err = bloom_init_generate_salt(bloom);
  ASSERT_OK_FATAL("bloom_init_generate_salt()", err);
  // FIXME callaback
  //
  // tinyhash_free(tiny);
  // bloom_free(bloom);

  int in;
  for (in = 0; in < size; in++) {
    bloom_insert(bloom, tiny, (const char *)&in, sizeof(int));
  }

  int bad = 0;
  for (int trial = 0; trial < trials; trial++) {
    in++;
    int res = bloom_get(bloom, tiny, (const char *)&in, sizeof(int));
    bad += (res == 1);
  }

  double sim_fp_rate = (double)bad / trials;
  ASSERT_WARNING("false-positive rate higher than expected",
      fabs(sim_fp_rate - fp_rate) < err_tolerance);

  tinyprf_free(tiny);
  bloom_free(bloom);
  return ret;
}

int main() {
  if (test_init()==OK &&
      test_insert_get()==OK &&
      test_many()==OK) {
    printf("pass\n");
    return 0;
  } else {
    return 1;
  }

}
