package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/ini.v1"
)

type User struct {
	Name      string `json:"name"`
	Followers int    `json:"size"`
	Color     string `json:"color"`
}

type Link struct {
	Source    string `json:"source"`
	Target    string `json:"target"`
	CreatedAt string `json:"created_at"`
	Key       string `json:"key"`
}

func readUsers(file string) []*User {
	var us []*User
	raw, err := ioutil.ReadFile(file)
	if err == nil {
		json.Unmarshal(raw, &us)
	}
	return us
}

func main() {
	t := time.Now()

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	cfgfile := flags.String("conf", "", "Mandatory ini config file")
	nodes := flags.String("n", "./data/nodes.json", "Nodes file.\n-n \"\" to create new file\n")
	searchQuery := flags.String("q", "", "Search")
	countFlag := flags.Int("c", 5, "Count by hundred steps")
	flags.Parse(os.Args[1:])

	cfg, err := ini.Load(*cfgfile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	consumerKey := cfg.Section("").Key("TWITTER_CONSUMER_KEY").String()
	consumerSecret := cfg.Section("").Key("TWITTER_CONSUMER_SECRET").String()
	accessToken := cfg.Section("").Key("TWITTER_accessToken").String()
	accessSecret := cfg.Section("").Key("TWITTER_accessSecret").String()

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		log.Fatal("Error: Consumer key/secret and Access token/secret required")
	}

	if *searchQuery == "" {
		log.Fatal("Error: Missing mandatory query")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	u := make(map[string]*User)
	var users []*User
	usersFile := *nodes
	if usersFile != "" {
		users = readUsers(usersFile)
		for _, ou := range users {
			u[ou.Name] = ou
		}
	}

	var links []Link

	var maxID = int64(0)

	for k := 0; k < *countFlag; k++ {
		// search tweets
		searchTweetParams := &twitter.SearchTweetParams{
			Query:      *searchQuery,
			TweetMode:  "compat",
			ResultType: "recent",
			MaxID:      maxID,
			Count:      100, // Max 100

		}
		search, a, b := client.Search.Tweets(searchTweetParams)
		//fmt.Printf("SEARCH TWEETS:\n%+v\n", search.Statuses)
		fmt.Printf("a: %+v\nb: %+v\n", a, b)
		fmt.Printf("SEARCH METADATA:\n%+v\n", search.Metadata)
		//os.Exit(0)
		//k := 0
		for _, s := range search.Statuses {
			//fmt.Printf("SEARCH TWEETS:\n%s\n", s.Text)
			tc, _ := s.CreatedAtTime()
			//fmt.Printf("SEARCH TWEETS: %s\n", tc.Format("2006-01-02-15h04"))
			maxID = s.ID
			fmt.Printf("*ID %d\n", maxID)
			fmt.Printf("User : %s (%d)\n", s.User.ScreenName, s.User.FollowersCount)
			if _, ok := u[s.User.ScreenName]; ok == false {
				newu := User{Name: s.User.ScreenName, Followers: s.User.FollowersCount}
				users = append(users, &newu)
				u[newu.Name] = &newu
			}
			for _, um := range s.Entities.UserMentions {
				fmt.Printf(" mention: %s\n", um.ScreenName)
				if _, ok := u[um.ScreenName]; ok == false {
					newu := User{Name: um.ScreenName, Followers: 0}
					users = append(users, &newu)
					u[newu.Name] = &newu
				}
				newl := Link{Source: s.User.ScreenName,
					Target:    um.ScreenName,
					CreatedAt: tc.Format("2006-01-02-15h04"),
					Key:       s.IDStr}
				links = append(links, newl)
			}
			if s.RetweetedStatus != nil {
				fmt.Printf(" Retweeted: %s (%d)\n", s.RetweetedStatus.User.ScreenName, s.RetweetedStatus.User.FollowersCount)
				u[s.RetweetedStatus.User.ScreenName].Followers = s.RetweetedStatus.User.FollowersCount
			}
		}
	}

	txt := *searchQuery
	txt = regexp.MustCompile(`[\(\)@# ]+`).ReplaceAllString(txt, "")

	linksFile := fmt.Sprintf("./data/links-%s-%s.json", t.Format("2006-01-02-15h04"), txt)
	if usersFile == "" {
		usersFile = fmt.Sprintf("./data/users-%s-%s.json", t.Format("2006-01-02-15h04"), txt)
	}

	prettyJSON, _ := json.MarshalIndent(users, "", "\t")
	//fmt.Printf("\n\nUsers :\n %s\n", string(prettyJSON))
	//		os.Exit(0)
	ioutil.WriteFile(usersFile, prettyJSON, 0644)

	prettyJSON, _ = json.MarshalIndent(links, "", "\t")
	//fmt.Printf("\n\nLinks :\n %s\n", string(prettyJSON))
	ioutil.WriteFile(linksFile, prettyJSON, 0644)

	fmt.Printf("%s\n", t.AddDate(0, 0, -4).Format("2006-01-02"))
}
