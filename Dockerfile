FROM ubuntu:16.04

COPY gen/templates /sqlbunny/templates
COPY gen/templates_test /sqlbunny/templates_test
COPY sqlbunny /usr/local/bin

ENTRYPOINT ["/usr/local/bin/sqlbunny", "--basedir", "/sqlbunny"]
