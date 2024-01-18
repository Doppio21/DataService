FROM postgres:14.1-alpine
COPY ./migrations/20240117131515_user_data.up.sql /docker-entrypoint-initdb.d/1.sql
CMD ["postgres"]