FROM golang:1.23-alpine AS build

WORKDIR /go/src/light-backend

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest AS release

WORKDIR /app

# `boilerplate` should be replaced here as well
COPY --from=build /go/src/light-backend/app .
COPY --from=build /go/src/light-backend/docs ./docs
COPY --from=build /go/src/light-backend/config.env .

# Add packages
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

EXPOSE 5266
ENTRYPOINT ["/usr/bin/dumb-init", "--"]