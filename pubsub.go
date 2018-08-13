package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"
)

type (
	writer struct {
		sync.Mutex
		kw *kafka.Writer
	}
)

var (
	writers map[string]*writer

	writerLock sync.Mutex
)

func initPubSub() {
	writers = make(map[string]*writer)
}

func getBrokers() []string {
	return []string{"localhost:9092"}
}

func getWriter(topic string) *writer {
	writerLock.Lock()
	defer writerLock.Unlock()

	w, ok := writers[topic]
	if !ok {
		w = &writer{}
		w.kw = kafka.NewWriter(kafka.WriterConfig{
			Topic:   topic,
			Brokers: getBrokers(),
		})
		writers[topic] = w
	}

	return w
}

func subscribe(ctx context.Context, topic string, offset int64) (<-chan kafka.Message, <-chan error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Topic:    topic,
		Brokers:  getBrokers(),
		MinBytes: 0,
		MaxBytes: 10e6,
	})
	messages := make(chan kafka.Message, 1)
	errors := make(chan error, 1)

	sendErr := func(ctx context.Context, errors chan error, err error) {
		select {
		case <-ctx.Done():
			select {
			case errors <- err:
			default:
			}
		case errors <- err:
		}
	}
	sendMsg := func(ctx context.Context, messages chan kafka.Message, msg kafka.Message) {
		select {
		case <-ctx.Done():
			select {
			case errors <- ctx.Err():
			default:
			}
		case messages <- msg:
		}
	}
	go func() {
		defer close(messages)
		defer close(errors)
		err := reader.SetOffset(offset)
		if err != nil {
			sendErr(ctx, errors, err)
			return
		}
		for {
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				sendErr(ctx, errors, err)
				continue
			}
			sendMsg(ctx, messages, msg)
		}
	}()
	return messages, errors
}

func publish(ctx context.Context, key []byte, data interface{}, topic string) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w := getWriter(topic)

	return w.kw.WriteMessages(ctx, kafka.Message{
		Value: buf,
		Key:   key,
	})
}
