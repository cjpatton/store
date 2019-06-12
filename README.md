Oblivous Go maps
================

This package provides secure storage of `map[string]string` objects. The
contents of the structure cannot be deduced from its public representation, and
querying it requires knowledge of a secret key. It is suitable for client/server
protocols where the service is trusted only to provide storage. In addition to
providing confidentiality, it allows the client to verify the integrity of the
server's responses.

An overview and installation instructions follow; the package documentation is
indexed on [GoDoc](http://godoc.org/github.com/cjpatton/store).

![#b8b8b8](https://placehold.it/15/b8b8b8/000000?text=+) **DISCLAIMER:** This
code is related to a research paper, which will be up on ePrint at some point.

![#b8b8b8](https://placehold.it/15/b8b8b8/000000?text=+) **FUTURE WORK:**
Currently this package only provides _immutable_ storage of maps, meaning once
you've created a data store, you can't change its contents. However, with some
modifications it should be possible to insert, remove, and update input/output
pairs securely. I'll be working on this next.

The `store` package
-------------------

The main package provides two data structures: **Store** and **Dict**. The
former offers _confidentiality_ and _integrity_ for `map[string]string` objects
with arbitrary-length inputs and outputs. Its security follows from the
combination of _authenticated encryption with associated data_
([AEAD](https://en.wikipedia.org/wiki/Authenticated_encryption)) and the latter
structure, which offers only _confidentiality_ and is only suitable for maps
who's outputs are of length _at most_ 60.

**Store.**
The client possesses a secret key `K` and data `M` (of type `map[string]string`).
It executes:
```
pub, priv, err := store.NewStore(K, M)
```

and transmits `pub`, the public representation of `M`, to the server.
To compute `M[input]`, the client executes:
```
x, y, err := priv.GetIdx(input)
```

and sends `x` and `y` (both of type `int`) to the server. The pair `(x,y)` is
called the _index_ and is used by the server to compute its share of the
output.  The server computes:
```
pubShare, err := pub.GetShare(x, y)
```
and sends `pubShare` (of type `[]byte`) to the client. Finally, the client
executes:
```
output, err := priv.GetOutput(input, pubShare)
```

This combines `pubShare` with a private share computed from `input` and `K`.
The result is `output = M[input]`.  Note that the server is not entrusted with
the key; its only job is to look up the index requested by the client. The
underlying data structure is designed so that _no_ information about `input` or
`output` is leaked to any party not in possession of the secret key.

For convenience, this package also provides an interface for querying `pub`
directly:
```
output, err := priv.Get(pub, input)
```

**Dict.**
This light-weight structure is the core of **Store.** The Go package is an
interface for the underlying C implementation.  It can be used in exactly the
same way as **Store**, but is only suitable for short (60 byte) outputs. See the
package documentation for an explanation of this limitation. To construct it,
the client executes:
```
pub, priv, err := store.NewDict(K, M)
```

The remaining functions are as above.

The `store/pb` package
----------------------
File `pb/store.proto` specifies a bare-bones [remote procedure
call](http://www.grpc.io/docs/quickstart/go.html) for the client and server
roles in the protocol above.

**The `StoreProvider` service.**
The `user` computes `pub` from its map `M` and key `K` and provisions the
service provider (out-of-band) with `pub`.  The request consists of the `user`
and `(x,y)`, and the response consists of the `pubShare` computed from `x`, `y`,
and `pub`.

This simple RPC provides no authentication of the user, so any *anyone* can get
the *entire* public store of *any* user. This is not a problem, however, as long
as the adversary doesn't know (or can't guess) `K`. But if `K` is derived from a
password, for example, then the contents of `pub` are susceptible to dictionary
attacks.

For documentation of this package, check out the
[GoDoc](http://godoc.org/github.com/cjpatton/store/pb) index.

Installation
------------
First, you'll need Go. To get the latest version on Ubuntu, do

```
$ sudo add-apt-repository ppa:longsleep/golang-backports
$ sudo apt update
$ sudo apt install golang-go
```

On Mac, download the [pkg](https://golang.org/dl/) and install it. Next, add the
following lines to the end of`.bashrc` on Ubuntu or `.bash_profile` on Mac:

```
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

In a new terminal, make the directory `$HOME/go`, go to the directory and type:
```
go get github.com/cjpatton/store
```
This downloads this repository and puts it in
`go/src/github.com/cjpatton/store`.

Next, the core data structures are implemented in C. (Navigate to
`go/src/github.com/cjpatton/store/c/`.)  The `Makefile` compiles a shared object
for which the Go code has bindings. They depend on OpenSSL (SHA512 and
HMAC-SHA512), so you'll need to install this library in advance. On Ubuntu:
```
$ sudo apt-get install libssl-dev
```
On Mac via Homebrew:
```
$ brew install openssl
```
(Homebrew puts the includes in `/usr/local/opt/openssl/include`, which is a
non-standard location. `Makefile` passes this directory to the compioler via
`-I`, so this shouldn't be a problem.) To build the C code and run tests do:
```
$ make && make test
```
Note that, since the data strucutres are probabilistic, the tests will produce
warnings from time to time. (This is OK as long as it doesn't produce **a lot**
of warnings.) To install, do
```
$ sudo make install && sudo ldconfig
```

This builds a file called `libstructsec.so` and moves it to `/usr/local/lib` and
copies the header files to `/usr/local/include/structsec`.

Now you should be able to build the package. To run tests, do
```
$ go test github.com/cjpatton/store
```

Running the toy application
---------------------------
`hadee/server/hadee_server.go` implements the RPC service and serves a single
user. It takes as input the user name and a file containing the public store.
To run it, first generate a sample store by doing:
```
$ cd hadee/gen && go install && hadee_gen
```
It will prompt you for a "master password" used to derive a key, which is used
to generate the structure. This writes a file `store.pub` to the current
directory. (The map it represents is hard-coded in the Go code.) To run the
server, do:
```
$ cd hadee/server && go install && hadee_server cjpatton store.pub
```
This opens a TCP socket on localhost:50051 and begins serving requests. To run
the client, do:
```
$ cd hadee/client && go install && hadee_client cjpatton
```

![#f03c15](https://placehold.it/15/f03c15/000000?text=+) **SECURITY WARNING:**
Do NOT use this for anything real. As is, the protocol is susceptible to
dictionary attacks on the master password.

Modifying `store.proto`
----------------------
**You only need to do this if you want to modify the protocol buffers or RPC.**
This project uses protcool buffers and remote procedure calls. To build you'll
first need the lastest version of `protoc`. Go to [protobuf
documentation](https://developers.google.com/protocol-buffers/docs/gotutorial)
for instructions. To build `store.pb.go`, go to
`$HOME/go/src/github.com/cjpatton/store/pb` and run
```
  $ protoc -I . store.proto --go_out=plugins=grpc:.
```
Note that you only need to do this if you modify `store.proto`.

Copyright notice
----------------
This software is distributed under the terms of the 3-Clause BSD License; see
`LICENSE` in the root project directory for details.
