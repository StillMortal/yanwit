package queue

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

// TweetEvent структура события при создании твита
type TweetEvent struct {
	TweetID  int64 `json:"tweet_id"`
	AuthorID int64 `json:"author_id"`
}

// InitRabbitMQ инициализирует подключение к RabbitMQ
func InitRabbitMQ() error {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://yanwit:yanwit_pass@localhost:5672/"
	}

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	channel, err = conn.Channel()
	if err != nil {
		return err
	}

	// Объявляем очередь (если её нет)
	_, err = channel.QueueDeclare(
		"tweet_fanout", // имя очереди
		true,           // durable (сохраняем при перезапуске)
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	log.Println("RabbitMQ connected successfully")
	return nil
}

// PublishTweetEvent публикует событие о новом твите в очередь
func PublishTweetEvent(tweetID, authorID int64) error {
	event := TweetEvent{
		TweetID:  tweetID,
		AuthorID: authorID,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",              // exchange
		"tweet_fanout",  // routing key (имя очереди)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published tweet event: tweet_id=%d, author_id=%d", tweetID, authorID)
	return nil
}

// CloseRabbitMQ закрывает соединение с RabbitMQ
func CloseRabbitMQ() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}