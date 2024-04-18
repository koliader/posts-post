package rabbitmq

import (
	"fmt"

	"github.com/koliader/posts-post.git/internal/util"
	"github.com/streadway/amqp"
)

type UpdateEmailMessage struct {
	Email    string `json:"email"`
	NewEmail string `json:"newEmail"`
}

type Client struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewClient(config util.Config) (*Client, error) {
	conn, err := amqp.Dial(config.RbmUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		conn:    conn,
		Channel: channel,
	}, nil
}
func (s *Client) CreateQueue(channelName string) error {
	_, err := s.Channel.QueueDeclare(channelName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendMessage(queueName string, message []byte) error {
	err := c.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing message to RabbitMQ: %w", err)
	}
	return nil
}

func (c *Client) GetMessages(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := c.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (c *Client) Close() error {
	if c.Channel != nil {
		if err := c.Channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
