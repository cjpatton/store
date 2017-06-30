/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "bits.h"

int ceillog2(int radix) {
  int y = radix, n = 0;
  while ((y >> 1) > 0) {
    y >>= 1; n ++;
  }
  if ((1 << n) != radix) {
    n ++;
  }
  return n;
}

void set_bit(char *bytes, int idx) {
  int byte = idx >> 3;
  int bit = (byte << 3) ^ idx;
  bytes[byte] |= 1 << bit;
}

void unset_bit(char *bytes, int idx) {
  int byte = idx >> 3;
  int bit = (byte << 3) ^ idx;
  bytes[byte] &= ~(1 << bit);
}

int get_bit(const char *bytes, int idx) {
  int byte = idx >> 3;
  int bit = (byte << 3) ^ idx;
  return (bytes[byte] & (1 << bit)) > 0;
}

// TODO See if there's a way to do this more efficiently.
int get_chunk(const char *bytes, int k, int i) {
  int chunk = 0;
  for (int j = 0; j < k; j++) {
    chunk = (chunk<<1) | get_bit(bytes, (i*k)+j);
  }
  return chunk;
}
