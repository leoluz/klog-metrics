FROM golang:1.11.0-alpine3.8

RUN apk --no-cache add git curl gcc musl-dev

ENV PROJECT_NAME github.com/AppDirect/kubelogs

WORKDIR /go/src/${PROJECT_NAME}
ADD . /go/src/${PROJECT_NAME}

#RUN go test ./... -cover

RUN go build -v

# Runtime
FROM alpine:3.8

ENV PROJECT_NAME github.com/AppDirect/kubelogs

WORKDIR /app/
COPY --from=0 /go/src/${PROJECT_NAME}/kubelogs .

ENTRYPOINT ["/app/kubelogs"]

