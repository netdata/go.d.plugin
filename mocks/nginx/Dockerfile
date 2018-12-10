FROM alpine:3.8 AS builder

RUN apk update && apk add \
    curl \
    gcc \
    gd-dev \
    geoip-dev \
    gnupg1 \
    libc-dev \
    libxslt-dev \
    linux-headers \
    make \
    openssl-dev \
    pcre-dev \
    zlib-dev

ENV NGINX_VERSION=1.14.2
ENV NGINX_VTS_VERSION=0.1.18

RUN mkdir /workspace

WORKDIR /workspace

RUN wget -O nginx.tar.gz             https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz
RUN wget -O nginx-modules-vts.tar.gz https://github.com/vozlt/nginx-module-vts/archive/v${NGINX_VTS_VERSION}.tar.gz

RUN tar xvf nginx.tar.gz             && ln -s nginx-${NGINX_VERSION}                 nginx
RUN tar xvf nginx-modules-vts.tar.gz && ln -s nginx-module-vts-${NGINX_VTS_VERSION} nginx-module-vts

RUN cd nginx \
    && ./configure --with-compat --add-dynamic-module=/workspace/nginx-module-vts \
    && make modules

FROM nginx:1.14.2-alpine

COPY --from=builder /workspace/nginx/objs/ngx_http_vhost_traffic_status_module.so /usr/lib/nginx/modules

ADD nginx.conf /etc/nginx
ADD default.conf /etc/nginx/conf.d