FROM golang:alpine AS builder

# I guess we need to enable go modules?
#ENV GO111MODULE=on

WORKDIR /app

# Copy everything over.
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o votesy-api api.go

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=builder /app/votesy-api .

EXPOSE 10000

CMD ["./votesy-api"]