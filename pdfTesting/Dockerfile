# To be built from the / directory using:
# docker build -f Dockerfile -t pdf_testing .
# docker run pdf_testing

FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN apk add --update --no-cache git

RUN go get rsc.io/pdf

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main .
WORKDIR .
CMD ["./main"]


# go clean -i -n github.com/motemen/gore