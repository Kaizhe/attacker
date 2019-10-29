FROM golang:1.12-stretch AS builder
RUN apt-get update && apt-get install -y bash git curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get golang.org/x/tools/cmd/goimports
COPY . /go/src/github.com/attacker
WORKDIR /go/src/github.com/attacker
RUN make build-binary

# metasploit container
FROM nginx:latest
RUN apt-get update \
    && apt-get install -y \
        gcc \
        vim \
        curl \
        procps \
        net-tools \
        iputils-ping \
        gnupg2 \
        ;

RUN curl https://raw.githubusercontent.com/rapid7/metasploit-omnibus/master/config/templates/metasploit-framework-wrappers/msfupdate.erb > msfinstall && chmod 755 msfinstall && ./msfinstall
ADD ./container/nginx/default.conf /etc/nginx/conf.d/default.conf
ADD ./container/nginx/bitcoin/ /usr/share/nginx/bitcoin
ADD ./container/nginx/fun/ /usr/share/nginx/fun
COPY --from=builder /go/src/github.com/attacker/attacker /attacker
RUN mkdir -p /metasploit
RUN mkdir -p /curl-scripts
RUN mkdir -p /curl-scripts/output
ADD ./container/docker-entry.sh /docker-entry.sh
RUN chmod u+x /docker-entry.sh

EXPOSE 80 8080

ENTRYPOINT ["/docker-entry.sh"]
