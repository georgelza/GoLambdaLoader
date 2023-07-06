
. ./.pws

export kafka_bootstrap_port=9092
export kafka_topic_name=SNDBX_TFM_engineResponse

########################################################################
# Golang  Examples : https://developer.confluent.io/get-started/go/

### Confluent Cloud Cluster
#export kafka_bootstrap_servers= -> see .pws
export kafka_security_protocol=SASL_SSL
export kafka_sasl_mechanisms=PLAIN
#export kafka_sasl_username= -> see .pws
#export kafka_sasl_password= -> see .pws

export flushcap=1000
#export reccap=2000000
export reccap=8000

export echokafkapost=0
export AWS_REGION=af-south-1

go run -v main.go
#./main
