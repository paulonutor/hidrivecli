package commands

import (
	"github.com/paulonutor/hidrivecli/hidrive"

	"os"
	"text/tabwriter"

	"fmt"

	"github.com/spf13/cobra"
)

const userListFormat = "%s\t%s\t%s\t%s\t%v\t%v\t%s\n"

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users of the account",
	RunE:  listUsers,
}

func listUsers(cmd *cobra.Command, args []string) error {
	params := &hidrive.UserListParams{
		Scope: "all",
		Fields: []string{
			"alias",
			"account",
			"descr",
			"email",
			"home",
			"is_admin",
			"is_owner",
		},
	}
	users, err := api.Users.List(params)

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(
		w,
		userListFormat,
		"Account",
		"Alias",
		"Description",
		"Email",
		"Owner",
		"Admin",
		"Home",
	)

	for _, user := range users {
		fmt.Fprintf(
			w,
			userListFormat,
			user.AccountId,
			user.Alias,
			user.Description,
			user.Email,
			user.IsOwner,
			user.IsAdmin,
			user.Home,
		)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(userCmd)
}
