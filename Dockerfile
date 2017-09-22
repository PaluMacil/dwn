FROM golang:1.9.1 as goget
RUN go get .\..

FROM golang:1.9.1-alpine as gobuild
COPY --from=goget /go/src /go/src
RUN go install github.com/palumacil/dwn

FROM alpine
LABEL maintainer="dcwolf@gmail.com"
COPY --from=gobuild /go/bin /go/bin
# TODO: COPY dist from ng build
EXPOSE 80
CMD /go/bin/dwn