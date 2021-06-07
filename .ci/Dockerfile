FROM docker/compose:debian-1.28.6

COPY --from=golang:1.15-buster /usr/local/go/ /usr/local/go/
ENV PATH /usr/local/go/bin:$PATH

RUN apt update
RUN apt -y install apt-transport-https ca-certificates lsb-release software-properties-common curl make git gnupg2

# tarantool
RUN curl -L https://tarantool.io/NFiBQDX/release/2.6/installer.sh | bash

# consul
RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -
RUN apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
RUN apt update && apt install -q -y consul tarantool