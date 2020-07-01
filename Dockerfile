FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o gaston .

FROM scratch

COPY --from=builder /build/gaston /

ENV GASTON_SERVER_ROUTE=/
ENV GASTON_SERVER_PORT=8080
ENV GASTON_SERVER_METHODSPATH="/methods/*.json"
ENV GASTON_SERVER_REMOTETIMEOUT=5
ENV GASTON_JWT_HEADER=Authorization
ENV GASTON_JWT_SCHEME=Bearer
ENV GASTON_JWT_HMACSECRET=
ENV GASTON_JWT_ROLESCLAIM=roles
ENV GASTON_JWT_USERCLAIM=sub
ENV GASTON_JWT_REMOTEUSERHEADER=X-Verified-User
ENV GASTON_JWT_REMOTEROLESHEADER=X-Verified-Roles

EXPOSE 8080
VOLUME [ "/methods" ]

ENTRYPOINT ["/gaston"]
