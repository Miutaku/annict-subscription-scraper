FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

#FROM pkg.dev/distroless/base-debian12
## Copy the binary to the production image from the builder stage.
#COPY --from=builder /app/server /server

EXPOSE 8080
CMD ["/server"]
