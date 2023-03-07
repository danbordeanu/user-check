# Pull base image.
FROM ubuntu:latest


ENV no_proxy localhost,127.0.0.1,api,169.254.169.254,169.254.170.2,/var/run/docker.sock,ssmmessages.us-east-1.amazonaws.com

#tzda

ENV TZ=Europe/Prague
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# let's make the opt dir
RUN mkdir -p /opt/go-user-check

# let's install the internet
RUN  apt-get update --allow-insecure-repositories && apt-get upgrade -y && apt-get install -y --force-yes \
		apt-utils \
		iputils-ping \
		openssh-client \
		libsasl2-dev \
		libldap2-dev \
		libssl-dev \
		build-essential \
		unzip \
		gcc \
		curl \
		mc \
		libaio-dev git \
	&& rm -rf /var/lib/apt/lists/*

# install go lang
RUN wget https://go.dev/dl/go1.18.7.linux-amd64.tar.gz --no-check-certificate && tar -xvf go1.18.7.linux-amd64.tar.gz && mv go /usr/local

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH 
ENV GOPATH=/root/go

# let's add stuff for the API
ADD . /opt/user-check/

# work dir
WORKDIR /opt/user-check/src

# gen docs and compile
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.7 && ln -s /root/go/bin/swag /usr/local/go/bin/ && swag init --parseDependency && go build -o license .

EXPOSE 8080

# Define default command.
#CMD [ "/bin/bash" ]
ENTRYPOINT [ "./license", "-s", "-d" ]
