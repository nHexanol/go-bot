package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var (
	RemoveCommand = flag.Bool("rmcmd", true, "hi")
	GuildID       = flag.String("guild", "", "Guild ID")
)

func setInterval(cb func(), t time.Duration) chan bool {
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(t):
				cb()
				fmt.Println("tick")
			case <-stop:
				return
			}
		}
	}()

	return stop
}

type Player struct {
	Kind      string `json:"kind"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
	Data      []struct {
		Username string `json:"username"`
		UUID     string `json:"uuid"`
		Rank     string `json:"rank"`
		Meta     struct {
			FirstJoin time.Time `json:"firstJoin"`
			LastJoin  time.Time `json:"lastJoin"`
			Location  struct {
				Online bool   `json:"online"`
				Server string `json:"server"`
			} `json:"location"`
			Playtime float64 `json:"playtime"`
			Tag      struct {
				Display bool   `json:"display"`
				Value   string `json:"value"`
			} `json:"tag"`
			Veteran bool `json:"veteran"`
		} `json:"meta"`
		Characters map[string]struct {
			Type     string `json:"type"`
			Level    int    `json:"level"`
			Dungeons struct {
				Completed int `json:"completed"`
				List      []struct {
					Name      string `json:"name"`
					Completed int    `json:"completed"`
				} `json:"list"`
			} `json:"dungeons"`
			Raids struct {
				Completed int `json:"completed"`
				List      []struct {
					Name      string `json:"name"`
					Completed int    `json:"completed"`
				} `json:"list"`
			} `json:"raids"`
			Quests struct {
				Completed int      `json:"completed"`
				List      []string `json:"list"`
			} `json:"quests"`
			ItemsIdentified int `json:"itemsIdentified"`
			MobsKilled      int `json:"mobsKilled"`
			Pvp             struct {
				Kills  int `json:"kills"`
				Deaths int `json:"deaths"`
			} `json:"pvp"`
			BlocksWalked int `json:"blocksWalked"`
			Logins       int `json:"logins"`
			Deaths       int `json:"deaths"`
			Playtime     int `json:"playtime"`
			Gamemode     struct {
				Craftsman bool `json:"craftsman"`
				Hardcore  bool `json:"hardcore"`
				Ironman   bool `json:"ironman"`
				Hunted    bool `json:"hunted"`
			} `json:"gamemode"`
			Skills struct {
				Strength     int `json:"strength"`
				Dexterity    int `json:"dexterity"`
				Intelligence int `json:"intelligence"`
				Defence      int `json:"defence"`
				Defense      int `json:"defense"`
				Agility      int `json:"agility"`
			} `json:"skills"`
			Professions struct {
				Alchemism struct {
					Level int `json:"level"`
					Xp    int `json:"xp"`
				} `json:"alchemism"`
				Armouring struct {
					Level int `json:"level"`
					Xp    int `json:"xp"`
				} `json:"armouring"`
				Combat struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"combat"`
				Cooking struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"cooking"`
				Farming struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"farming"`
				Fishing struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"fishing"`
				Jeweling struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"jeweling"`
				Mining struct {
					Level int `json:"level"`
					Xp    int `json:"xp"`
				} `json:"mining"`
				Scribing struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"scribing"`
				Tailoring struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"tailoring"`
				Weaponsmithing struct {
					Level int `json:"level"`
					Xp    int `json:"xp"`
				} `json:"weaponsmithing"`
				Woodcutting struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"woodcutting"`
				Woodworking struct {
					Level int     `json:"level"`
					Xp    float64 `json:"xp"`
				} `json:"woodworking"`
			} `json:"professions"`
			Discoveries      int  `json:"discoveries"`
			EventsWon        int  `json:"eventsWon"`
			PreEconomyUpdate bool `json:"preEconomyUpdate"`
		} `json:"characters"`
		Guild struct {
			Name string `json:"name"`
			Rank string `json:"rank"`
		} `json:"guild"`
		Global struct {
			BlocksWalked    int `json:"blocksWalked"`
			ItemsIdentified int `json:"itemsIdentified"`
			MobsKilled      int `json:"mobsKilled"`
			TotalLevel      struct {
				Combat     int `json:"combat"`
				Profession int `json:"profession"`
				Combined   int `json:"combined"`
			} `json:"totalLevel"`
			Pvp struct {
				Kills  int `json:"kills"`
				Deaths int `json:"deaths"`
			} `json:"pvp"`
			Logins      int `json:"logins"`
			Deaths      int `json:"deaths"`
			Discoveries int `json:"discoveries"`
			EventsWon   int `json:"eventsWon"`
		} `json:"global"`
		Ranking struct {
			Guild  interface{} `json:"guild"`
			Player struct {
				Solo struct {
					Combat         interface{} `json:"combat"`
					Woodcutting    interface{} `json:"woodcutting"`
					Mining         interface{} `json:"mining"`
					Fishing        interface{} `json:"fishing"`
					Farming        interface{} `json:"farming"`
					Alchemism      interface{} `json:"alchemism"`
					Armouring      interface{} `json:"armouring"`
					Cooking        interface{} `json:"cooking"`
					Jeweling       interface{} `json:"jeweling"`
					Scribing       interface{} `json:"scribing"`
					Tailoring      interface{} `json:"tailoring"`
					Weaponsmithing interface{} `json:"weaponsmithing"`
					Woodworking    interface{} `json:"woodworking"`
					Profession     interface{} `json:"profession"`
					Overall        interface{} `json:"overall"`
				} `json:"solo"`
				Overall struct {
					All        interface{} `json:"all"`
					Combat     interface{} `json:"combat"`
					Profession interface{} `json:"profession"`
				} `json:"overall"`
			} `json:"player"`
			Pvp interface{} `json:"pvp"`
		} `json:"ranking"`
	} `json:"data"`
}

type Guild struct {
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
	Members []struct {
		Name           string    `json:"name"`
		UUID           string    `json:"uuid"`
		Rank           string    `json:"rank"`
		Contributed    int64     `json:"contributed"`
		Joined         time.Time `json:"joined"`
		JoinedFriendly string    `json:"joinedFriendly"`
	} `json:"members"`
	Xp              int       `json:"xp"`
	Level           int       `json:"level"`
	Created         time.Time `json:"created"`
	CreatedFriendly string    `json:"createdFriendly"`
	Territories     int       `json:"territories"`
	Banner          struct {
		Base      string `json:"base"`
		Tier      int    `json:"tier"`
		Structure string `json:"structure"`
		Layers    []struct {
			Colour  string `json:"colour"`
			Pattern string `json:"pattern"`
		} `json:"layers"`
	} `json:"banner"`
	Request struct {
		Timestamp int `json:"timestamp"`
		Version   int `json:"version"`
	} `json:"request"`
}

func main() {

	debug.SetPanicOnFault(false)

	flag.Parse()
	client, e := discordgo.New()
	if e != nil {
		panic(e)
	}

	sh := make(chan os.Signal, 1)
	signal.Notify(sh, syscall.SIGSEGV, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		select {
		case <-sh:
			var errmsg = "```golang\nInterrupted\n\n" + string(debug.Stack())
			client.ChannelMessageSend("1010645092232614002", errmsg+"\n```")
			fmt.Println(string(debug.Stack()))
		}
	}()

	client.AddHandler(message)

	client.Identify.Intents = discordgo.IntentsGuildMessages

	err2 := client.Open()
	if err2 != nil {
		panic(err2)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	if ee := recover(); ee != nil {
		client.ChannelMessageSend("1010645092232614002", "```golang\n"+e.Error()+"\n```")
	}
	/*
		var query = func() {

			res, err := http.Get("https://api.wynncraft.com/public_api.php?action=guildStats&command=Sins+of+Seedia")
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			body, e := ioutil.ReadAll(res.Body)
			if e != nil {
				panic(e)
			}

			var guild Guild
			json.Unmarshal(body, &guild)
			if len(oldMember) == 0 {
				oldMember = guild.Members
				return
			}
			// member joined
			if len(oldMember) < len(guild.Members) {
				for k, v := range guild.Members {
					if v.Name != oldMember[k].Name {
						fmt.Println(guild.Members[k].Name + " joined the guild")
					}
				}
			}

			// member left
			if len(oldMember) > len(guild.Members) {
				for k, v := range guild.Members {
					if v.Name != oldMember[k].Name {
						fmt.Println(oldMember[k].Name + " left the guild")
					}
				}
			}
		}
		setInterval(query, 60000)
	*/
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, "!stats ") {
		var u = strings.Replace(m.Content, "!stats ", "", -1)
		query(u, s, m)
	} else if strings.Contains(m.Content, "!char") {
		var u = strings.Replace(m.Content, "!char ", "", -1)
		queryCharacters(u, s, m)
	} else if strings.Contains(m.Content, "!guild") {
		var guildName = strings.Replace(m.Content, "!guild ", "", -1)
		var guildNameREAL = strings.Replace(guildName, " ", "+", -1)

		onlineGuildMember(guildNameREAL)
	} else if strings.Contains(m.Content, "!eval") && m.Author.ID == "246865469963763713" {

		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)

		done, err := i.Eval(strings.Replace(m.Content, "!eval ", "", -1))
		if err != nil {
			var errEmbed = &discordgo.MessageEmbed{
				Title:       "Evaluate error",
				Description: "```golang\n> " + strings.Replace(m.Content, "!eval ", "", -1) + "\n\n< " + strings.Replace(err.Error(), "_.go", "<eval>", -1) + "\n\nstacktrace:\n" + string(debug.Stack()) + "\n```",
				Color:       0xff0000,
			}
			s.ChannelMessageSendEmbed(m.ChannelID, errEmbed)
		} else {
			var evalEmbed = &discordgo.MessageEmbed{
				Title:       "Evaluate completed",
				Description: "```golang\n> " + strings.Replace(m.Content, "!eval ", "", -1) + "\n\n< " + fmt.Sprint(done) + "\n\nstacktrace :\n\n" + string(debug.Stack()) + "\n```",
				Color:       0x00ff00,
				Timestamp:   time.Now().Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text: m.Author.Username,
				},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, evalEmbed)

		}
	}
}

