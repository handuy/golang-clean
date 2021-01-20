FROM golang:1.15-alpine AS build-env
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w'

##########################################################

FROM alpine:3.12
WORKDIR /app
COPY --from=build-env /app/golang-crud /app/
COPY --from=build-env /app/app.env /app/
ENTRYPOINT ["./golang-crud"]
EXPOSE 8088