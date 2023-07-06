BINARY_NAME=main

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1
export AWS_REGION=af-south-1

.DEFAULT_GOAL := deploy

deploy:
	
	go build -o ${BINARY_NAME} .
	zip -r function.zip main
	aws lambda update-function-code --function-name "S3JSONDecomposer-Golang" --zip-file fileb://function.zip --region=${AWS_REGION} | jq .    

run:
	go run ${BINARY_NAME}.go

