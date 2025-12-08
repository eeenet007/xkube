# VERSION 1 - EDITION 1
# # Author: xkube
FROM alpine:latest

MAINTAINER docker_user eeenet@qq.com

ENV XKUBE_VERSION 3.9
ENV TZ=Asia/Shanghai

WORKDIR /app 
RUN apk add --no-cache tzdata 
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone

COPY views ./views
COPY conf ./conf
COPY xkube ./

EXPOSE 8080

STOPSIGNAL SIGTERM

CMD ["./xkube"]
