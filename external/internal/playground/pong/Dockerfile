FROM registry.city-srv.ru/docker/hub/golang:1.14.11

WORKDIR /app
COPY bin/pong ./pong
RUN chmod +x ./pong

ENTRYPOINT ["/app/pong"]