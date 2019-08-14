package cmd

import (
	"fmt"
	"github.com/buckket/strahlemann/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get information about database entries",
	Run:   getStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getStatus(cmd *cobra.Command, args []string) {
	db, err := database.New(viper.GetString("DATABASE_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.CreateSchema()
	if err != nil {
		log.Fatal(err)
	}

	status, err := db.GetStatus()
	if err != nil {
		log.Fatal(err)
	}

	if status.AllEntries == status.DoneEntries {
		fmt.Printf("Current status:\n- %d/%d entries completed, all done!\n", status.DoneEntries, status.AllEntries)
	} else {
		fmt.Printf("Current status:\n"+
			"- %d/%d entries completed (~%d%%)\n"+
			"- Material sufficient for about %d more tweets\n"+
			"- Last Tweet: %d\n"+
			"- Current posstion: %d/%d\n"+
			"- Next Tweet: %q\n",
			status.DoneEntries, status.AllEntries, int(float64(status.DoneEntiresLength)/float64(status.AllEntriesLength)*100),
			int(float64(status.AllEntriesLength-status.DoneEntiresLength)/280),
			status.LastTweet,
			status.CurrentPosition, status.CurrentLenght,
			status.NextTweet)
	}
}
