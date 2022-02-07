# docker build -f Dockerfile -t personWebsite .
# docker run -p 80:80 personWebsite
FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main .
COPY /src src
COPY /static static
WORKDIR .
#CMD ["./main"]

EXPOSE 80