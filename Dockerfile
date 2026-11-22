FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY src/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o booking-engine main.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/booking-engine /app/booking-engine
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/booking-engine"]
