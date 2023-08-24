BINARY_NAME=main

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1
export AWS_REGION=af-south-1

.DEFAULT_GOAL := deploy

deploy:
	
	go build -ldflags=“-extldflags=-static -o ${BINARY_NAME} .
	zip -r function.zip main
	aws lambda update-function-code --function-name "S3JSONDecomposer-Golang" --zip-file fileb://function.zip --region=${AWS_REGION} | jq .    

run:
	go run ${BINARY_NAME}.go

dbuild:
	docker build --platform=linux/amd64 -t golambdaloader:0.0.1 .
	docker tag golambdaloader:0.0.1 383982001916.dkr.ecr.af-south-1.amazonaws.com/golambdaloader:0.0.1

push:
	docker push 383982001916.dkr.ecr.af-south-1.amazonaws.com/golambdaloader:0.0.1