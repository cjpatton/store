// Copyright (c) 2017, Christopher Patton
// All rights reserved.

/*
Package store provides secure storage of `map[string]string` objects. The
contents of the structure cannot be deduced from its public representation,
and querying it requires knowledge of a secret key. It is suitable for
settings in which the service is trusted to provide storage, but nothing
else.

The client possesses a secret key `K` and data `M` (of type
`map[string]string`. It executes:

		pub, priv, err := store.New(K, M)

and transmits `pub`, the public representation of `M`, to the server.
To compute `M[input]`, the client executes:

		x, y, err := priv.GetIdx(input)

and sends `x` and `y` to the server. The are integers corresponding to rows
in a table (of type `[][]byte`) encoded by `pub`. The server executes:

		xRow, err := pub.GetRow(x)

and sends `xRow` (of type `[]byte`) to the client. (Similarly for `y`.)
Finally, the client executes:

		output, err := priv.GetValue(input, [][]byte{xRow, yRow})

and `M[input] = output`.

Note that the server is not entrusted with the key; it simply looks up the
rows of table requested by the client. The data structure is designed so that
no information about `input` or `output` is leaked any party not in
possession of the secret key.

For convenience, we also provide a means for the client to query the
structure directly:

		output, err := store.Get(pub, priv, input)
*/
package store
