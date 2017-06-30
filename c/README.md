struct
======

Implementations of (keyed) data structures, including **Bloom filters** and
**Bloomier filters**. See `README.md` in the root project directory for
installation instructions.


TODO
----
 * Rename. There's likely a lot of collisions with `libstruct.so` and
   `/usr/local/include/struct`.
 * `bits.c`: More efficient `get_chunk()`?
 * Function for computing BF parameters.
