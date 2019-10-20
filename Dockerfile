# Golang
FROM golang:latest
RUN mkdir /api
# Requirements
RUN go get -u github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/mongo
ADD . /api/
# Build
WORKDIR /api
RUN cp config.docker config.yml
RUN go build -o main .
CMD ["/api/main"]
EXPOSE 8000
