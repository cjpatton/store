# Add to to the end of cc when building on Mac.
#
# FIXME This line is a hack.
mac = -I/usr/local/opt/openssl/include

cc = gcc -Wall -ansi -fpic -std=c11
lflags = -lcrypto

tinytest = test-tiny
bitstest = test-bits
bloomtest = test-bloom
dicttest = test-dict

libstruct.so: bits.o tiny.o bloom.o dict.o
	$(cc) -shared -o libstruct.so bits.o tiny.o bloom.o dict.o

install: libstruct.so const.h bits.h tiny.h bloom.h dict.h
	cp libstruct.so /usr/local/lib
	mkdir /usr/local/include/struct
	cp const.h bits.h tiny.h bloom.h dict.h /usr/local/include/struct

test: $(bitstest) $(tinytest) $(bloomtest) $(dicttest)

$(dicttest): $(dicttest).out
	./$(dicttest).out

$(bloomtest): $(bloomtest).out
	./$(bloomtest).out

$(tinytest): $(tinytest).out
	./$(tinytest).out

$(bitstest): $(bitstest).out
	./$(bitstest).out

$(dicttest).out: dict.o bits.o tiny.o $(dicttest).c
	$(cc) $(dicttest).c bits.o tiny.o dict.o -o $(dicttest).out $(lflags)

$(bloomtest).out: bloom.o bits.o tiny.o $(bloomtest).c
	$(cc) $(bloomtest).c bits.o tiny.o bloom.o -o $(bloomtest).out $(lflags)

$(tinytest).out: bits.o tiny.o $(tinytest).c
	$(cc) $(tinytest).c bits.o tiny.o -o $(tinytest).out $(lflags)

$(bitstest).out: bits.o $(bitstest).c
	$(cc) $(bitstest).c bits.o -o $(bitstest).out $(lflags)

dict.o: dict.h dict.c const.h tiny.h
	$(cc) -c dict.c

bloom.o: bloom.h bloom.c const.h tiny.h bits.h
	$(cc) -c bloom.c

tiny.o: tiny.h tiny.c const.h
	$(cc) -c tiny.c

bits.o: bits.h bits.c const.h
	$(cc) -c bits.c

clean:
	rm -f *.so *.o $(tinytest).out $(tinyhashtest).out $(bitstest).out $(bloomtest).out $(dicttest).out
