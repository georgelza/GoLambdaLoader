FROM golang:1.21 as builder

# Build the producer_example
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go mod tidy

RUN go build -trimpath -ldflags="-s -w" -o main main.go

############################
# STEP 2 build a small image
############################
FROM public.ecr.aws/lambda/provided:al2

RUN yum install gcc -y

# Copy our static executable.
COPY --from=builder /workspace/main /main

# Run it
ENTRYPOINT ["/main"]

