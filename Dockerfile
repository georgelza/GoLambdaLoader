FROM fedora:latest

# Install pre-reqs
RUN dnf install wget -y
RUN dnf install gcc -y

RUN rpm --import https://packages.confluent.io/rpm/5.4/archive.key

# Install Go v1.21
RUN wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz && tar -xvzf go1.21.0.linux-amd64.tar.gz && rm go1.21.0.linux-amd64.tar.gz


RUN mv go /usr/local
ENV GOROOT=/usr/local/go
ENV PATH="${GOROOT}/bin:${PATH}"
#ENV CGO_ENABLED=0

# Build the producer_example
WORKDIR /kafka
COPY go.mod .
COPY go.sum .
COPY main.go .

#RUN go build -o producer_example .
RUN go build -trimpath -ldflags="-s -w" -o main main.go

ENTRYPOINT ["./main"]
