go-tiff2pdf [![GoDoc](https://godoc.org/github.com/companieshouse/go-tiff2pdf?status.svg)](https://godoc.org/github.com/companieshouse/go-tiff2pdf) [![Build Status](https://travis-ci.org/companieshouse/go-tiff2pdf.svg?branch=master)](https://travis-ci.org/companieshouse/go-tiff2pdf)
===========

tiff2pdf (from libtiff) as a service.

### Getting started

- Run `make deps` to download and install libtiff
- Run `make test` (optional: converts `t2p-test/tifs/*` to PDFs in `t2p-test/pdfs/`)
- Run `make` to build `go-tiff2pdf` library and service
- Run `./build/go-tiff2pdf` or `make run` to start the service

This has been tested on:
* Mac OS X 10.10 with Go 1.3.1 and Go 1.3.3
* Ubuntu 14.04 with Go 1.2.1

### TIFF to PDF request example

To convert a TIFF to PDF, `POST` the TIFF bytes as the request body to the `/convert` endpoint.

You can set PDF metadata using the headers:
* `PDF-Subject`
* `PDF-Author`
* `PDF-Creator`
* `PDF-Title`
* `PDF-PageSize`
* `PDF-FullPage`

Log messages for a request are linked together using a randomly generated unique ID.
You can set the `X-Request-ID` header to override the default random ID.

#### Example request

```
POST /convert HTTP/1.1
Content-Type: image/tiff
PDF-Subject: pdf subject line
PDF-Author: pdf author
PDF-Creator: pdf creator
PDF-Title: pdf title
PDF-PageSize: A4
PDF-FullPage: true
X-Request-ID: my-unique-id

[TIFF bytes]
```

#### Example response

```
HTTP/1.1 200 Ok
Content-Type: application/pdf
Content-Length: [n]
PDF-Pages: [n]

[PDF bytes]
```
