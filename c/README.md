struct
======

Implementations of (keyed) data structures, including **Bloom filters** and
**Bloomier filters**.

To build and run tests:

```$ make && make test```

Note that many tests are probabilistic and will produce warnings from time to
time. (There's an issue if it produces **a lot** of warnings.) To install:

```$ sudo make install && sudo ldconfig```

This builds a file called libstruct.so and moves it to `/usr/local/lib` and copies
the header files to `/usr/local/include/struct`.

TODO
----
 * Rename. There's likely a lot of collisions with `libstruct.so` and
   `/usr/local/include/struct`.
 * `bits.c`: More efficient `get_chunk()`?
 * Function for computing BF parameters.
