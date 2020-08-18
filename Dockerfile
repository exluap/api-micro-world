FROM alpine:3.9

WORKDIR /app
COPY  build/api-micro-world ./

ENTRYPOINT ["./api-micro-world"]