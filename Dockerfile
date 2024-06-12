# build a main binary using the golang container
FROM golang:1.22 as builder

WORKDIR /build/
COPY . .
RUN go build

# copy the main binary to a separate container based on ubuntu
FROM ubuntu:latest
WORKDIR /bin/
COPY --from=builder /build/tcping .
ENTRYPOINT [ "/bin/tcping" ]
CMD [ "-server" ]
EXPOSE 9000