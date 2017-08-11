FROM golang:1.8.3 as goget
LABEL maintainer="dcwolf@gmail.com"
RUN go get .\..

FROM golang:1.8.3-alpine as gobuild
LABEL maintainer="dcwolf@gmail.com"
COPY --from=goget /go/src /go/src
RUN go install github.com/palumacil/dwn

FROM alpine
LABEL maintainer="dcwolf@gmail.com"
COPY --from=gobuild /go/bin /go/bin
# COPY dist
EXPOSE 80
CMD /go/bin/dwn