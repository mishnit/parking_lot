FROM postgres:10.3

COPY up.test.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
