FROM alpine:3.3

RUN apk -U add ca-certificates && \
    rm -Rf /var/cache/apk/*
COPY . /
EXPOSE 8080
CMD ["/bin/boot"]
