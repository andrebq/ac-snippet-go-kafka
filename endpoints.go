package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/rs/xid"
)

func channelPost(e echo.Context) error {
	channel := e.Param("channel")
	if channel == "" {
		return badRequest("missing channel name")
	}

	user, ok := asString(e.Request().Context().Value(authKey))
	if !ok || user == "" {
		panic("protected function called without proper authentication")
	}

	var newMessage Message
	if err := e.Bind(&newMessage); err != nil {
		return err
	}
	if newMessage.Text == "" {
		return badRequest("missing text")
	}
	newMessage.Destination = "chan:" + channel
	newMessage.Sender = "user:" + user
	newMessage.ID = xid.New().String()
	newMessage.ServerTime = time.Now()

	if err := publish(e.Request().Context(), []byte(newMessage.ID), newMessage, "global-inbox"); err != nil {
		return badGateway("unable to publish your message try again later")
	}

	e.Response().Header().Add("Location", fmt.Sprintf("//messages/%v", newMessage.ID))
	e.String(http.StatusCreated, newMessage.ID)
	return nil
}

func asString(in interface{}) (string, bool) {
	str, ok := in.(string)
	return str, ok
}
