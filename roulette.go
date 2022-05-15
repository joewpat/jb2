//module for russian roullette functionality
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

//taunt the user before deciding their fate
func taunt() string {
	tauntmsg := "Let's decide your fate"
	//let's decide your fate
	return tauntmsg
}

//either kick or save the user
func roulette(m *discordgo.Message, session *discordgo.Session) string {
	taunt()
	rand.Seed(time.Now().UnixNano()) //init random
	safemsg := "Looks like you live... this time"
	killmsg := "nothing personnel, kid"
	safegif := searchGifs("relieved")
	killgif := searchGifs("gunshot")
	sendLog(fmt.Sprintln("Russian roulette triggered for user", m.Author))
	bullet := rand.Intn(5)
	guild := m.GuildID
	sendLog(fmt.Sprintln("suicide mode triggered"))
	if strings.Contains(m.Content[3:], "suicide") {
		sendLog(fmt.Sprintln("revolver landed on chamber: ", bullet))
		bullet = 0
	}
	if bullet == 0 {
		//kick the user - Will only work if the bot has higher permissions than the user
		session.ChannelMessageSend(m.ChannelID, killgif)
		fmt.Println(session.GuildMemberDeleteWithReason(guild, m.Author.ID, killmsg))
		sendLog(fmt.Sprintln("Kicking user with ID", m.Author.ID))
		//add functionality to PM the user to bring them back to life
		go revive(m, session)
		return killmsg
	} else {
		//tell the user they're safe
		session.ChannelMessageSend(m.ChannelID, safegif)
		return safemsg
	}
}

func revive(m *discordgo.Message, session *discordgo.Session) {
	reviveSeconds := 5                                              //Time between killing the user and reviving them
	reviveMessage := "Stay away from the light. We still need you!" //Message to send to user before sending them the invite link

	author := m.Author.ID
	userChannel, err := session.UserChannelCreate(author)
	if err != nil {
		fmt.Println("error DMing user,", err)
		return
	}
	var invite discordgo.Invite
	invite.MaxAge = 20
	invite.MaxUses = 1
	invite.Temporary = true
	userInvite, err := session.ChannelInviteCreate(m.ChannelID, invite)
	if err != nil {
		fmt.Println("error creating channel invite,", err)
		return
	}
	time.Sleep(time.Duration(reviveSeconds) * time.Second)
	sendLog("Revive timer up, sending revive message")
	userInviteLink := "https://discord.gg/" + userInvite.Code
	session.ChannelMessageSend(userChannel.ID, reviveMessage)
	time.Sleep(time.Duration(reviveSeconds/2) * time.Second)
	session.ChannelMessageSend(userChannel.ID, userInviteLink)
	sendLog("Sent invite message with link: " + userInviteLink)
}
