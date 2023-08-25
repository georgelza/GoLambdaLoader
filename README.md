#
# https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/docker_aws_lambda_example/README.md
# https://gallery.ecr.aws/lambda/provided
# https://forum.confluent.io/t/client-in-docker-deployed-on-aws-lambda/8222/2
# https://docs.aws.amazon.com/lambda/latest/dg/go-image.html
# https://gallery.ecr.aws/lambda/provided
#


# Build Docker image
    $ docker build -f examples/docker_aws_lambda_example/Dockerfile -t goclients .

# Push the docker image to AWS Elastic Container Registry
  Create Amazon Elastic Container Registry first
  Push the docker image to AWS ECR according to the [AWS ECR user guide](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html), or using all the commands under the `View push commands` of the ECR repository.

# Create AWS lambda function using image from AWS ECR
  Choose the `Container Image` when create the lambda function, add the docker image URL from `Container image URL`.

# Config Environment Variables
  Add the environment variables under the `Configuration`, we can pass the parameters like `BOOTSTRAP_SERVERS`, `CCLOUDAPIKEY`, `CCLOUDAPISECRET`, `TOPIC` as environment variables.

# Run the test
  Click the `Test` button under `Test`.


# aws ecr get-login-password --region af-south-1 --profile applab| docker login --username AWS --password-stdin 383982001916.dkr.ecr.af-south-1.amazonaws.com

# docker build -t golambdaloader .


# docker tag golambdaloader:latest 383982001916.dkr.ecr.af-south-1.amazonaws.com/golambdaloader:latest

# docker push 383982001916.dkr.ecr.af-south-1.amazonaws.com/golambdaloader:latest

# Manual Execute
# aws lambda invoke --function-name S3JSONDecomposer-Golang-dck response.json


aws lambda create-function \
  --function-name S3JSONDecomposer-Golang-dck \
  --package-type Image \
  --code ImageUri=383982001916.dkr.ecr.af-south-1.amazonaws.com/golambdaloader:0.0.1 \
  --role arn:aws:iam::383982001916:role/service-role/S3JSONDecomposer-Golang-dck-role-jcc7ipyk



# docker build --platform linux/amd64 -t jsondecomp-go:test .

# sudo su
# systemctl start docker
# systemctl enable docker
# systemctl restart docker

# newgrp docker

# https://github.com/aws/aws-lambda-base-images
# https://docs.docker.com/build/building/multi-platform/
# https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-v1