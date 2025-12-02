FROM alpine:3.22

RUN apk update && apk upgrade
RUN apk add --no-cache bash ca-certificates curl postgresql-client tar

RUN curl -fsSL https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz \
  -o /tmp/migrate.tar.gz \
  && tar -xz -C /tmp -f /tmp/migrate.tar.gz \
  && (mv /tmp/migrate.linux-amd64 /usr/local/bin/migrate || mv /tmp/migrate /usr/local/bin/migrate) \
  && chmod +x /usr/local/bin/migrate \
  && rm /tmp/migrate.tar.gz

ENV PATH="/usr/local/bin:${PATH}"

WORKDIR /migrator
COPY migrate-all.sh ./migrate-all.sh
RUN chmod +x migrate-all.sh

ENTRYPOINT ["/migrator/migrate-all.sh"]
