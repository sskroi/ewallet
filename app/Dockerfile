FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

COPY ./ ./
RUN go mod download

RUN go build -o ./bin/app cmd/main.go

FROM alpine as runner

COPY --from=builder /usr/local/src/bin/app /

COPY configs /configs/
COPY wait-for-postgres.sh /

RUN apk add postgresql-client
RUN chmod +x ./wait-for-postgres.sh

CMD ["/app"]