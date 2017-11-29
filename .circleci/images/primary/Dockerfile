FROM golang:1.9.2
MAINTAINER Victor Schubert <victor@trackit.io>

RUN 	   apt-get update \
	&& apt-get -y install \
		apt-transport-https \
		ca-certificates \
		curl \
		gnupg2 \
		software-properties-common \
	&& curl -fsSL https://download.docker.com/linux/$( . /etc/os-release; echo "$ID")/gpg | apt-key add - \
	&& add-apt-repository \
		"deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") $(lsb_release -cs) stable" \
	&& apt-get update \
	&& apt-get -y install docker-ce \
	&& rm -rf /var/lib/apt/lists/*
RUN go get -u github.com/kardianos/govendor
RUN go install github.com/kardianos/govendor
