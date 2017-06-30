Good times with data structures
===============================

A super light-weight password management solution for paranoid people who like
Go.

store
-----

Package store provides secure storage of `map[string]string` objects. The
contents of the structure cannot be deduced from its public representation,
and querying it requires knowledge of a secret key. It is suitable for
settings in which the service is trusted to provide storage, but nothing
else.

The client possesses a secret key `K` and data `M` (of type
`map[string]string`. It executes:

```pub, priv, err := store.New(K, M)```

and transmits `pub`, the public representation of `M`, to the server.
To compute `M[input]`, the client executes:

```x, y, err := priv.GetIdx(input)```

and sends `x` and `y` to the server. The are integers corresponding to rows
in a table (of type `[][]byte`) encoded by `pub`. The server executes:

```xRow, err := pub.GetRow(x)```

and sends `xRow` (of type `[]byte`) to the client. (Similarly for `y`.)
Finally, the client executes:

```output, err := priv.GetValue(input, [][]byte{xRow, yRow})```

Then `M[input] = output`.  Note that the server is not entrusted with the key;
it only looks up the rows of the table requested by the client. The data
structure is designed so that no information about `input` or `output` is leaked
to any party not in possession of the secret key. For convenience, we also
provide a means for the client to query the structure directly:

```output, err := store.Get(pub, priv, input)```


Installation
------------

The core data structures are implemented in C. (See `c/`.) These are compiled
into a share object for which the Go code has bindings. To begin, you'll need
to install the C To build and run tests, go to `c/` and do

```$ make && make test```

Note that since the data structure are probabilistic, the tests will produce
warnings from time to time. (This OK as long ast it doesn't produce **a lot** of
warnings.) To install, do

```$ sudo make install && sudo ldconfig```

This builds a file called libstruct.so and moves it to `/usr/local/lib` and copies
the header files to `/usr/local/include/struct`.

Finally, you'll need Go. To get the latest version on Ubuntu, do

```
$ sudo add-apt-repository ppa:longsleep/golang-backports
$ sudo apt-get update
$ sudo apt-get install golang-go
```
