package main

import (
	"fmt"
	// "net/http"
	"github.com/bwmarrin/discordgo"	
)

var (
	Session           *discordgo.Session 
	Token             string
)

func leaveAllGuildsWorker(Session *discordgo.Session) {
	for _, guild := range Session.State.Guilds {
		err := Session.GuildLeave(guild.ID)
		if err != nil {
			fmt.Println("Error leaving server,", err)
		}
	}
}

func main() {
	Session, _ = discordgo.New("")
	// Session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36"
	err := Session.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	leaveAllGuildsWorker(Session)
}
