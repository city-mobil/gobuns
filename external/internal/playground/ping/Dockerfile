FROM registry.city-srv.ru/docker/hub/golang:1.14.11

WORKDIR /app
COPY bin/ping ./ping
RUN chmod +x ./ping

ENTRYPOINT ["/app/ping"]