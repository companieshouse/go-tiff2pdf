LIBTIFF_PATH=vadz/libtiff
LIBTIFF_REL=../../$(LIBTIFF_PATH)
TIFF2PDF_C=tiff2pdf/c/tiff2pdf.c
T2P_TEST_PATH=t2p-test
# the last known commit that worked with this build (June 2016)
LIBTIFF_COMMIT=c421b993abe1d6792252833c3bc8b3252b015fb9
all: build

lib:
	go build -work .

$(TIFF2PDF_C): $(LIBTIFF_REL)/tools/tiff2pdf.c
	sed -e 's/^t2p_enable(/__not_&/' -e 's/^t2p_disable(/__not_&/' -e '/^int main(/,/^}/d' < $< > $@.tmp
	mv $@.tmp $@
build: deps $(TIFF2PDF_C)
	CGO_ENABLED=1 GOOS=linux go build -work -o build/go-tiff2pdf ./tiff2pdf-service
run: build
	./build/go-tiff2pdf

test: deps $(TIFF2PDF_C)
	go build -work -o build/t2p-test ./$(T2P_TEST_PATH)
	test -d $(T2P_TEST_PATH)/tifs || mkdir $(T2P_TEST_PATH)/tifs
	test -d $(T2P_TEST_PATH)/pdfs || mkdir $(T2P_TEST_PATH)/pdfs
	if ! ls $(T2P_TEST_PATH)/tifs/* > /dev/null 2>&1; then echo To test, put sample TIFF files into $(T2P_TEST_PATH)/tifs/; false; fi
	cd $(T2P_TEST_PATH) && ../build/t2p-test
	echo See PDFs in $(T2P_TEST_PATH)/pdfs/

getdeps:
	test -d $(LIBTIFF_REL) || git clone https://github.com/$(LIBTIFF_PATH).git $(LIBTIFF_REL) && cd $(LIBTIFF_REL) && git checkout $(LIBTIFF_COMMIT)
cleandeps:
	cd $(LIBTIFF_REL) && make distclean
configdeps: getdeps
	cd $(LIBTIFF_REL) && ( test -f Makefile || ./configure --disable-pixarlog --disable-zlib)
deps: configdeps
	cd $(LIBTIFF_REL) && make

clean:
	rm -r build $(TIFF2PDF_C)

# .PHONY: all lib build run test deps configdeps getdeps cleandeps clean
.PHONY: all lib build run test deps configdeps getdeps cleandeps
