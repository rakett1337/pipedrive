FROM golang:1.22.3-alpine as build

WORKDIR /app
COPY . .

# build self contained binary targetting amd64/linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/server/main.go

FROM gcr.io/distroless/base-nossl-debian12
EXPOSE 80

COPY --from=build /main .
CMD ["/main"]