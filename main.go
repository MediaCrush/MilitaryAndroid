package main

import (
	"fmt"
	"github.com/jdiez17/go-irc"
	"strings"
	"time"
)

func main() {
    err := loadConfig("config.json")
    if err != nil {
        fmt.Println("Error reading the configuration: " + err.Error())
        return
    }

	conn, err := irc.NewConnection(Config.IRC.Server, int(Config.IRC.Port))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer conn.Close()

	conn.LogIn(irc.Identity{Nick: Config.Nick})

	conn.AddHandler(irc.MOTD_END, func(c *irc.Connection, e *irc.Event) {
        if Config.NickServPassword != "" {
            c.Privmsg("NickServ", "identify " + Config.NickServPassword)
        }

        for _, channel := range Config.Channels {
            c.Join(channel)
        }
	})
    conn.AddHandler(irc.PRIVMSG, expandGithubIssue) 

	bot := irc.NewBot(conn)
	bot.AddCommand("echo", func(c *irc.Connection, e *irc.Event) {
		message := strings.Join(e.Params, " ")
		e.React(c, message)
	})

	bot.AddCommand("portal", portalCommandHandler)
	bot.AddCommand("compliment", complimentCommandHandler)

	for {
		<-time.After(1 * time.Second)
	}
}
