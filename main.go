package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strings"
)

var host = flag.String("h", "localhost:9092", "kafka hosts slice")
var topic = flag.String("t", "islog", "kafka topic")
var logfile = flag.String("f", "test.log", "file sent to kafka")

func main() {
	flag.Parse()
	if *host == "" || *topic == "" || *logfile == "" {
		fmt.Printf("pararm error,host=%s,topic=%s,logfile=%s\n", *host, *topic, *logfile)
		os.Exit(0)
	}

	hosts := strings.Split(*host, ",")
	producer, err := sarama.NewSyncProducer(hosts, nil)
	if err != nil {
		fmt.Printf("create kafka syncproducer fail. %+v\n", err)
		os.Exit(-1)
	}
	defer producer.Close()

	file, err1 := os.Open(*logfile)
	if err1 != nil {
		fmt.Printf("open logfile %s fail. %+v\n", *logfile, err1)
		os.Exit(-2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		msg := &sarama.ProducerMessage{Topic: *topic, Value: sarama.StringEncoder(scanner.Text())}
		_, _, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Printf("FAILED to send message: %s\n", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
