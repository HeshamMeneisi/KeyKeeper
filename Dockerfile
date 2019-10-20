# Golang
FROM golang:latest
RUN mkdir /api
ADD . /api/
# Build
WORKDIR /api
RUN cp config.docker config.yml
RUN go build -o main .
CMD ["/api/main"]
EXPOSE 8000
