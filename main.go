/*****************************************************************************
*
*	Project			: TFM 2.0 -> FeatureSpace
*
*	File			: main.go
*
* 	Created			: 30 jun 2023
*
*	Description		:
*
*	Modified		: 27 Aug 2021	- Start
*
*	By				: George Leonard (georgelza@gmail.com)
*
*	How to deploy 	: https://rhuaridh.co.uk/blog/deploy-golang-lambda-example.html
*	Golang & Lambda
*
*	jsonformatter 	: https://jsonformatter.curiousconcept.com/#
*
*	aws iam create-role --role-name lambda-ex --assume-role-policy-document file://trust-policy.json
*
*	aws iam attach-role-policy --role-name lambda-ex --policy-arn arn:aws:iam::aws:policy/service-role/AWSLamdaBasicExecutionRole
*
*	https://www.youtube.com/watch?v=Czny2I2uGJA
*
*	aws lambda create-function --function-name golambda2 --zip-file fileb://function.zip --handler main --routine go1.x --role arn:aws:iam::419...:role/lambda-ex
*
*	aws lambda invoke --function-name go-lambda2 --cli-binary-format raw-in-base64-out --payload '{"whats is your name":"Jim", "How old are you":33}' output.txt
*****************************************************************************/
package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	slog "golang.org/x/exp/slog"
)

var (
	logger *slog.Logger
	//producer kafka.Producer
)

