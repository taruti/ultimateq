// The ultimateq bot framework.
package main

import (
	"bytes"
	"github.com/aarondl/ultimateq/bot"
	"github.com/aarondl/ultimateq/config"
	"github.com/aarondl/ultimateq/irc"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Handler struct {
}

func (h Handler) PrivmsgUser(m *irc.Message, endpoint irc.Endpoint) {
	if strings.Split(m.Sender, "!")[0] == "Aaron" {
		endpoint.Send(m.Message())
	}
}

func (h Handler) PrivmsgChannel(m *irc.Message, endpoint irc.Endpoint) {
	if m.Message() == "hello" {
		endpoint.Privmsg(m.Target(), "Hello to you too!")
	} else {
		lock.Lock()
		chain.Build(
			bytes.NewBuffer(
				[]byte(m.Message()),
			),
		)
		endpoint.Privmsg(m.Target(), chain.Generate(100))
		lock.Unlock()
	}
}

func conf(c *config.Config) *config.Config {
	c. // Defaults
		Nick("nobody__").
		Altnick("nobody_").
		Realname("there").
		Username("guy").
		Userhost("friend").
		NoReconnect(true)

	c. // First server
		Server("irc.gamesurge.net1").
		Host("localhost").
		Nick("Aaron").
		Altnick("nobody1").
		ReconnectTimeout(5)

	c. // Second Server
		Server("irc.gamesurge.net2").
		Host("localhost").
		Nick("nobody2")

	return c
}

var chain = NewChain(2)
var lock = sync.Mutex{}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
	log.SetOutput(os.Stdout)

	b, err := bot.CreateBot(bot.ConfigureFunction(conf))
	if err != nil {
		log.Println(err)
	}

	b.Register(irc.PRIVMSG, Handler{})

	ers := b.Connect()
	if len(ers) != 0 {
		log.Println(ers)
		return
	}
	b.Start()

	<-time.After(20 * time.Second)
	c := conf(config.CreateConfig())
	c.RemoveServer("irc.gamesurge.net1")
	c.GlobalContext().
		Channels("#test").
		Server("irc.gamesurge.net3").
		Host("localhost").
		Nick("super").
		ServerContext("irc.gamesurge.net2").
		Nick("heythere").
		Channels("#hithere")
	if !c.IsValid() {
		c.DisplayErrors()
		log.Fatal("Config error")
	}
	//b.ReplaceConfig(c)

	b.WaitForHalt()
	b.Stop()
	b.Disconnect()
	<-time.After(10 * time.Second)
}
