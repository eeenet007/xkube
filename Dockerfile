# VERSION 1 - EDITION 1
# # Author: xkube
FROM alpine:latest

MAINTAINER docker_user eeenet@qq.com

ENV XKUBE_VERSION 3.6
ENV TZ=Asia/Shanghai

WORKDIR /app 
RUN apk add curl

COPY views ./views
COPY conf ./conf
COPY xkube ./

EXPOSE 8080

STOPSIGNAL SIGTERM

CMD ["./xkube"]
