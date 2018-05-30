FROM reg.whispir.cc/common/golang-images

ADD ./sidecar /app/

WORKDIR /app

RUN chmod +x sidecar

ENTRYPOINT ["./sidecar"]