
package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type UserAttributes struct {
	Name       string
	Password	string
	Attributes map[string][]string
}

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user yamlfile",
	Short: "add a user to the IdP",
	Long: `Parses the user's data to create an entry in the 
	configuration file.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userData, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		var newUser UserAttributes
		err = yaml.Unmarshal(userData, &newUser)
		if err != nil {
			return err
		}
		// Hash the password if it is not already in bcrypt format
		pwd := newUser.Password
		if len(pwd) > 0 && !strings.HasPrefix(pwd, `$2a`) {
			pwd, err = hashPassword([]byte(pwd))
			if err != nil {
				return err
			}
			newUser.Password = pwd
		}
		// Get the existing users
		users := []UserAttributes{}
		if err = viper.UnmarshalKey("users", &users); err != nil {
			return err
		}
		found := false
		for i, u := range users {
			if u.Name == newUser.Name {
				users[i] = newUser
				found = true
				break
			}
		}
		if !found {
			users = append(users, newUser)
		}
		viper.Set("users", users)
		return viper.WriteConfig()
	},
}

func init() {
	AddCmd.AddCommand(userCmd)
}

