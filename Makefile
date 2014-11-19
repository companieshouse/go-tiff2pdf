LIBTIFF_PATH=vadz/libtiff
all: service

lib:
	go build -work .

service:
	go build -work -o build/go-tiff2pdf ./tiff2pdf-service

test:
	go build -work -o build/t2p-test ./t2p-test

deps:
	-git clone git@github.com:$(LIBTIFF_PATH).git ../../$(LIBTIFF_PATH)
	cd ../../$(LIBTIFF_PATH) && ./configure --disable-pixarlog --disable-zlib && make

.PHONY: all lib service deps
