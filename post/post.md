TL;DR In this post we saw what libraies we can use to interact with Kafka and how to create a simple producer/consumer pair. This also paved the way to further posts expanding this client in order to create a simple, but functional, online-chat service.

# Introduction

You probably already heard about [https://kafka.apache.org/](Kafka), you known that Apache project started at LinkedIn that brought the concent of a distributed commit-log to the masses. With that tool and some Go programming we can build a very crude but fully functional online-chat application.

To start we need a way to run Kafka in our computers or find a Kafka cluster lingering around. Kafka usually runs in a clustered environmente managed by Zookeeper with one leader a one or more followers to allow high-availability, that kind of deployment comes with its onw set of problems so we will stay away from that and use a simpler solution running a single-node kafka cluster.

One way is to use a [https://raw.githubusercontent.com/andrebq/ac-snippet-go-kafka/master/docker-compose.yaml](docker-compose.yml) file and start it using **docker-compose**. One important detail is that Kafka will advertise its address to clients so your client needs to resolve the address/hostname. That means adding entries to your **/etc/hosts** or **c:\windows\system32\drivers\etc\hosts** to map the advertised hostname to a valid IP.

```
    127.0.0.1 kafkaserver
```

If you are using **docker machine** on windows or the docker host has a different IP you might need to map **kafkahost** to a different IP.

Now that kafka preparation is done, to compile the code you will need to download a [https://golang.org/dl/](recent) version of [https://golang.org](Go) and install either: [https://github.com/golang/vgo](vgo) or [https://github.com/golang/dep](dep). If you are using Go 1.11 or further vgo support is already built into the **go** binary so you don't need any extra download.

# A kafka client for Go

There are many kafka clients for Go:

* sarama: from Shopify is probably the most famous one but thas low-level API which makes things more complicated.
* confluent-kafka-go: is written by confluent.io and is a CGO-wrapper around librdkafka which is the official library, but [https://dave.cheney.net/2016/01/18/cgo-is-not-go](CGO is not Go) so lets avoid it.
* segment.io: this is a pure-go client from Segment.io, it exposes a simpler API and abstracts most of the complexity of dealing with Kafka. This is the one we will use.

# Kafka 101

Kafka basic data-structure is a append-only list of entries with optional keys. That list is stored on persistent disk. When Kafka is running out of space it will start a compact process to remove old entries from the list.

For entries without a key, each entry is considered unique and will be sent to consumers but entries associated with a key only the latest value is kept in the commit-log, this turns Kafka into a key-value database.

Producers (Writers in kafka-go parlance) will push messages to one Kafka broker informing: a partition; a key (if available); and the actual data. Depending on how the connection was established it will wait until the message gets replicated to at least N nodes (where N >= 1).

Consumers (Readers in kafka-go parlance) on the other hand will connect to Kafka and ask for messages from a specific partition, they may also provide:

* a GroupID to allow messages to be distribuited across a group of consumers;
* or an offset to control the starting point (oldest available, specific point in the log, most recent).

# Our demo

Our system consists of a simple HTTP api where an user can _POST_ messages to a _channel_, those messages are posted to a _global_ Kafka topic which might then be consumed by interested parties. To validate the solution we have one consumer reading the messages and measuring the delta between posting the message and receiving it. Other clients should do more interesting things like: indexing it for full-text search; routing them to other topics after some processing; sending push-notifications; etc...

The code responsible for publishing the message is quite simple:

```
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
```

Reading messages is also easy:

```
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
```

Notice we are reading them inside a **goroutine** so consumers of the api will interact with it over **go channels**. One important thing to notice is the use of **reader.FetchMessage** instead of **ReadMessage**, the later will commit the message Offset to the GroupID topic and when reconnected it will start from the last commited offset **FetchMessage** on the other hand only downloads the data and moves the pointer locally.

Another important call is **reader.SetOffset(offset)**, this method will change our starting point:

* -1: read from the oldest offset available
* -2: read from the newest offset available (aka, only new messages)

# Wrapping-up

In this post we learned about how to connect your Go application with a Kafka cluster to read/write messages so you can create your very own online-chat application. All the code is available on [https://github.com/andrebq/ac-snippet-go-kafka](Github) so you can download and play with it. 