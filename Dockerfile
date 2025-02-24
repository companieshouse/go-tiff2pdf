FROM 416670754337.dkr.ecr.eu-west-2.amazonaws.com/ci-golang-build-1.23:latest

RUN yum update \
    && yum install -y g++

WORKDIR /go/src

RUN git clone https://github.com/companieshouse/go-tiff2pdf . \
    && git checkout main \
    && make
    
EXPOSE 9090
CMD ["./build/go-tiff2pdf"]
