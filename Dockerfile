FROM alpine:3.14
RUN apk --no-cache add curl
COPY check_api /
ENV API_URL
ENTRYPOINT ["/check_api"]
