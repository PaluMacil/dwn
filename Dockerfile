FROM golang as build-server
ADD . /go/src/dwn
WORKDIR /go/src/dwn
RUN go build -o dwn-server

FROM alpine
COPY --from=build-server /go/src/dwn/dwn-server .

CMD ["dwn-server", "prod"]