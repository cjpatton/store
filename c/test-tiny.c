#include "openssl/rand.h"
#include "stdio.h"
#include "string.h"

#include "bits.h"
#include "const.h"
#include "tiny.h"

int test_tinyprf_bad_new() {
  int ret = OK;
  const char test [] = "test_tinyprf_bad_new";

  tiny_ctx *tiny = tinyprf_new(1);
  ASSERT_FATAL("tinyprf_new(1) is not NULL", tiny == NULL);

  return ret;
}

// Tests tinyprf_init() and tinyprf().
//
// Returns OK if tests passed and ERR otherwise.
int test_tinyprf_one() {
  int ret = OK;
  const char test [] = "test_tinyprf_one";

  const char test_key[] = "abcd123412341234"; // must be KEY_BYTES long.
  const char test_in[] = "Don't YOU wish your input was this cool?";
  int radix = 700;
  tiny_ctx *tiny = tinyprf_new(radix);

  int err = tinyprf_init(tiny, test_key);
  ASSERT_OK_FATAL("tinyprf_init()", err);
  // FIXME callback
  //
  // tinyprf_free(tiny);
  ASSERT_EQ_ERROR("_use_prf", tiny->_use_prf, 1);
  ASSERT_EQ_ERROR("radix", tiny->radix, radix);
  ASSERT_EQ_ERROR("chunk_bits", tiny->chunk_bits, 10);
  ASSERT_EQ_ERROR("chunks", tiny->chunks, 51);

  // Check output.
  int res = tinyprf(tiny, test_in, strlen(test_in), NULL, 0);
  ASSERT_EQ_ERROR("tinyprf()", res, 436);

  res = tinyprf(tiny, test_in, strlen(test_in), test_in, strlen(test_in));
  ASSERT_EQ_ERROR("tinyhash()", res, 534);

  // Test prf()
  char buf [HASH_BYTES];
  err = prf(tiny, test_in, strlen(test_in), NULL, 0, buf, HASH_BYTES);
  const char expected_buf [] = "\x23\xe5\x1f\xfc\xbc\x60\xe4\xb9";
  ASSERT_OK_ERROR("prf()", err);
  ASSERT_EQ_ERROR("strncmp(buf, expected_buf, 8)", strncmp(buf, expected_buf, 8), 0);

  err = prf(tiny, test_in, strlen(test_in), NULL, 0, NULL, 0);
  ASSERT_OK_ERROR("second prf()", err);

  // Test again after reinitializing the context.
  err = tinyprf_init(tiny, test_key);
  ASSERT_OK_FATAL("second tinyprf_init()", err);
  // FIXME callback
  //
  // tinyprf_free(tiny);

  tinyprf_free(tiny);
  return ret;
}

// Tests multiple executions of tinyprf(), beginning with a random, 4-byte
// value and counting up.
//
// NOTE This test is flaky.
//
// Tests for systematic bias in the output. Assuming the output of tinyprf() is
// uniform, the expected number of outputs in each bin is '(double)trials /
// radix'. For each bin, check if the bin count differs from the expected value
// by more than 'err_tolerance'.
int test_tinyprf_many() {
  int ret = OK;
  const char test [] = "test_tinyprf_many";

  // Parameters for test.
  int radix = 1000, trials = 100000, err_tolerance = 0.00035*trials;
  tiny_ctx *tiny = tinyprf_new(radix);

  // Generate key and initialize context.
  int err = tinyprf_init_generate_key(tiny);
  ASSERT_OK_FATAL("tinyprf_init_generate_key()", err);
  // FIXME callback
  //
  // tinyprf_free(tiny);

  int *bins = malloc(radix * sizeof(int));
  memset(bins, 0, radix * sizeof(int));

  int bad = 0, seed;
  RAND_bytes((unsigned char *)&seed, sizeof(int));
  for (int trial = 0; trial < trials; trial++) {
    seed ++;
    int res = tinyprf(tiny, (const void *)&seed, sizeof(int), NULL, 0);
    ASSERT_FATAL("tinyprf() failed", res != ERR_SHA512);
    // FIXME callback
    //
    // free(bins);
    if (res == BAD) {
      bad ++;
    } else {
      bins[res] ++;
    }
  }
  if (bad > 0) {
    fprintf(stderr, "warning: %d bad tinyprf() outputs\n", bad);
  }

  int first = 0;
  RAND_bytes((unsigned char *)&first, sizeof(int));
  int skewed = 0;
  double expected_bin_ct = ((double)trials / radix);
  for (int ct = 0; ct < radix; ct++) {
    if ((bins[ct] - expected_bin_ct) > err_tolerance ||
        (expected_bin_ct - bins[ct]) > err_tolerance) {
      fprintf(stderr,
          "warning: tinyprf(): bin #%d count is %d, expected ~%0.1f\n",
          ct, bins[ct], expected_bin_ct);
      skewed ++;
    }
  }
  if (skewed > 0) {
    fprintf(stderr,
        "warning: tinyprf(): found %d possibly skewed bins\n", skewed);
  }

  tinyprf_free(tiny);
  free(bins);
  return ret;
}

