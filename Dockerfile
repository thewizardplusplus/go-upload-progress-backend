# syntax=docker/dockerfile:1.2

FROM golang:1.19-alpine AS builder

WORKDIR /go/src/github.com/thewizardplusplus/go-upload-progress-backend
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go install -a -ldflags='-w -s -extldflags "-static"' ./...

FROM scratch

# Problem: https://blog.cubieserver.de/2020/go-debugging-why-parsemultipartform-returns-error-no-such-file-or-directory/
# Solution: https://www.reddit.com/r/docker/comments/8y2zyx/comment/hneaaqe/
# Note: https://gist.github.com/thaJeztah/cfd929a31976b745e3f7515ae37eb192?permalink_comment_id=4397393#gistcomment-4397393
RUN --mount=from=busybox:1.36-uclibc,source=/bin,target=/bin \
  mkdir -m 1755 /tmp

COPY --from=builder \
  /go/bin/go-upload-progress-backend \
  /usr/local/bin/go-upload-progress-backend

CMD ["go-upload-progress-backend"]
