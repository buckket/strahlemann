package cmd

import (
	"github.com/buckket/strahlemann/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new blog post",
	Run:   addPost,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addPost(cmd *cobra.Command, args []string) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	post := database.Post{
		Content: string(b),
	}

	db, err := database.New(viper.GetString("DATABASE_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.CreateSchema()
	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertPost(&post)
	if err != nil {
		log.Fatal(err)
	}
}
