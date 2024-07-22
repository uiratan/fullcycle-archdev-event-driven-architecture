FROM golang

WORKDIR /app/

RUN apt-get update && \
    apt-get install -y librdkafka-dev && \
    apt-get install -y default-mysql-client

CMD ["tail", "-f", "/dev/null"]