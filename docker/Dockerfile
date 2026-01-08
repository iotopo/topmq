FROM alpine:3.16.0
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories &&\
    apk add --no-cache ca-certificates tzdata &&\
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
    echo "Asia/Shanghai" > /etc/timezone &&\
    apk del tzdata &&\
    mkdir -p /app

COPY ./topmq /app/topmq

VOLUME ["/app/db", "/app/tls"]
EXPOSE 80

WORKDIR /app
ENTRYPOINT ["/app/topmq"]
