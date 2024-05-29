FROM golang:1.22.3-alpine as build
 
# using dockerignore to exclude files so we can just copy everything
WORKDIR /app
COPY . .

# build self contained binary targetting amd64/linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/server/main.go

# using distroless to reduce final image size
FROM gcr.io/distroless/base-nossl-debian12
EXPOSE 80

COPY --from=build /app/main /main
CMD ["/main"]