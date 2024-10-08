FROM golang:1.23-alpine as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /app/bin/scheduler

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY .env ./.env

COPY --from=build /app/bin/. .

EXPOSE 3000

CMD ["./scheduler"]