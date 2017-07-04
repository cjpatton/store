struct
======

Implementations of (keyed) data structures, including a [Bloom
filters](https://en.wikipedia.org/wiki/Bloom_filter) (`bloom.{h,c}`) and
[Bloomier filters](https://arxiv.org/abs/0807.0928) (`dicet.{h,c}`).

See `README.md` in the root project directory for installation instructions.

TODO(cjpatton) Give details about tiny, bloom, and dict and design rationale.

TODO
----
 * Rename. There's likely a lot of collisions with `libstruct.so` and
   `/usr/local/include/struct`.
 * Function for computing BF parameters.
