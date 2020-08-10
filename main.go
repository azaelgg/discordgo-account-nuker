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

func getAllGuildsWorker(Session *discordgo.Session) []string {
	var guildIDS []string
	for _, guild := range Session.State.Guilds {
		guildIDS = append(guildIDS, guild.ID)
	}
	return guildIDS
}

func getAllFriendsWorker(Session *discordgo.Session) []string {
	var friendIDS []string
	relationships, err := Session.RelationshipsGet()
	if err != nil {
		fmt.Println("[\u001b[31m-\u001b[0m] := Error retrieving friend list,", err)
	}

	for _, friend := range relationships {
		friendIDS = append(friendIDS, friend.ID)
	}
	return friendIDS
}

func main() {
	
	Session, _ = discordgo.New("")
	Session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36"
	err := Session.Open()
	if err != nil {
		fmt.Println("[\u001b[31m-\u001b[0m] := Error opening connection,", err)
		return
	}
	friendList := getAllFriendsWorker(Session)
	for _, friend := range friendList {
		Session.RelationshipDelete(friend)
		/*
		if err != nil {
			fmt.Println("[\u001b[31m-\u001b[0m] := Error removing relationship : " + friend, err)
		}
		*/
		fmt.Println("[\u001b[32m+\u001b[0m] := Relationship removed : " + friend)
	}

	guildsList := getAllGuildsWorker(Session)
	for _, guild := range guildsList {
		err := Session.GuildLeave(guild)
		if err != nil {
			fmt.Println("[\u001b[31m-\u001b[0m] := Error leaving guild,", err)
		}
		fmt.Println("[\u001b[32m+\u001b[0m] := Guild left : " + guild)
	}

}