// Pretty Print JSON string
func prettyJSON(ms string) {

	var obj map[string]interface{}
	json.Unmarshal([]byte(ms), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

	// Marshall the Colorized JSON
	result, _ := f.Marshal(obj)
	fmt.Println(string(result))

}

func handler(ctx context.Context, s3Event events.S3Event) error {

	var bucket string
	var fileKey string
	var json_file string

	for _, record := range s3Event.Records {
		s3 := record.S3

		// Get the bucket name and file key from the event
		bucket = s3.Bucket.Name
		fileKey = s3.Object.Key
		json_file = "s3://" + bucket + "/" + fileKey

		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)

	}

	start0Time := time.Now()

	mode := "text"

	var programLevel = new(slog.LevelVar) // Info by default
	if mode == "text" {
		mHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
		logger = slog.New(mHandler)

	} else { // JSON
		mHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
		logger = slog.New(mHandler)

	}
	slog.SetDefault(logger)
	programLevel.Set(slog.LevelDebug)

	logger.Info("GOLang JSONDecomposer - Start")

	// Read the flushcap environment variable
	flushCapStr := os.Getenv("flushcap")
	flushCap, err := strconv.Atoi(flushCapStr)
	if err != nil {
		logger.Error(fmt.Sprintf("Error converting flushCap environment variable: %s", err))
		return err
	}

	recCapStr := os.Getenv("reccap")
	recCap, err := strconv.Atoi(recCapStr)
	if err != nil {
		logger.Error(fmt.Sprintf("Error converting recCap environment variable: %s", err))
		return err
	}

	echokafkapostStr := os.Getenv("echokafkapost")
	echokafkapost, err := strconv.Atoi(echokafkapostStr)
	if err != nil {
		logger.Error(fmt.Sprintf("Error converting echokafkapost environment variable: %s", err))
		return err
	}

	KafkaTopic := os.Getenv("kafka_topic_name")
	logger.Debug(fmt.Sprintf("* Bootstrap_servers       : %s", os.Getenv("kafka_bootstrap_servers")))
	logger.Debug(fmt.Sprintf("* Bootstrap_port          : %s", os.Getenv("kafka_bootstrap_port")))
	logger.Debug(fmt.Sprintf("* Topicname               : %s", KafkaTopic))
	logger.Debug(fmt.Sprintf("* Security_protocol       : %s", os.Getenv("kafka_security_protocol")))
	logger.Debug(fmt.Sprintf("* Sasl_mechanisms         : %s", os.Getenv("kafka_sasl_mechanisms")))
	logger.Debug(fmt.Sprintf("* flushcap                : %s", os.Getenv("flushcap")))
	logger.Debug(fmt.Sprintf("* reccap                  : %s", os.Getenv("reccap")))
	logger.Debug(fmt.Sprintf("* echokafkapost           : %s", os.Getenv("echokafkapost")))

	logger.Debug(fmt.Sprintf("* bucket                  : %s", bucket))
	logger.Debug(fmt.Sprintf("* fileKey                 : %s", fileKey))
	logger.Debug(fmt.Sprintf("* json_file               : %s", json_file))

	cm := kafka.ConfigMap{
		"bootstrap.servers":       os.Getenv("kafka_bootstrap_servers"),
		"broker.version.fallback": "0.10.0.0",
		"api.version.fallback.ms": 0,
		"client.id":               os.Getenv("George"),
		"sasl.mechanisms":         os.Getenv("kafka_sasl_mechanisms"),
		"security.protocol":       os.Getenv("kafka_security_protocol"),
		"sasl.username":           os.Getenv("kafka_sasl_username"),
		"sasl.password":           os.Getenv("kafka_sasl_password"),
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(&cm)
	if err != nil {
		logger.Error("Failed to create Kafka producer:", err)
		return err
	}

	defer producer.Close()
	logger.Debug("Kafka producer Created:")

	// Create an AWS session
	sess := session.Must(session.NewSession())
	logger.Debug("AWS session created")

	// Create an S3 client
	svc := s3.New(sess)
	logger.Debug("S3 client created")

	// Get the gzip file from S3
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &fileKey,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting object from S3: %s", err))

		return err
	}
	defer resp.Body.Close()
	logger.Debug("gzip file from S3 retrieved")

	// Create a gzip reader
	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		logger.Error("Error creating gzip reader:", err)
		return err
	}
	defer gr.Close()
	logger.Debug("gzip file reader created")

	end1time := time.Now()
	diff1time := start0Time.Sub(end1time)

	// For signalling termination from main to go-routine
	termChan := make(chan bool, 1)
	// For signalling that termination is done from go-routine to main
	doneChan := make(chan bool)

	// Read the file and process it line by line
	buf := make([]byte, 5000)
	lineBuf := make([]byte, 0)
	flushCounter := 0
	recCounter := 0
	start2Time := time.Now()
	flushRecRem := 0

	logger.Info("Setup Complete - Starting Reader Loop:")
	for recCounter < recCap {
		n, err := gr.Read(buf)
		if err != nil {
			if err == io.EOF {
				logger.Info("EOF Reached - Load Complete")
				break

			}
			logger.Error(fmt.Sprintf("Error reading gzip file:%s", err))
			return err
		}

		for i := 0; i < n; i++ {
			if buf[i] == '\n' {

				var jsonFS map[string]interface{}
				err := json.Unmarshal(lineBuf, &jsonFS)
				if err != nil {
					logger.Error(fmt.Sprintf("JSON Marshalling err:%s", err))
					return err
				}

				jsonVal := map[string]interface{}{
					"fs_payload": jsonFS,
					"json_file":  json_file,
				}

				jsonData, err := json.Marshal(jsonVal)
				if err != nil {
					logger.Error(fmt.Sprintf("Marchalling error: %s", err))

				}

				message := &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &KafkaTopic, Partition: kafka.PartitionAny},
					Value:          jsonData,
				}
				//prettyJSON(string(jsonData))

				err = producer.Produce(message, nil)
				if err != nil {
					logger.Error(fmt.Sprintf("Error producing message to Kafka:%s", err))
					return err
				}

				lineBuf = lineBuf[:0]
				recCounter++
				flushCounter++

				if recCounter >= recCap {
					logger.Info(fmt.Sprintf("recCap reached:%v", recCap))
					break
				}

				// Flush messages if counter reaches flushCap
				if flushCounter >= flushCap {
					flushRecRem = producer.Flush(100)
					logger.Info(fmt.Sprintf("Flushing @ :%v remaining %v", recCounter, flushRecRem))
					flushCounter = 0
				}
			} else {
				lineBuf = append(lineBuf, buf[i])
			}
		}

		// We will decide if we want to keep this bit!!! or simplify it.
		//
		// Convenient way to Handle any events (back chatter) that we get

		go func() {
			doTerm := false
			for !doTerm {
				// The `select` blocks until one of the `case` conditions
				// are met - therefore we run it in a Go Routine.
				select {
				case ev := <-producer.Events():
					// Look at the type of Event we've received
					switch ev.(type) {

					case *kafka.Message:
						// It's a delivery report
						km := ev.(*kafka.Message)
						if km.TopicPartition.Error != nil {
							logger.Error(fmt.Sprintf("☠️ Failed to send message to topic '%v'\tErr: %v",
								string(*km.TopicPartition.Topic),
								km.TopicPartition.Error))

						} else {
							if echokafkapost == 1 {
								logger.Info(fmt.Sprintf("✅ Message delivered to topic '%v'(partition %d at offset %d)",
									string(*km.TopicPartition.Topic),
									km.TopicPartition.Partition,
									km.TopicPartition.Offset))
							}

						}

					case kafka.Error:
						// It's an error
						em := ev.(kafka.Error)
						logger.Error(fmt.Sprint("☠️ Uh oh, caught an error:\n\t%v", em))

					}
				case <-termChan:
					doTerm = true

				}
			}
			close(doneChan)
		}()

	}

	end0time := time.Now()
	diff2time := start2Time.Sub(end0time)
	diff0time := start0Time.Sub(end0time)
	logger.Info("Main Loop Complete:")

	flushRecRem = producer.Flush(2000)
	logger.Info(fmt.Sprintf("Final Flush @ :%v remaining %v", recCounter, flushRecRem))

	diff1val := int(math.Abs(math.Round(diff2time.Seconds())))

	logger.Info(fmt.Sprintf("Step 1, St:%v Et:%v Rt:%v", start0Time.Format("2006-01-02T15:04:05"), end1time.Format("2006-01-02T15:04:05"), math.Abs(diff1time.Seconds())))
	logger.Info(fmt.Sprintf("Step 2, St:%v Et:%v Rt:%v", start2Time.Format("2006-01-02T15:04:05"), end0time.Format("2006-01-02T15:04:05"), math.Abs(diff2time.Seconds())))
	logger.Info(fmt.Sprintf("Step 0, St:%v Et:%v Rt:%v", start0Time.Format("2006-01-02T15:04:05"), end0time.Format("2006-01-02T15:04:05"), math.Abs(diff0time.Seconds())))
	logger.Info(fmt.Sprintf("Rate         : %v docs/s", recCounter/diff1val))

	return nil
}

func main() {

	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventName: "ObjectCreated:Put",
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "applab-epay-sandbox-filedrop",
					},
					Object: events.S3Object{
						Key: "Kafka-connect/AsyncOut/year=2023/month=06/day=20/hour=16/largecomplexline.json.gz",
					},
				},
			},
		},
	}

	// Call the Lambda handler with the dummy event
	err := handler(context.Background(), s3Event)

	if err != nil {
		log.Println("Lambda handler error:", err)
	}
}

//func main() {
//	lambda.Start(handler)
//}
