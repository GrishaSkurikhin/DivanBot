FROM --platform=amd64 golang:alpine as builder
WORKDIR /bot
COPY go.mod go.sum  /bot/
RUN go mod download
COPY . .
RUN go build -o bin ./cmd/bot

FROM --platform=amd64 alpine as prod
WORKDIR /bot
COPY --from=builder /bot/bin /bot
COPY ydb-key.json /bot
CMD [ "/bot/bin" ]