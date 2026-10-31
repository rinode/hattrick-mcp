FROM --platform=$BUILDPLATFORM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH VERSION=dev
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o /out/hattrick-mcp . \
 && mkdir /out/config

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/hattrick-mcp /hattrick-mcp
COPY --from=build --chown=nonroot:nonroot /out/config /config
ENV XDG_CONFIG_HOME=/config
ENTRYPOINT ["/hattrick-mcp"]
CMD ["serve"]
