LIBTIFF_PATH=vadz/libtiff
LIBTIFF_REL=../../$(LIBTIFF_PATH)
TIFF2PDF_C=tiff2pdf/c/tiff2pdf.c
all: build

lib:
	go build -work .

$(TIFF2PDF_C): $(LIBTIFF_REL)/tools/tiff2pdf.c
	sed -e 's/^t2p_enable(/__not_&/' -e 's/^t2p_disable(/__not_&/' -e '/^int main(/,/^}/d' < $< > $@.tmp
	mv $@.tmp $@
build: $(TIFF2PDF_C)
	go build -work -o build/go-tiff2pdf ./tiff2pdf-service
run: build
	./build/go-tiff2pdf

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

clean:
	rm -r build $(TIFF2PDF_C)

.PHONY: all lib build run test deps configdeps getdeps cleandeps clean
