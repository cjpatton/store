# Add to to the end of cc when building on Mac.
mac_cppflags = -I/usr/local/opt/openssl/include
cc = gcc -Wall -ansi -fpic -std=c11 $(mac_cppflags)
lflags = -lcrypto

tinytest = test-tiny
bitstest = test-bits
bloomtest = test-bloom
dicttest = test-dict

libstructsec.so: bits.o tiny.o bloom.o dict.o
	$(cc) -shared -o libstructsec.so bits.o tiny.o bloom.o dict.o $(lflags)

install: libstructsec.so const.h bits.h tiny.h bloom.h dict.h
	cp libstructsec.so /usr/local/lib
	mkdir -p /usr/local/include/structsec
	cp const.h bits.h tiny.h bloom.h dict.h /usr/local/include/structsec

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
