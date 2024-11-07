FROM 416670754337.dkr.ecr.eu-west-2.amazonaws.com/ci-golang-build-1.23:latest

RUN yum update
RUN yum install -y g++

WORKDIR /go/src
RUN go version
RUN git clone https://github.com/companieshouse/go-tiff2pdf github.com/companieshouse/go-tiff2pdf

WORKDIR /go/src/github.com/companieshouse/go-tiff2pdf
RUN git checkout feature/CC-144
RUN make

# EXPOSE 9090
CMD ["./build/go-tiff2pdf", "-bind-addr=0.0.0.0:9090"]
