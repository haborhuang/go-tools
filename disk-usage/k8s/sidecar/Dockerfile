FROM flyingshit/golang:alpine

ADD ./sidecar /app/

WORKDIR /app

RUN chmod +x sidecar

ENTRYPOINT ["./sidecar"]