func queryCharacters(username string, s *discordgo.Session, m *discordgo.MessageCreate) {

	var sb strings.Builder

	var getString string = "https://api.wynncraft.com/v2/player/" + username + "/stats"
	var respond, error = http.Get(getString)
	if error != nil {
		fmt.Println("Error")
	}
	defer respond.Body.Close()
	var body, err = ioutil.ReadAll(respond.Body)
	if err != nil {
		fmt.Println("Error")
	}

	var user Player
	json.Unmarshal(body, &user)

	var current int = 1
	sb.WriteString("```\n")
	for k := range user.Data[0].Characters {
		sb.WriteString("\n[ ")
		sb.WriteString(fmt.Sprint(current))
		sb.WriteString(" ] ")
		sb.WriteString(user.Data[0].Characters[k].Type)
		sb.WriteString("\n")
		sb.WriteString("Total Level : ")
		sb.WriteString(fmt.Sprint(user.Data[0].Characters[k].Level))
		sb.WriteString("\nCombat Level : ")
		sb.WriteString(fmt.Sprint(user.Data[0].Characters[k].Professions.Combat.Level))
		sb.WriteString(" [ ")
		sb.WriteString(fmt.Sprint(user.Data[0].Characters[k].Professions.Combat.Xp))
		sb.WriteString("% ]\n")
		current++
	}
	sb.WriteString("\n```")

	var embed = &discordgo.MessageEmbed{
		Title:       user.Data[0].Username + "'s Character(s)",
		Color:       0xffff00,
		Description: sb.String(),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func query(username string, s *discordgo.Session, m *discordgo.MessageCreate) {

	var getString string = "https://api.wynncraft.com/v2/player/" + username + "/stats"
	var respond, error = http.Get(getString)
	if error != nil {
		fmt.Println("Error")
	}
	defer respond.Body.Close()
	var body, err = ioutil.ReadAll(respond.Body)
	if err != nil {
		fmt.Println("Error")
	}

	var user Player
	json.Unmarshal(body, &user)

	var rn = time.Now().UnixMilli()

	var lastSeenStr = ""
	var lastSeenDay, lastSeenHour, lastSeenMinute, lastSeenSecond, lastSeenMilisec int64
	if user.Data[0].Meta.Location.Online {
		lastSeenStr = "Online at " + user.Data[0].Meta.Location.Server
	} else {
		var lastSeenMilis = time.Unix(0, user.Data[0].Meta.LastJoin.UnixMilli()*int64(time.Millisecond)).UnixMilli()
		var howLong = rn - lastSeenMilis
		lastSeenDay = howLong / 86400000
		lastSeenHour = (howLong - lastSeenDay*86400000) / 3600000
		lastSeenMinute = (howLong - lastSeenDay*86400000 - lastSeenHour*3600000) / 60000
		lastSeenSecond = (howLong - lastSeenDay*86400000 - lastSeenHour*3600000 - lastSeenMinute*60000) / 1000
		lastSeenMilisec = howLong - lastSeenDay*86400000 - lastSeenHour*3600000 - lastSeenMinute*60000 - lastSeenSecond*1000
		lastSeenStr = strconv.FormatInt(lastSeenDay, 10) + "d " + strconv.FormatInt(lastSeenHour, 10) + "h " + strconv.FormatInt(lastSeenMinute, 10) + "m " + strconv.FormatInt(lastSeenSecond, 10) + "s " + strconv.FormatInt(lastSeenMilisec, 10) + "ms"
	}
	fmt.Println(lastSeenStr)
	var firstJoinMilis = time.Unix(0, user.Data[0].Meta.FirstJoin.UnixMilli()*int64(time.Millisecond)).UnixMilli()
	var howLong2 = rn - firstJoinMilis
	var firstJoinDay = howLong2 / 86400000
	var firstJoinHour = (howLong2 - firstJoinDay*86400000) / 3600000
	var firstJoinMinute = (howLong2 - firstJoinDay*86400000 - firstJoinHour*3600000) / 60000
	var firstJoinSecond = (howLong2 - firstJoinDay*86400000 - firstJoinHour*3600000 - firstJoinMinute*60000) / 1000
	var firstJoinMilisec = howLong2 - firstJoinDay*86400000 - firstJoinHour*3600000 - firstJoinMinute*60000 - firstJoinSecond*1000

	type PlayerGuild struct {
		Name string `json:"name"`
		Rank string `json:"rank"`
	}

	var plgu = PlayerGuild{
		Name: user.Data[0].Guild.Name,
		Rank: user.Data[0].Guild.Rank,
	}

	var guildName, guildStar, guildRank string

	if plgu.Name != "" {
		guildName = user.Data[0].Guild.Name
		guildRank = user.Data[0].Guild.Rank
		switch guildRank {
		case "RECRUIT":
			guildStar = " **-** "
		case "RECRUITER":
			guildStar = " **\\*** "
		case "CAPTAIN":
			guildStar = " **\\*\\*** "
		case "STRATEGIST":
			guildStar = " **\\*\\*\\*** "
		case "CHIEF":
			guildStar = " **\\*\\*\\*\\*** "
		case "OWNER":
			guildStar = " **\\*\\*\\*\\*\\*** "
		default:
			guildStar = user.Data[0].Guild.Rank
		}
	} else {
		guildName = "<nil>"
		guildRank = "<nil>"
		guildStar = "<nil>"
	}
	var name string
	if user.Data[0].Meta.Tag.Display {
		name = user.Data[0].Username + "  **[ " + user.Data[0].Meta.Tag.Value + " ]**"
	} else {
		name = user.Data[0].Username
	}

	playtime := user.Data[0].Meta.Playtime
	var newTime float64 = playtime * 4.7 / 60
	newTime = math.Round(newTime*100) / 100

	var embed = &discordgo.MessageEmbed{

		Title: name,
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Guild",
				Value: guildName + " " + " [" + guildStar + "]",
			},
			{
				Name:  "Playtime",
				Value: fmt.Sprintf("%d", int32(math.Round(newTime))) + " hours",
			},
			{
				Name:  "Logins / Deaths",
				Value: fmt.Sprint(user.Data[0].Global.Logins) + " / " + fmt.Sprint(user.Data[0].Global.Deaths),
			},
			{
				Name:  "Last Seen",
				Value: lastSeenStr,
			},
			{
				Name:  "First joined",
				Value: fmt.Sprint(firstJoinDay) + " days " + fmt.Sprint(firstJoinHour) + " hrs " + fmt.Sprint(firstJoinMinute) + " min " + fmt.Sprint(firstJoinSecond) + " s " + fmt.Sprint(firstJoinMilisec) + " ms ago",
			},
		},
	}

	var _, err3 = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err3 != nil {
		fmt.Println("Error")
	}
}

func onlineGuildMember(GuildName string) {

	var member []string

	// query guild object
	res, e := http.Get("https://api.wynncraft.com/public_api.php?action=guildStats&command=" + GuildName)
	if e != nil {
		panic(e)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var guild Guild

	json.Unmarshal(body, &guild)

	for _, m := range guild.Members {
		_ = append(member, m.Name)
	}

	fmt.Println(member, 12341234)

	// query online player and match all online username to guild member list
	res, e = http.Get("https://api.wynncraft.com/public_api.php?action=onlinePlayers")
	if e != nil {
		panic(e)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data = make(map[string][]string)

	json.Unmarshal([]byte(body), &data)

	delete(data, "request")

	var keys = []string{}

	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	var str = new(strings.Builder)
	str.Write([]byte("Online guild member: "))
	// find online guild member in online player list
	for _, k := range guild.Members {
		for _, v := range data["result"] {
			if k.Name == v {
				fmt.Println(k.Name, " is online")
			}
		}
	}

	for _, k := range keys {
		for _, v := range data[k] {
			fmt.Println(k, " : ", v)
		}
	}

}
