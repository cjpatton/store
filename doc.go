// Copyright (c) 2017, Christopher Patton
// All rights reserved.

/*
Package store provides secure storage of map[string]string objects. It combines
an AEAD ("authenticated encryption with associated data") scheme and a data
structure for representing functions called a Bloomier filter. It is suited
client/server protocols in which the server is trusted only to provide storage.
The client possesses a secret key K and data M (of type map[string]string). It
executes:

		pub, priv, err := store.NewStore(K, M)

and transmits pub, the public representation of M, to the server.
To compute M[input], the client executes:

		x, y, err := priv.GetIdx(input)

and sends x and y to the server. These are integers corresponding to rows in a
table encoded by pub. The server executes:

		pubShare, err := pub.GetShare(x, y)

and sends pubShare (of type []byte) to the client. (This is the XOR of the two
rows requested by the client.) Finally, the client executes:

		output, err := priv.GetOutput(input, pubShare)

This combines pubShare with the private share computed from input and the
key. The result is M[input]. Note that the server is not entrusted with the key;
its only job is to look up the rows of the table requested by the client. The
data structure is designed so that no information about input or output (except
for the length of the output) is leaked to any party not in possession of the
client's secret key.

At the core of data structure is a Bloomier filter, a variant of a technique of
Charles and Chellapilla for representing functions. (See "Bloomier Filters: A
second look", appearing at ESA 2008.) It is implemented in C, but this package
provides Go bindings. You can use it similarly to Store:

		pub, priv, err := store.NewDict(K, M)
		x, y, err := priv.GetIdx(input)
		pubShare, err := pub.GetShare(x, y)
		output, err := priv.GetOutput(input, pubShare)

However, the length of the outputs is limited to 60 bytes, a limitation we now
explain. The idea of Dict is that each input is mapped to two rows of a table so
that, when these rows are added together (i.e., their bitwise-XOR is computed),
the result is equal to the output. The look up is performed using hash
functions. The table has L rows and each row is 64 bytes in length. Roughly
speaking, the query is evaluated as follows:

		x := H1(input) // An integer in range [1..L]
		y := H2(input) // An integer in range [1..L]
		pad := H3(input) // A string of length 64 bytes
		output := Table[x] ^ Tabley[y] ^ pad

(Note that the above is pseudocode; the functions H1, H2, and H3 are not
provided.) In our setting, the hash functions are implemented using HMAC-SHA512,
which is a keyed, pseudorandom function. The output of HMAC-SHA512 is 64 bytes
in length.

If some query is not in the table, then the result of the query should indicate
as much. This is accomplished by appending a tag to the output. After adding up
the pad and table rows, we check if the last few bytes are 0; if so, then the
input/output pair is in the map; otherwise the input/output pair is not in the
map. 4 bytes of each row are allocated for the tag and padding of the output;
hence, each output must be at most 60 bytes long.  Note that this makes the data
structure probabilistic, since there is a small chance that, when the query is
evaluated, the tag bytes will all equal 0, even though the input is not correct.

NOTE: Dict does not on its own provide integrity protection, as Store does. It's
meant to be extremely light weight, and in fact is a core component of Store.

NOTE: The underlying data structure is implemented in C. The source can be found
in github.com/cjpatton/store/c; refer to github.com/cjpatton/store/README.md for
installation instructions.
*/
package store
