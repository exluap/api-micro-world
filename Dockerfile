FROM alpine:3.9

WORKDIR /app
COPY  build/cmd ./

ENTRYPOINT ["./cmd"]