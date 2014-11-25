#go-tiff2pdf

tiff2pdf (from libtiff) as a service.

### Getting started

- Run `make deps` to download and install libtiff
- Run `make test` (optional: converts `t2p-test/tifs/*` to PDFs in `t2p-test/pdfs/`)
- Run `make` to build `go-tiff2pdf` library and service
- Run `./build/go-tiff2pdf` or `make run` to start the service

This has been tested on:

    - Mac OS X 10.10 with Go 1.3.1
    - Ubuntu 14.04 with Go 1.2.1
