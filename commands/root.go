package commands

import (
	"github.com/paulonutor/hidrivecli/hidrive"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	authConfigPath string
	tokenSource    oauth2.TokenSource
	api            *hidrive.Client
)

var rootCmd = &cobra.Command{
	Use:                "hidrive",
	Short:              "A commandline interface to the HiDrive API",
	Long:               `This CLI can be used to control many of the HiDrive cloud storage features.`,
	PersistentPreRunE:  initClient,
	PersistentPostRunE: storeToken,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func loadToken() *oauth2.Token {
	data, err := ioutil.ReadFile(authConfigPath)

	if err != nil {
		return nil
	}

	var token oauth2.Token

	if err := json.Unmarshal(data, &token); err != nil {
		return nil
	}

	return &token
}

func storeToken(cmd *cobra.Command, args []string) (err error) {
	if _, err = os.Stat(authConfigPath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(authConfigPath), 0700)

		if err != nil {
			return
		}
	}

	token, err := tokenSource.Token()

	if err != nil {
		return
	}

	data, err := json.Marshal(token)

	if err != nil {
		return
	}

	err = ioutil.WriteFile(authConfigPath, data, 0600)

	return
}

func initClient(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "--------------------------------",
		ClientSecret: "--------------------------------",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.hidrive.strato.com/oauth2/authorize",
			TokenURL: "https://www.hidrive.strato.com/oauth2/token",
		},
		Scopes: []string{"admin,rw"},
	}

	homePath, err := homedir.Dir()

	if err != nil {
		return
	}

	authConfigPath = filepath.Join(homePath, ".config", "hidrivecli", "auth.json")

	tokenSource = conf.TokenSource(ctx, loadToken())

	if _, tokenErr := tokenSource.Token(); tokenErr != nil {
		fmt.Printf("Open the following URL and login to your HiDrive account: \n\n")
		fmt.Printf("%v\n\n", conf.AuthCodeURL("should-not-change"))
		fmt.Printf("After granting access copy the shown authorization code.\n")
		fmt.Printf("Enter the authorization code: ")

		var code string

		if _, err = fmt.Scan(&code); err != nil {
			return
		}

		var token *oauth2.Token

		token, err = conf.Exchange(ctx, code)

		if err != nil {
			return
		}

		tokenSource = conf.TokenSource(ctx, token)
	}

	httpClient := oauth2.NewClient(ctx, tokenSource)
	api = hidrive.NewClient(httpClient, "https://api.hidrive.strato.com/2.1/")

	return
}
