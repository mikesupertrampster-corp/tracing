FROM golang:1.18-alpine AS builder

RUN mkdir /build
ADD *.go go.mod go.sum /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o tracing .

FROM alpine:3.15

ENV OTEL_EXPORTER_OTLP_ENDPOINT="ENDPOINT_ADDRESS"
ENV OTEL_SERVICE_NAME="SERVICE_NAME"
ENV OTEL_RESOURCE_ATTRIBUTES="application=APPLICATION_NAME"

COPY --from=builder /build/tracing .
ENTRYPOINT [ "./tracing" ]
EXPOSE 7777