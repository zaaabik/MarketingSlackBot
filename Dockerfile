FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY bot .
RUN chmod +x bot
CMD ["./bot"]
