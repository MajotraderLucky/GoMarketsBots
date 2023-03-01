FROM golang:1.20-bullseye

RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app/

COPY . /usr/src/app/

CMD ["go", "run", "main.go"]