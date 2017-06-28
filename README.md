Implementations of (keyed) Bloom filters and dictionaries.

TODO
====

* Go interface
  * Wrap cdict_t in protobuf
  * RPC for serving queries from untrusted server.
* bits.c: More efficient get_chunk()?
* Function for computing BF parameters.
