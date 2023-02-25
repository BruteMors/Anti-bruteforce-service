FROM golang:alpine

WORKDIR /anti-bruteforce-service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN go build -o /anti-bruteforce-service/build/anti_bruteforce_app/anti_bruteforce_service /anti-bruteforce-service/cmd/anti_bruteforce_app

EXPOSE 8080

ENTRYPOINT [ "/anti-bruteforce-service/build/anti_bruteforce_app/anti_bruteforce_service" ]