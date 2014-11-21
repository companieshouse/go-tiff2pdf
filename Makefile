LIBTIFF_PATH=vadz/libtiff
LIBTIFF_REL=../../$(LIBTIFF_PATH)
all: service

lib:
	go build -work .

service:
	go build -work -o build/go-tiff2pdf ./tiff2pdf-service

test:
	go build -work -o build/t2p-test ./t2p-test
	cd t2p-test && ../build/t2p-test

getdeps:
	test -d $(LIBTIFF_REL) || git clone git@github.com:$(LIBTIFF_PATH).git $(LIBTIFF_REL)
cleandeps:
	cd $(LIBTIFF_REL) && make distclean
configdeps: getdeps
	cd $(LIBTIFF_REL) && ( test -f Makefile || ./configure --disable-pixarlog --disable-zlib )
deps: configdeps
	cd $(LIBTIFF_REL) && make

.PHONY: all lib service test deps configdeps getdeps cleandeps
