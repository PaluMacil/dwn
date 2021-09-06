FROM golang as build-server
ADD . /go/src/dwn
WORKDIR /go/src/dwn
RUN go build -o dwn-server

FROM ubuntu
COPY --from=build-server /go/src/dwn/dwn-server /opt/dwn/dwn-server

CMD ["/opt/dwn/dwn-server", "prod"]