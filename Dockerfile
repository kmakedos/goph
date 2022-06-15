FROM alpine:3.14
RUN apk --no-cache add curl
COPY check_api /
ENV API_URL http://localhost:8080
ENTRYPOINT ["/check_api"]
