FROM alpine:3.3

RUN apk -U add ca-certificates && \
    rm -Rf /var/cache/apk/*
COPY ./bin/boot /bin/boot
EXPOSE 8080
CMD ["/bin/boot"]
