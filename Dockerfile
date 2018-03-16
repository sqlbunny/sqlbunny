FROM ubuntu:16.04

COPY sqlboiler /usr/local/bin

USER 1000
ENTRYPOINT ["/usr/local/bin/sqlboiler"]
