package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	initPubSub()
	ctx := watchSigKill(context.Background())

	go printGlobalInbox(ctx)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logrus.StandardLogger().Writer(),
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(user, pwd string, c echo.Context) (bool, error) {
		if user == "" {
			return false, nil
		}
		c.SetRequest(c.Request().WithContext(context.WithValue(
			c.Request().Context(), authKey, user)))
		return true, nil
	}))

	e.POST("/channels/:channel/message", channelPost)

	go shutdown(ctx, e)

	e.Start("localhost:9099")
}

func shutdown(ctx context.Context, e *echo.Echo) {
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}

func printGlobalInbox(ctx context.Context) {
	msgs, err := subscribe(ctx, "global-inbox", -2)
	for {
		select {
		case <-ctx.Done():
			logrus.WithError(ctx.Err()).Error("done")
			return
		case e := <-err:
			logrus.WithError(e).Error("error reading data")
		case m := <-msgs:
			var data Message
			err := json.Unmarshal(m.Value, &data)
			if err != nil {
				logrus.WithField("offset", m.Offset).WithError(err).Error("error decoding message")
			} else {
				logrus.WithField("offset", m.Offset).WithField("message", data).
					WithField("delay", time.Now().Sub(data.ServerTime).String()).Info()
			}
		}
	}
}

func watchSigKill(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		cancel()
		signal.Stop(ch)
		signal.Ignore(os.Interrupt, os.Kill)
	}()
	return ctx
}
