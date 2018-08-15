package main

import (
	"net/http"
	"flag"
	"os"
	"github.com/dghubble/sling"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	user = flag.String("user", "", "Username")
	channel = flag.String("channel", "demo", "Destination channel")
	host = flag.String("host", "http://localhost:9099", "API host")
	message = flag.String("msg", "", "Message to send")
	help = flag.Bool("h", false, "Help")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	body := struct {
		Message string `url:"Text"`
	}{
		Message: *message,
	}
	resp, err := sling.New().Base(*host).
		SetBasicAuth(*user, *user).
		Post(fmt.Sprintf("/channels/%s/message", *channel)).BodyForm(body).Receive(nil, nil)
	if err != nil {
		logrus.WithError(err).Fatal("unable to send message")
	} else if resp.StatusCode != http.StatusCreated {
		logrus.WithField("statusCode", resp.StatusCode).Fatal("Unexpected status code")
	}
	io.Copy(os.Stdout, resp.Body)
}