all: service

lib:
	go build -work .

service:
	go build -work -o build/go-tiff2pdf ./tiff2pdf-service

test:
	go build -work -o build/t2p-test ./t2p-test
	cd t2p-test; ../build/t2p-test

deps:
	-git clone git@github.com:vadz/libtiff.git ../../vadz/libtiff
	pushd ../../vadz/libtiff; ./configure --disable-pixarlog --disable-zlib; make; popd

.PHONY: all lib service deps
