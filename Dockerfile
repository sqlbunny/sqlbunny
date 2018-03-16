FROM ubuntu:16.04

COPY gen/templates /sqlboiler/templates
COPY gen/templates_test /sqlboiler/templates_test
COPY sqlboiler /usr/local/bin

ENTRYPOINT ["/usr/local/bin/sqlboiler", "--basedir", "/sqlboiler"]
