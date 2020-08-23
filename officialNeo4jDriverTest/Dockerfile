FROM golang:1.13-alpine3.11 AS builder

# add our cgo dependencies
RUN apk add --update --no-cache ca-certificates cmake make g++ openssl-dev git curl pkgconfig openssl-libs-static
# clone seabolt-1.7.0 source code
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
# invoke cmake build and install artifacts - default location is /usr/local
WORKDIR /seabolt/build
# CMAKE_INSTALL_LIBDIR=lib is a hack where we override default lib64 to lib to workaround a defect
# in our generated pkg-config file
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

# install dep
#RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/app
ADD . /go/src/app

# install dependencies
#RUN dep init

RUN go get github.com/neo4j/neo4j-go-driver/neo4j

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /go/src/app .
WORKDIR .
CMD ["./main"]