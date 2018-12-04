FROM golang:1.10.5-alpine AS build
RUN apk add --no-cache git && go get github.com/golang/dep/cmd/dep
ENV GOPATH /go
ENV PATH $PATH:/go/bin
COPY Gopkg.lock Gopkg.toml /go/src/xsurl/
WORKDIR /go/src/xsurl/api
RUN dep ensure -vendor-only
COPY . /go/src/xsurl/api/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goapp ./cmd/shortenurl


FROM scratch
WORKDIR /app
COPY --from=build /go/src/xsurl/api/goapp .
ENTRYPOINT ["./goapp"]
# CMD ["--help"]
