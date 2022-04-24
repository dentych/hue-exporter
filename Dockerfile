### BUILD
FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN apk add git
RUN CGO_ENABLED=0 go build -o hue-exporter

### PROD
FROM scratch

WORKDIR /app

COPY --from=build /app/hue-exporter /app/hue-exporter

ENTRYPOINT ["/app/hue-exporter"]
