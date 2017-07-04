Good times with data structures
===============================


Package `store`
---------------
Package store provides secure storage of `map[string]string` objects. The
contents of the structure cannot be deduced from its public representation, and
querying it requires knowledge of a secret key. It is suitable for client/server
protocols where the service is trusted to provide storage, but is otherwise
not trusted.

The client possesses a secret key `K` and data `M` (of type `map[string]string`.
It executes:
```
pub, priv, err := store.New(K, M)
```

and transmits `pub`, the public representation of `M`, to the server.
To compute `M[input]`, the client executes:
```
x, y, err := priv.GetIdx(input)
```

and sends `x` and `y` to the server. These are integers corresponding to rows in
a table (of type `[][]byte`) encoded by `pub`. The server executes:
```
pubShare, err := pub.GetShare(x, y)
```

and sends `pubShare` (of type `[]byte`) to the client. (This is the XOR of the
two rows corresponding to `x` and `y`.) Finally, the client executes:
```
output, err := priv.GetValue(input, pubShare)
```

This combines `pubShare` with a private share computed from `input` and the key.
The result is `output = M[input]`.  Note that the server is not entrusted with
the key; its only job is to look up the rows of the table requested by the
client. The data structure is designed so that _no_ information about `input` or
`output` is leaked to any party not in possession of the secret key.

Note that the length of each `output` is limited to 60 bytes; see the Go
documentation for details.

The `StoreProvider` RPC service
-------------------------------
`store.proto` provides a simple [remote procedure
call](http://www.grpc.io/docs/quickstart/go.html) for requesting public shares.
The `user` computes `pub` from its map `M` and key `K` and provisions the
service provider (out-of-band) with `pub`.  The requests consists of the `user`
and the table rows `x` and `y`, and the response consists of the `pubShare`
computed from `x`, `y`, and the `pub` associated with `user`.

This simple RPC provides no authentication, so any *anyone* can get the *entire*
public store of *any* user. This is not a problem, however, as long as the
adversary doesn't know (or can't guess) `K`. But if `K` is derived from a
password, for example, then the contents of `pub` are susceptible to dictionary
attacks.

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

Next, the core data structures are implemented in C. (Naviagte to
`go/src/github.com/cjpatton/store/c/`.)  The `Makefile` compiles a shared object
for which the Go code has bindings. They depend on OpenSSL, so you'll need to
install this library in advance. On Ubuntu:
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
warnings from time to time. (This is OK as long ast it doesn't produce **a lot**
of warnings.) To install, do
```
$ sudo make install && sudo ldconfig
```

This builds a file called `libstruct.so` and moves it to `/usr/local/lib` and
copies the header files to `/usr/local/include/struct`.

Now you should be able to build the package. To run tests, do
```
$ go test github.com/cjpatton/store
```

Running the toy application
---------------------------
`hadee_server/hadee_server.go` implements the RPC service and serves a single
user. It takes as input the user name and a file containing the public store.
To run it, first generate a sample store by doing:
```
$ cd hadee_gen && go install && hadee_gen
```
It will prompt you for a "master password" used to derive a key, which is used
to generate the structure. This writes a file `store.pub` to the current
directory. (The map it represents is hard-coded in the Go code.) To run the
server, do:
```
$ cd hadee_server && go install && hadee_server cjpatton store.pub
```
This opens a TCP socket on localhost:50051 and begins serving requests. To run
the client, do:
```
$ cd hadee_client && go install && hadee_client cjpatton
```

**SECURITY WARNING:** Do NOT use this for anything real. As is, the protocol is
susceptible to dictionary attacks on the master password.

Modifying `store.proto`
----------------------
**You only need to do this if you want to modify the protocol buffers or RPC.**
This project uses protcool buffers and remote procedure calls. To build you'll
first need the lastest version of `protoc`. Go to [protobuf
documentation](https://developers.google.com/protocol-buffers/docs/gotutorial)
for instructions. To build `store.pb.go`, go to
`$HOME/go/src/github.com/cjpatton/store/` and run
```
  $ protoc -I . store.proto --go_out=plugins=grpc:.
```
Note that you only need to this if you modify `store.proto`.


Copyright notice
----------------
This software is distributed under the terms of the 3-Clause BSD License; see
`LICENSE` in the root project directory for details.
