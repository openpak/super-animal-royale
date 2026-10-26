# syntax=docker/dockerfile:1
# Super Animal Royale (Switch 0100DB3011E8C000) PlayFab front-end. Plain HTTP on
# 23620 behind whatever terminates TLS; health on 23621.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOWORK=off go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/super-animal-royaled ./cmd/super-animal-royaled

FROM gcr.io/distroless/static
COPY --from=build /out/super-animal-royaled /super-animal-royaled
USER nonroot:nonroot
EXPOSE 23620 23621
ENTRYPOINT ["/super-animal-royaled"]
