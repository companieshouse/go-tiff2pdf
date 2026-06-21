# libtiff is pulled from the maintained upstream release (not the abandoned
# vadz/libtiff GitHub mirror, whose newest tag is from 2017). Pin to a tagged
# release tarball and verify its SHA256 so the build is reproducible and the
# library carries current security fixes.
LIBTIFF_VERSION=4.7.1
LIBTIFF_SHA256=b92017489bdc1db3a4c97191aa4b75366673cb746de0dce5d7a749d5954681ba
LIBTIFF_TARBALL=tiff-$(LIBTIFF_VERSION).tar.xz
LIBTIFF_URL=https://download.osgeo.org/libtiff/$(LIBTIFF_TARBALL)
# extracted, repo-local (git-ignored) libtiff source tree
LIBTIFF_REL=libtiff-src
TIFF2PDF_C=tiff2pdf/c/tiff2pdf.c
T2P_TEST_PATH=t2p-test

# Only the codecs that need no external libraries are enabled, matching the
# original build (which linked nothing beyond -lm). Disabling the rest keeps
# the build self-contained and deterministic regardless of what dev/CI hosts
# happen to have installed.
LIBTIFF_CONFIGURE_FLAGS=--disable-jpeg --disable-old-jpeg --disable-lzma \
	--disable-jbig --disable-zstd --disable-lerc --disable-webp \
	--disable-pixarlog --disable-zlib

all: build

lib:
	CGO_ENABLED=1 go build -work .

$(TIFF2PDF_C): $(LIBTIFF_REL)/tools/tiff2pdf.c
	# Copy the tiff2pdf tool into the build, dropping its main() and routing
	# the output enable/disable toggle to the Go I/O layer (see tif_golang.c).
	sed -e '/^int main(/,/^}/d' \
	    -e '/^#include "libport.h"/d' \
	    -e 's/t2p->outputdisable = 1;/GoOutputDisable((int)(intptr_t)t2p);/' \
	    -e 's/t2p->outputdisable = 0;/GoOutputEnable((int)(intptr_t)t2p);/' \
	    < $< > $@.tmp
	mv $@.tmp $@

build: deps $(TIFF2PDF_C)
	CGO_ENABLED=1 go build -work -o build/go-tiff2pdf ./tiff2pdf-service
run: build
	./build/go-tiff2pdf

test: deps $(TIFF2PDF_C)
	CGO_ENABLED=1 go build -work -o build/t2p-test ./$(T2P_TEST_PATH)
	test -d $(T2P_TEST_PATH)/tifs || mkdir $(T2P_TEST_PATH)/tifs
	test -d $(T2P_TEST_PATH)/pdfs || mkdir $(T2P_TEST_PATH)/pdfs
	if ! ls $(T2P_TEST_PATH)/tifs/* > /dev/null 2>&1; then echo To test, put sample TIFF files into $(T2P_TEST_PATH)/tifs/; false; fi
	cd $(T2P_TEST_PATH) && ../build/t2p-test
	echo See PDFs in $(T2P_TEST_PATH)/pdfs/

getdeps:
	test -f $(LIBTIFF_REL)/libtiff/tiffio.h || ( \
	    curl -fsSL -o $(LIBTIFF_TARBALL) $(LIBTIFF_URL) && \
	    echo "$(LIBTIFF_SHA256)  $(LIBTIFF_TARBALL)" | \
	        (command -v sha256sum >/dev/null 2>&1 && sha256sum -c - || shasum -a 256 -c -) && \
	    rm -rf $(LIBTIFF_REL) && mkdir -p $(LIBTIFF_REL) && \
	    tar xf $(LIBTIFF_TARBALL) -C $(LIBTIFF_REL) --strip-components=1 && \
	    rm -f $(LIBTIFF_TARBALL) )
cleandeps:
	rm -rf $(LIBTIFF_REL)
configdeps:
	cd $(LIBTIFF_REL) && ( test -f libtiff/tif_config.h || ./configure $(LIBTIFF_CONFIGURE_FLAGS) )
deps: configdeps
	$(MAKE) -C $(LIBTIFF_REL)/libtiff

clean:
	rm -r build $(TIFF2PDF_C)

.PHONY: all lib build run test deps configdeps cleandeps clean
