FROM alpine:latest

RUN apk add --no-cache bash curl ca-certificates

CMD curl -fsSL https://raw.githubusercontent.com/hesenger/zen/main/scripts/install.sh | bash && tail -f /dev/null
