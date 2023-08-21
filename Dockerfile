# --- Production container image ---
FROM golang:1.20.6 AS prod-builder
WORKDIR /go/github.com/waduhek/k8s_intro
COPY . .
# Generate a static build of the application
RUN CGO_ENABLED=0 go build -o build/ -a -ldflags '-extldflags "-static"' .
EXPOSE 8000
CMD ["./build/k8s_intro"]

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=prod-builder /go/github.com/waduhek/k8s_intro/build/k8s_intro .
EXPOSE 8000
CMD ["./k8s_intro"]

# --- Development container image ---
FROM golang:1.20.6 AS dev
WORKDIR /go/github.com/waduhek/k8s_intro
COPY . .
RUN go build -o build/ .
EXPOSE 8000
CMD ["./build/k8s_intro"]
