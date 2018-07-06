FROM golang:1.10

RUN apt-get update && \
    apt-get install -y curl git zip

# Go Dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | bash
