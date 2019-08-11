package cmd

import (
	"database/sql"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/buckket/strahlemann/database"
	"github.com/buckket/strahlemann/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strconv"
)

var tweetCmd = &cobra.Command{
	Use:   "tweet",
	Short: "Tweet the next piece of information",
	Run:   postTweet,
}

func init() {
	rootCmd.AddCommand(tweetCmd)
}

func postTweet(cmd *cobra.Command, args []string) {
	db, err := database.New(viper.GetString("DATABASE_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.CreateSchema()
	if err != nil {
		log.Fatal(err)
	}

	post, err := db.GetNextPost()
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No more data, we're done here!\n")
			return
		} else {
			log.Fatal(err)
		}
	}

	tweetText, pos := utils.ExtractTweet(post.Content[post.Position:])
	if pos == 0 {
		post.Complete = true
	} else {
		post.Position += pos
	}

	tapi := anaconda.NewTwitterApiWithCredentials(viper.GetString("TWITTER_ACCESS_TOKEN"),
		viper.GetString("TWITTER_ACCESS_TOKEN_SECRET"),
		viper.GetString("TWITTER_CONSUMER_KEY"),
		viper.GetString("TWITTER_CONSUMER_SECRET"))
	_, err = tapi.GetSelf(url.Values{})
	if err != nil {
		log.Fatal(err)
	}

	v := url.Values{}
	if post.LastTweet > 0 {
		v.Add("in_reply_to_status_id", strconv.FormatInt(post.LastTweet, 10))
		v.Add("auto_populate_reply_metadata", "true")
	}
	tweet, err := tapi.PostTweet(tweetText, v)
	if err != nil {
		log.Fatal(err)
	}

	post.LastTweet = tweet.Id
	fmt.Printf("Posted: %q\nURL: %s\n", tweetText, utils.GenerateTweetURL(viper.GetString("TWITTER_USERNAME"), tweet.Id))

	err = db.UpdatePost(post)
	if err != nil {
		log.Fatal(err)
	}
}
