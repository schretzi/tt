package tt

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

var digitCheck = regexp.MustCompile(`^[0-9]+$`)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientListCmd)
	clientCmd.AddCommand(clientCreateCmd)
	clientCmd.AddCommand(clientGetCmd)
	clientCmd.AddCommand(clientModifyCmd)
	clientCmd.AddCommand(clientActivateCmd)
	clientCmd.AddCommand(clientDeactivateCmd)

	clientListCmd.PersistentFlags().BoolVarP(&flagAll, "all", "a", false, "list all (including inactive) clients")

}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Management Client Items",
	Long:  `Management Client Items`,
}

var clientListCmd = &cobra.Command{
	Use:   "list",
	Short: "List (active) Clients by name(pattern)",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return clientListNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		var clientName string
		if len(args) > 0 {
			clientName = args[0]
		}

		cl := clientListNew()
		if clientName == "" {
			cl.listAll()

		} else {
			_ = cl.listByName(clientName)
		}

		clFiltered := clientListNew()
		if flagAll {
			clFiltered = cl
		} else {
			for _, c := range cl {
				if c.Active {
					clFiltered = append(clFiltered, c)
				}
			}
		}
		clientPrintTable(clFiltered)

	},
}

var clientCreateCmd = &cobra.Command{
	Use:   "create <client>",
	Short: "Create new Client",
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)
		clientName := args[0]

		if clientName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
		}
		client := clientNew()
		client.create(clientName)
		fmt.Println(msgCreateClient, client.ID, "/", client.Name)
	},
}

var clientModifyCmd = &cobra.Command{
	Use:   "modify <id> <client>",
	Short: "Change name of the Client (id or old name)",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return clientListNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 2, cmd)

		clientIdInput := args[0]

		if clientIdInput == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientId))
		}

		client := clientNew()

		if digitCheck.MatchString(clientIdInput) {
			clientId, _ := strconv.ParseUint(clientIdInput, 10, 32)
			client.getById(uint(clientId))
		} else {
			cl := clientListNew()
			rows := cl.listByName(clientIdInput)
			if rows != 1 {
				log.Fatalf(ErrorString, CharError, errors.New(errorUnambiguously))
			}
			client.getById(cl[0].ID)
		}

		clientName := args[1]

		if clientName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
		}

		client.modify(clientName)

		fmt.Println(msgModifiedClient, client.ID, "/", client.Name)
	},
}

var clientGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Display client by ID (uint) or Name (string)",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return clientListNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)

		clientIdInput := args[0]

		if clientIdInput == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientId))
		}

		client := clientNew()
		client.getByIdOrName(clientIdInput)

		cl := clientListNew()
		cl = append(cl, client)

		clientPrintTable(cl)
	},
}

func clientPrintTable(cl []Client) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Active"},
		},
	}

	for _, client := range cl {
		var activeText string
		if client.Active {
			activeText = msgActive
		} else {
			activeText = msgInactive
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", client.ID)},
			{Text: client.Name},
			{Text: activeText},
		}

		table.Body.Cells = append(table.Body.Cells, r)

	}
	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}

func clientListNames() []string {
	var listNames []string
	clientList := clientListNew()
	clientList.listAll()

	for _, client := range clientList {
		if client.Active {
			listNames = append(listNames, client.Name)
		}
	}
	return listNames
}

var clientActivateCmd = &cobra.Command{
	Use:   "activate <name>",
	Short: "Make",
	Long:  `Activate the Client`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return clientListNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)

		clientName := args[0]

		if clientName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
		}

		client := clientNew()
		cl := clientListNew()
		rows := cl.listByName(clientName)
		if rows != 1 {
			log.Fatalf(ErrorString, CharError, errors.New(errorUnambiguously))
		}
		client.getById(cl[0].ID)
		client.activate()

		fmt.Println(msgModifiedClient, client.ID, "/", client.Name)
	},
}

var clientDeactivateCmd = &cobra.Command{
	Use:   "deactivate <name>",
	Short: "Make",
	Long:  `Deactivate the Client`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return clientListNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)

		clientName := args[0]

		if clientName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
		}

		client := clientNew()
		cl := clientListNew()
		rows := cl.listByName(clientName)
		if rows != 1 {
			log.Fatalf(ErrorString, CharError, errors.New(errorUnambiguously))
		}
		client.getById(cl[0].ID)
		client.deactivate()

		fmt.Println(msgModifiedClient, client.ID, "/", client.Name)
	},
}
