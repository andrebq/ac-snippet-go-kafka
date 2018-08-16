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
* confluent-kafka-go: is written by confluent.io and is a CGO-wrapper around librdkafka which is the official library, but [https://dave.cheney.net/2016/01/18/cgo-is-not-go](CGO is not Go) so we'll skip it.
* segment.io: this is a pure-go client from Segment.io, it exposes a simpler API and abstracts most of the complexity of dealing with Kafka. This is the one we will use.

