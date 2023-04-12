FROM golang:1.4

RUN apt-get update
RUN apt-get install -y g++

WORKDIR /go/src
RUN git clone https://github.com/companieshouse/go-tiff2pdf github.com/companieshouse/go-tiff2pdf

WORKDIR /go/src/github.com/companieshouse/go-tiff2pdf
RUN make
RUN go install ./tiff2pdf-service

EXPOSE 9090

ENTRYPOINT ["/go/bin/tiff2pdf-service"]
