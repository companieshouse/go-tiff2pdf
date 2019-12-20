TESTS ?= ./...
LIBTIFF_PATH=vadz/libtiff
LIBTIFF_REL=../../$(LIBTIFF_PATH)
TIFF2PDF_C=tiff2pdf/c/tiff2pdf.c
T2P_TEST_PATH=t2p-test
# the last known commit that worked with this build (June 2016)
LIBTIFF_COMMIT=c421b993abe1d6792252833c3bc8b3252b015fb9

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

.PHONY: all
all: build

.PHONY: lib
lib:
	go build -work .

$(TIFF2PDF_C): $(LIBTIFF_REL)/tools/tiff2pdf.c
	sed -e 's/^t2p_enable(/__not_&/' -e 's/^t2p_disable(/__not_&/' -e '/^int main(/,/^}/d' < $< > $@.tmp
	mv $@.tmp $@

.PHONY: build
build: deps $(TIFF2PDF_C)
	go build -work -o build/go-tiff2pdf ./tiff2pdf-service

.PHONY: run
run: build
	./build/go-tiff2pdf

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test $(TESTS) -run 'Unit' -coverprofile=coverage.out

.PHONY: test-integration
test-integration: deps $(TIFF2PDF_C)
	go build -work -o build/t2p-test ./$(T2P_TEST_PATH)
	test -d $(T2P_TEST_PATH)/tifs || mkdir $(T2P_TEST_PATH)/tifs
	test -d $(T2P_TEST_PATH)/pdfs || mkdir $(T2P_TEST_PATH)/pdfs
	if ! ls $(T2P_TEST_PATH)/tifs/* > /dev/null 2>&1; then echo To test, put sample TIFF files into $(T2P_TEST_PATH)/tifs/; false; fi
	cd $(T2P_TEST_PATH) && ../build/t2p-test
	echo See PDFs in $(T2P_TEST_PATH)/pdfs/

.PHONY: getdeps
getdeps:
	test -d $(LIBTIFF_REL) || git clone https://github.com/$(LIBTIFF_PATH).git $(LIBTIFF_REL) && cd $(LIBTIFF_REL) && git checkout $(LIBTIFF_COMMIT)

.PHONY: cleandeps
cleandeps:
	cd $(LIBTIFF_REL) && make distclean

.PHONY: configdeps
configdeps: getdeps
	cd $(LIBTIFF_REL) && ( test -f Makefile || ./configure --disable-pixarlog --disable-zlib )

.PHONY: deps
deps: configdeps
	cd $(LIBTIFF_REL) && make

.PHONY: clean
clean:
	rm -r build $(TIFF2PDF_C)
