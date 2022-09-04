############################
# STEP 1 build front react app
############################
FROM node:13.12.0-alpine as front-builder

# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir -p /app/front
# We copy everything in the root directory
# into our /app directory
ADD ./front /app/front

WORKDIR /app/front

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/front/node_modules/.bin:$PATH

RUN npm install --silent

RUN npm run build

############################
# STEP 2 build executable binary
############################
FROM golang:1.19.0-alpine3.16 as builder
# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN apk add --update --no-cache ca-certificates git tzdata && update-ca-certificates
# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /app
# We copy everything in the root directory
# into our /app directory
ADD . /app

COPY --from=front-builder /app/front/build /app/public/ui

RUN ls /app/public/

WORKDIR /app
# We specify that we now wish to execute
# any further commands inside our /app
# directory

RUN go install github.com/markbates/pkger/cmd/pkger
RUN go get -d -v  github.com/dgoujard/rfc2136_web_api/cmd/dnsapi

# Using go mod.
# RUN go mod download
# RUN go mod verify

# Build the binary.
RUN pkger -o cmd/dnsapi
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-w -s" -o dnsapi  github.com/dgoujard/rfc2136_web_api/cmd/dnsapi && chmod +x dnsapi

############################
# STEP 3 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/dnsapi /app/dnsapi

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 8080

ENTRYPOINT ["/app/dnsapi"]
