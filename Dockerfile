FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY svc .
RUN chmod +x app
CMD ["./bot"]
