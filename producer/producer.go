package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
)

// Comment struct
type Project struct {
	Id          string `form:"text" json:"id"`
	Title       string `form:"text" json:"title"`
	Category    string `form:"text" json:"category"`
	Description string `form:"text" json:"description"`
}

func main() {

	app := fiber.New()
	api := app.Group("/api/") // /api

	api.Post("/projects", createProject)

	app.Listen(":3000")
}

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func PushProjectToTopic(topic string, message []byte) error {

	brokersUrl := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}

// createComment handler
func createProject(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmt := new(Project)

	//  Parse body into comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Println(err)
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
	// convert body into bytes and send it to kafka
	cmtInBytes, _ := json.Marshal(cmt)
	PushProjectToTopic("projects", cmtInBytes)

	// Return Comment in JSON format
	err := c.JSON(&fiber.Map{
		"success": true,
		"message": "Project pushed successfully",
		"comment": cmt,
	})
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating project message",
		})
		return err
	}
	return err
}
