# To be built from the /Services directory using:
# docker build -f Dockerfile -t docker-test .
# docker run -p 8881:8881 docker-test

FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main .
WORKDIR .
CMD ["./main"]