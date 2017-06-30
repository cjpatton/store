/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "stdio.h"
#include "stdlib.h"

#include "bits.h"
#include "const.h"

int test_get_chunk() {
  int n = 512;
  int k = 10;
  int t = n/k;

  // Set up a byte array for testing.
  //
  // Chunks alternate between 1^k and 0^k.
  char *bytes = malloc(n/8*sizeof(char));

  int ret = OK;
  for (int i = 0; i < t; i++) { // Chunk
    for (int j = 0; j < k; j++) { // Chunk bit
      if ((i%2) == 0) { // Even chunks get 0^k
        set_bit(bytes, (i*k)+j);
      } else { // Odd chunks get 1^k
        unset_bit(bytes, (i*k)+j);
      }
    }
    int chunk = get_chunk(bytes, k, i);
    if ((i%2) == 0 && chunk != (1<<k)-1) {
      fprintf(stderr, "error: test_get_chunk() fails on even chunk\n");
      ret = ERR;
      break;
    } else if ((i%2) == 1 && chunk != 0) {
      fprintf(stderr, "error: test_get_chunk() fails on odd chunk\n");
      ret = ERR;
      break;
    }
  }
  free(bytes);
  return ret;
}

int main() {
  if (test_get_chunk()==OK) {
    printf("pass\n");
    return 0;
  } else {
    return 0;
  }
}
