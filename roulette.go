//module for russian roullette functionality
package main

import (
	"fmt"
	"math/rand"
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
	fmt.Println("Russian roulette triggered for user", m.Author)
	sendLog(fmt.Sprintln("Russian roulette triggered for user", m.Author))
	bullet := rand.Intn(5)
	guild := m.GuildID
	fmt.Println("revolver landed on chamber: ", bullet)
	sendLog(fmt.Sprintln("revolver landed on chamber: ", bullet))
	if bullet == 0 {
		//kick the user - Will only work if the bot has higher permissions than the user
		session.ChannelMessageSend(m.ChannelID, killgif)
		fmt.Println(session.GuildMemberDeleteWithReason(guild, m.Author.ID, killmsg))
		fmt.Println("Kicking user with ID", m.Author.ID)
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
	author := m.Author.ID
	reviveSeconds := 5
	reviveMessage := "Stay away from the light. We still need you!"
	userChannel, err := session.UserChannelCreate(author)
	var invite discordgo.Invite
	invite.MaxAge = 20
	invite.MaxUses = 1
	invite.Temporary = true
	userInvite, err := session.ChannelInviteCreate(m.ChannelID, invite)
	fmt.Println(err)
	time.Sleep(time.Duration(reviveSeconds) * time.Second)
	userInviteLink := "https://discord.gg/" + userInvite.Code
	session.ChannelMessageSend(userChannel.ID, reviveMessage)
	time.Sleep(time.Duration(reviveSeconds/2) * time.Second)
	session.ChannelMessageSend(userChannel.ID, userInviteLink)
}