int test_tinyhash_bad_new() {
  int ret = OK;
  const char test [] = "test_tinyhash_bad_new";

  tiny_ctx *tiny = tinyhash_new(1);
  ASSERT_FATAL("tinyhash_new(1) is not NULL", tiny == NULL);

  return ret;
}

// Tests tinyhash_init() and tinyhash().
//
// Returns OK if tests passed and ERR otherwise.
int test_tinyhash_one() {
  int ret = OK;
  const char test [] = "test_tinyhash_one";

  const char test_in[] = "Don't YOU wish your input was this cool?";
  int radix = 30000;
  tiny_ctx *tiny = tinyhash_new(radix);

  int err = tinyhash_init(tiny);
  ASSERT_OK_FATAL("tinyhash_init()", err);
  // FIXME callback
  //
  // tinyhash_free(tiny);
  ASSERT_EQ_ERROR("_use_prf", tiny->_use_prf, 0);
  ASSERT_EQ_ERROR("radix", tiny->radix, radix);
  ASSERT_EQ_ERROR("chunk_bits", tiny->chunk_bits, 15);
  ASSERT_EQ_ERROR("chunks", tiny->chunks, 34);

  // Check output.
  int res = tinyhash(tiny, test_in, strlen(test_in), NULL, 0);
  ASSERT_EQ_ERROR("tinyhash()", res, 8643);

  res = tinyhash(tiny, test_in, strlen(test_in), test_in, strlen(test_in));
  ASSERT_EQ_ERROR("tinyhash()", res, 23300);

  // Test tiny->g
  char buf [HASH_BYTES];
  err = hash(tiny, test_in, strlen(test_in), NULL, 0, buf, HASH_BYTES);
  const char expected_buf [] = "\x6b\x94\x47\xbc\x35\x17\x95\xcb";
  ASSERT_OK_ERROR("hash()", err);
  ASSERT_EQ_ERROR("strncmp(buf, expected_buf, 8)", strncmp(buf, expected_buf, 8), 0);

  err = hash(tiny, test_in, strlen(test_in), NULL, 0, NULL, 0);
  ASSERT_OK_ERROR("second hash()", err);

  // Test again after reinitializing the context.
  err = tinyhash_init(tiny);
  ASSERT_OK_FATAL("second tinyhash_init()", err);
  // FIXME callback
  //
  // tinyhash_free(tiny);

  tinyhash_free(tiny);
  return ret;
}

// Tests multiple executions of tinyhash(), beginning with a random, 4-byte
// value and counting up.
//
// NOTE This test is flaky.
//
// Tests for systematic bias in the output. Assuming the output of tinyhash() is
// uniform, the expected number of outputs in each bin is '(double)trials /
// radix'. For each bin, check if the bin count differs from the expected value
// by more than 'err_tolerance'.
int test_tinyhash_many() {
  int ret = OK;
  const char test [] = "test_tinyhash_one";

  // Parameters for test.
  int radix = 1000, trials = 100000, err_tolerance = 0.00035*trials;
  tiny_ctx *tiny = tinyhash_new(radix);

  // Generate key and initialize context.
  int err = tinyhash_init(tiny);
  ASSERT_OK_FATAL("tinyhash_init()", err);
  // FIXME callback
  //
  // tinyhash_free(tiny);

  int *bins = malloc(radix * sizeof(int));
  memset(bins, 0, radix * sizeof(int));

  int bad = 0, seed;
  RAND_bytes((unsigned char *)&seed, sizeof(int));
  for (int trial = 0; trial < trials; trial++) {
    seed ++;
    int res = tinyhash(tiny, (const void *)&seed, sizeof(int), NULL, 0);
    ASSERT_FATAL("tinyhash() failed", res != ERR_SHA512);
    // FIXME callback
    //
    // free(bins);
    if (res == BAD) {
      bad ++;
    } else {
      bins[res] ++;
    }
  }
  if (bad > 0) {
    fprintf(stderr, "warning: %d bad tinyhash() outputs\n", bad);
  }

  int first = 0;
  RAND_bytes((unsigned char *)&first, sizeof(int));
  int skewed = 0;
  double expected_bin_ct = ((double)trials / radix);
  for (int ct = 0; ct < radix; ct++) {
    if ((bins[ct] - expected_bin_ct) > err_tolerance ||
        (expected_bin_ct - bins[ct]) > err_tolerance) {
      fprintf(stderr,
          "warning: tinyhash(): bin #%d count is %d, expected ~%0.1f\n",
          ct, bins[ct], expected_bin_ct);
      skewed ++;
    }
  }
  if (skewed > 0) {
    fprintf(stderr,
        "warning: tinyhash(): found %d possibly skewed bins\n", skewed);
  }

  tinyhash_free(tiny);
  free(bins);
  return ret;
}

int main() {
  if (test_tinyprf_bad_new()==OK &&
      test_tinyprf_one()==OK &&
      test_tinyprf_many()==OK &&
      test_tinyhash_bad_new()==OK &&
      test_tinyhash_one()==OK &&
      test_tinyhash_many()==OK) {
    printf("pass\n");
    return 0;
  } else {
    return 1;
  }
}
