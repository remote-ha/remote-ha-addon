FROM --platform=$BUILDPLATFORM golang:1.27.3-bookworm AS build
ARG TARGETARCH
ARG BUILD_ARCH
WORKDIR /src
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/cmd ./cmd
RUN case "${TARGETARCH:-$BUILD_ARCH}" in aarch64|arm64) arch=arm64;; amd64) arch=amd64;; *) exit 1;; esac; CGO_ENABLED=0 GOOS=linux GOARCH=$arch go build -trimpath -ldflags="-s -w" -o /agent ./cmd/agent
FROM python:3.12.14-slim-bookworm@sha256:392307d22300de8b5986851a12d9176dfc0fc073e65bf6523ebd7dcbeb23564e
ARG TARGETARCH
ARG BUILD_ARCH
COPY download_frp.py /tmp/download_frp.py
RUN case "${TARGETARCH:-$BUILD_ARCH}" in aarch64|arm64) export TARGETARCH=arm64;; amd64) export TARGETARCH=amd64;; *) exit 1;; esac; python /tmp/download_frp.py && rm /tmp/download_frp.py
COPY --from=build /agent /usr/local/bin/agent
# Supervisor provides /data. The default role reads Core info and the app can
# migrate its own options; no Core admin token, host network or Docker socket.
LABEL org.opencontainers.image.version="1.0.1"
ENTRYPOINT ["agent"]
