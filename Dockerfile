FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

COPY . .

RUN go build -o main .

RUN chmod +x main

EXPOSE 8080

CMD [ "./main" ]