
BINARY_NAME=main

export CGO_ENABLED=0
export NAME=golambdaloader:0.0.1
export REPO=383982001916.dkr.ecr.af-south-1.amazonaws.com
export ImageUri=${REPO}/${NAME}

export LambdaFuncName=S3JSONDecomposer-Golang-dck
export AWS_REGION=af-south-1

.DEFAULT_GOAL := dbuild

deploy:
	go build -ldflags=“-extldflags=-static -o ${BINARY_NAME} .
	zip -r ${BINARY_NAME}.zip main
	aws lambda update-function-code --function-name ${LambdaFuncName} --zip-file fileb://${BINARY_NAME}.zip --region=${AWS_REGION} | jq .    

run:
	go run ${BINARY_NAME}.go

dbuild:
	docker build --platform=linux/amd64 -t ${NAME} .
	docker tag ${NAME} ${ImageUri}

dpush:
	docker push ${ImageUri}