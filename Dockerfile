FROM golang:1.17

RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get install -y wget

ADD . /src
WORKDIR /src

EXPOSE 8081

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air