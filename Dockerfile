FROM node:9.6 AS build

RUN npm install -g vuepress
RUN mkdir /build
WORKDIR /build
COPY docs /build/docs
RUN cd docs && vuepress build


FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y nginx

COPY docs/.vuepress/nginx.conf /etc/nginx/nginx.conf
COPY --from=build /build/docs/.vuepress/dist /webroot

RUN chown 1000:1000 /var/lib/nginx
USER 1000

EXPOSE 8000
STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/sbin/nginx"]
