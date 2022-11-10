FROM golang:1.19-alpine AS builder

WORKDIR /go/src/github.com/thewizardplusplus/go-upload-progress-backend
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go install -a -ldflags='-w -s -extldflags "-static"' ./...

FROM scratch

COPY --from=builder \
  /go/bin/go-upload-progress-backend \
  /usr/local/bin/go-upload-progress-backend

CMD ["go-upload-progress-backend"]
