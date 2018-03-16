FROM centurylink/ca-certs

COPY sqlboiler /

USER 1000
ENTRYPOINT ["/sqlboiler"]
