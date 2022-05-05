//module for russian roullette functionality
package main

import (
    "math/rand"
	"github.com/bwmarrin/discordgo"
    "fmt"
)

//taunt the user before deciding their fate
func taunt(a string) string {
    tauntmsg := "Let's decide your fate"
    //let's decide your fate
    return tauntmsg
}
//either kick or save the user
func roulette(m *discordgo.Message, session *discordgo.Session) string { 
    safemsg := "Looks like you live... this time"
    killmsg := "nothing personnel, kid"
	fmt.Println("Russian roulette triggered for user", m.Author)
    bullet := rand.Intn(5)
    guild := m.GuildID
    if bullet == 0 {
        //kick the user - Will only work if the bot has higher permissions than the user
        fmt.Println(session.GuildMemberDeleteWithReason(guild,m.Author.ID,killmsg))
        fmt.Println("Kicking user with ID", m.Author.ID)
        return killmsg
    } else
    {
        //tell the user they're safe
        return safemsg
    }
}
