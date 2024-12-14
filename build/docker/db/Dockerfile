FROM postgres:16.3-alpine3.20

ADD /build/docker/db/init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
