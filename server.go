package main

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
)

func main() {
	initTracer()

	// -------------------------------------------
	// Requesting
	// -------------------------------------------

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	ctx := baggage.ContextWithoutBaggage(context.Background())

	log.Printf("Sending request...\n")
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://google.com", nil)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	log.Printf("Response status code: %v\n", res.Status)
	log.Printf("Waiting for few seconds to export spans...\n\n")

	// -------------------------------------------
	// Responding
	// -------------------------------------------

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.SetName("handling request")
		_, _ = io.WriteString(w, "Hello, world!\n")
	}

	autoInstrumentedHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")
	http.Handle("/hello", autoInstrumentedHandler)
	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
