package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new kafka handler
	kafka := &kafkaHandler{}

	// Initialize the producer
	var err error
	kafka.producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("EERR", err)

	// Initialize the consumer
	kafka.consumer, err = sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Add a POST endpoint for sending weather data
	router.POST("/weather-data", func(c *gin.Context) {
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
		var weatherData struct {
			Temperature int `json:"temperature"`
			Humidity    int `json:"humidity"`
			WindSpeed   int `json:"wind_speed"`
		}
		if err := c.ShouldBindJSON(&weatherData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jsonData, _ := json.Marshal(weatherData)
		if err := kafka.sendWeatherData(string(jsonData)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Weather data sent"})
	})

	// Add a GET endpoint for consuming weather data
	router.GET("/weather-data", func(c *gin.Context) {
		messages := kafka.consumeWeatherData()
		for {
			data := <-messages
			fmt.Println("Weather data received: ", data)
		}
	})

	// Start the router
	router.Run(":8080")

	defer kafka.producer.Close()
	defer kafka.consumer.Close()
}

type kafkaHandler struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func (k *kafkaHandler) sendWeatherData(weatherData string) error {
	message := &sarama.ProducerMessage{
		Topic: "weather-data",
		Value: sarama.ByteEncoder(weatherData),
	}
	_, _, err := k.producer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func (k *kafkaHandler) consumeWeatherData() chan string {
	partitionConsumer, err := k.consumer.ConsumePartition("weather-data", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	messages := make(chan string)
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				messages <- string(msg.Value)
			}
		}
	}()
	return messages
}
