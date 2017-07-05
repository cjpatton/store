// Copyright (c) 2017, Christopher Patton
// All rights reserved.

/*
Package store provides private storage of map[string]string objects. The set of
inputs are strings of arbitrary length and the set of outputs are strings of
length at most 60 bytes.

The data structure is designed for client/server protocols in which the server
is trusted only to provide storage. The client possesses a secret key K and
data M (of type map[string]string. It executes:

		pub, priv, err := store.New(K, M)

and transmits pub, the public representation of M, to the server.
To compute M[input], the client executes:

		x, y, err := priv.GetIdx(input)

and sends x and y to the server. These are integers corresponding to rows in a
table encoded by pub. The server executes:

		pubShare, err := pub.GetShare(x, y)

and sends pubShare (of type []byte) to the client. (This is the XOR of the two
rows requested by the client.) Finally, the client executes:

		output, err := priv.GetValue(input, pubShare)

This combines `pubShare` with the private share computed from `input` and the
key. The result is M[input]. Note that the server is not entrusted with the key;
its only job is to look up the rows of the table requested by the client. The
data structure is designed so that no information about input or output is
leaked to any party not in possession of the secret key. For convenience, we
also provide a means for the client to query the structure directly:

		output, err := store.Get(pub, priv, input)

The core structure is a variant of a technique of Charles and Chellapilla for
representing functions. (See "Bloomier Filters: A second look", appearing at ESA
2008.) The idea is that each input is mapped to two rows of a table so that,
when these rows are added together (i.e., their bitwise-XOR is computed), the
result is equal to the output. The look up is performed using hash functions.
The table has L rows and each row is 64 bytes in length. Roughly speaking, the
query is evaluated as follows:

		x := H1(input) // An integer in range [1..L]
		y := H2(input) // An integer in range [1..L]
		pad := H3(input) // A string of length 64 bytes
		output := Table[x] ^ Tabley[y] ^ pad

(Note that the above is pseudocode.) In our setting, the hash functions are
implemented using HMAC-SHA512, which is a keyed, pseudorandom function. The
output of HMAC-SHA512 is 64 bytes in length, hence the restriction on the length
on the table rows.

If some query is not in the table, then the result of the query should indicate
as much. This is accomplished by appending a tag to the output. After adding up
the pad and table rows, we check if the last few bytes are 0; if so, then the
input/output pair is in the map; otherwise the input/output pair is not in the
map. 4 bytes of each row are allocated for the tag and padding of the output;
hence, each output must be at most 60 bytes long.  Note that this makes the data
structure probabilistic, since there is a small chance that, when the query is
evaluated, the tag bytes will all equal 0, even though the input is not correct.

WARNIG: This package provides private storage of maps, but does not provide
integrity. Said another way, there's no way to detect if the service provider
has manipulated the contents of the map. Doing so, however, does not compromise
privacy.

NOTE: The underlying data structure is implemented in C. The source can be found
in github.com/cjpatton/store/c; refer to github.com/cjpatton/store/README.md for
installation instructions.
*/
package store
