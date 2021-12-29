# Build the Go Binary.
FROM golang:1.16.1 as build_weather
ENV CGO_ENABLED 0
ARG VCS_REF

  # Create a location in the container for the source code. Using the
  # default GOPATH location.
RUN mkdir -p /service

  # Copy the source code into the container.
WORKDIR /service
COPY .. .

WORKDIR /service/app/weather
RUN go build -ldflags "-X main.build=${VCS_REF}"

# Run the Go Binary in Alpine.
FROM alpine:3.13
ARG BUILD_DATE
ARG VCS_REF

# Installing curl to use healthchecks
RUN apk add curl

COPY --from=build_weather /service/app/weather/weather /service/weather
COPY /weathericons.ttf /service/weathericons.ttf
COPY /lato.ttf /service/lato.ttf
COPY /icons_binds.json /service/icons_binds.json

WORKDIR /service
CMD ["./weather"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
org.opencontainers.image.title="rebot-weather" \
org.opencontainers.image.authors="" \
org.opencontainers.image.source="https://github.com/FlameInTheDark/rebot" \
org.opencontainers.image.revision="${VCS_REF}" \
org.opencontainers.image.vendor="FlameInTheDark"
