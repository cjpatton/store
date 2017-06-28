/* bloom.h - A Bloom filter based on the double hashing scheme of Kirsch and
 * Mitzenmacher "Less hashing, same performance: Building a better bloom
 * filter." It uses tinyhash or tinyprf for (keyed) hashingo.
 *
 * TODO Add interfaces for insert/elemof using tinyprf and tinyhash.
 */

#ifndef BLOOM_H
#define BLOOM_H

#include "bits.h"
#include "tiny.h"

// The Bloom filter data structure, consisting of the filter array, salt, and
// paraemters.
typedef struct {

  char *filter, *salt;

  int hash_ct,      // Number of BF "hashes" (k)
      filter_bits,  // Filter length (in bits) (m)
      salt_bytes,   // Salt length (in byte)
      filter_bytes; // = ceil(filter_bits / 8)

} bloom_t;

// Creates a new bloom_new* and returns it.
//
// Must be destroyed with bloom_free().
bloom_t *bloom_new(int filter_bits, int salt_bytes, int hash_ct);

// Destroys BF context.
void bloom_free(bloom_t *ctx);

// Sets ctx->salt to 'salt' and zeros out 'filter'.
int bloom_init(bloom_t *ctx, const char *salt);

// Generates and stores a new salt of length 'ctx' and zeros out 'filter'.
int bloom_init_generate_salt(bloom_t *ctx);

// Inserts an element into the filter.
//
// Parameters:
//  ctx -- BF cotnext
//  hash_ctx -- Context for hashing, needed by 'hash'
//  hash -- Either tinyhash or tinyprf
//  in -- The input
//  in_bytes -- The length of the input
//
// Returns:
//  - OK if success,
//  - ERR_BLOOM_PARAMS_MISMATCH if bloom->filter_bits != tiny->radix,
int bloom_insert(
    bloom_t *bloom,
    tiny_ctx *tiny,
    const char *in, int in_bytes);

// Checks if the input is in the filter.
//
// Parameters:
//  ctx -- BF cotnext
//  hash_ctx -- Context for hashing, needed by 'hash'
//  hash -- Either tinyhash or tinyprf
//  in -- The input
//  in_bytes -- The length of the input
//
// Returns:
//  - 1 if found,
//  - 0 if not found, or
//  - an error.
int bloom_get(
    bloom_t *bloom,
    tiny_ctx *tiny,
    const char *in, int in_bytes);

#endif //BLOOM_H
