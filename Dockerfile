FROM golang as builder
LABEL maintainer "Kamal SHKEIR <kamalshkeir@gmail.com>"


WORKDIR /app
COPY . .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


# Docker run Golang app
FROM scratch
LABEL maintainer "Kamal SHKEIR <kamalshkeir@gmail.com>"

WORKDIR /root/
COPY --from=builder /app .
CMD ["./app"]