/* bits.h - Functions for manipulating arrays of bits.
 */
#ifndef BITS_H
#define BITS_H

// Computes ceil(log2(radix)) and returns the result.
//
// This is useful for determining the number of bits needed to encode a
// number x in range [0, radix).
int ceillog2(int radix);

// Sets a bit of 'bytes' to 1.
void set_bit(char *bytes, int idx);

// Sets a bit of 'bytes' to 0.
void unset_bit(char *bytes, int idx);

// Returns a bit of 'bytes'.
int get_bit(const char *bytes, int idx);

// Returns an integer corresponding to the i-th, k-bit chunk of 'bytes'. The
// first bit of the chunk is interpreted as the most significant bit of the
// output.
int get_chunk(const char *bytes, int k, int i);

#endif //BITS_H
