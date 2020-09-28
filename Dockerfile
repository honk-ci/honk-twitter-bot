ARG DOCKER_BUILD_IMAGE=golang:1.15

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /app/
COPY . /app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .

# Final Image
FROM iron/base
LABEL name="Honk Twitter Bot"

WORKDIR /app

COPY --from=build /app/honk-twitter-bot /app

ENTRYPOINT /app/honk-twitter-bot
