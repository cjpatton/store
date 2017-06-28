Implementations of (keyed) Bloom filters and dictionaries.

TODO
====

* Go interface
  * Wrap cdict_t in protobuf
  * RPC for serving queries from untrusted server.

* C library
  * Come up with better name than "struct". There may be a lot of collisions
    with libstruct.so and include/struct.
  * bits.c: More efficient get_chunk()?
  * Function for computing BF parameters.
