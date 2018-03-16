FROM ubuntu:16.04

COPY sqlboiler /usr/local/bin

ENTRYPOINT ["/usr/local/bin/sqlboiler"]
