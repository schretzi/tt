package tt

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.PersistentFlags().StringVarP(&paramClient, "client", "c", "", "list projects for client")
	projectCmd.RegisterFlagCompletionFunc("client", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return clientListNames(), cobra.ShellCompDirectiveNoFileComp
	})

	projectCmd.AddCommand(projectListCmd)
	projectListCmd.PersistentFlags().BoolVarP(&flagAll, "all", "a", false, "list all (including inactive) projects")

	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectGetCmd)

	/*	projectCmd.AddCommand(projectModifyCmd)
		projectCmd.AddCommand(projectActivateCmd)
		projectCmd.AddCommand(projectDeactivateCmd)
	*/

}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Management Project Items",
}

var projectListCmd = &cobra.Command{
	Use:   "list [--client|-c <Client>]",
	Short: "List (active) Project Items",
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string
		if len(args) > 0 {
			projectName = args[0]
		}

		pl := projectListNew()
		if projectName == "" && paramClient == ""{
			pl.listAll()
		} else if projectName == "" && paramClient != "" {
			client := clientNew()
			client.getByIdOrName(paramClient)
			pl.listByClient()
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

	},	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create -c|--client <client> <project>",
	Short: "Create new project for Client",
	Long:  `Create new project for Client`,
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)
		projectName := args[0]

		if projectName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorProjectName))
		}

		client := clientNew()
		if paramClient == "" {
			defaultClient := viper.GetString(cfgDefaultClient)
			if defaultClient == "" {
				cl := clientListNew()
				rows := cl.listByName(defaultClient)
				if rows != 1 {
					log.Fatalf(ErrorString, CharError, errors.New(errorUnambiguously))
				}
				client.getById(cl[0].ID)
			} else {
				client.getById(1)
			}
		} else {
			cl := clientListNew()
			rows := cl.listByName(paramClient)
			if rows != 1 {
				log.Fatalf(ErrorString, CharError, errors.New(errorUnambiguously))
			}
			client.getById(cl[0].ID)
		}

		project := projectNew()
		project.create(projectName, client.ID)
		fmt.Println(msgCreateProject, project.ID, "/", project.Name, "for client: ", project.Client.Name)
	},
}

/*
	var clientModifyCmd = &cobra.Command{
		Use:   "modify <id> <client>",
		Short: "Change name of the Client",
		Long:  `Change name of the Client`,
		Run: func(cmd *cobra.Command, args []string) {
			checkArguments(len(args), 2, cmd)

			clientIdInput := args[0]

			if clientIdInput == "" {
				log.Fatalf(ErrorString, CharError, errors.New(errorClientId))
			}

			u64, err := strconv.ParseUint(clientIdInput, 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			clientId := uint(u64)

			clientName := args[1]

			if clientName == "" {
				log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
			}

			client := clientNew()
			client.get(clientId)
			client.modify(clientName)

			fmt.Println(msgModifiedClient, client.ID, "/", client.Name)
		},
	}
*/
var projectGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Display project with ID",
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)

		projectIdInput := args[0]

		if projectIdInput == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorProjectId))
		}

		u64, err := strconv.ParseUint(projectIdInput, 10, 32)
		projectId := uint(u64)

		if err != nil {
			fmt.Println(err)
		}

		project := projectNew()
		project.get(projectId)

		projectList := projectListNew()
		projectList = append(projectList, project)

		projectPrintTable(projectList)
		/*
			fmt.Println("Project ID:", project.ID)
			fmt.Println("Name:", project.Name)
			fmt.Println("Client:", project.Client.Name)

			var activeText string
			if project.Active {
				activeText = msgActive
			} else {
				activeText = msgInactive
			}
			fmt.Println(msgActive+":", activeText)
		*/
	},
}

var projectFindCmd = &cobra.Command{
	Use:   "find <project>",
	Short: "Find project by Name",
	Long:  clientCmd.Short,
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)

		projectName := args[0]

		if projectName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
		}

		projectList := projectListNew()
		projectList.findByName(projectName)
		projectPrintTable(projectList)
	},
}

func projectList() {
	projectList := projectListNew()
	if paramClient != "" and {
		projectList.findForClient(1)
	}
	if flagAll {
		projectList.listAll()
	} else {
		projectList.listActive()
	}
	projectPrintTable(projectList)
}

/*
func clientListNames() []string {
	var listNames []string
	clientList := clientListNew()
	clientList.listActive()

	for _, client := range clientList {
		listNames = append(listNames, client.Name)
	}
	return listNames
}

var clientActivateCmd = &cobra.Command{
	Use:   "activate <name>",
	Short: "Make",
	Long:  `Activate the Client`,
	Run: func(cmd *cobra.Command, args []string) {
		checkArguments(len(args), 1, cmd)

		clientName := args[0]

		if clientName == "" {
			log.Fatalf(ErrorString, CharError, errors.New(errorClientName))
		}

		client := clientNew()
		client.findByName(clientName)
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
		client.findByName(clientName)
		client.deactivate()

		fmt.Println(msgModifiedClient, client.ID, "/", client.Name)
	},
}
*/

func projectPrintTable(projects []Project) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Client"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Active"},
		},
	}

	for _, project := range projects {
		var activeText string
		if project.Active {
			activeText = msgActive
		} else {
			activeText = msgInactive
		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", project.ID)},
			{Text: project.Client.Name},
			{Text: project.Name},
			{Text: activeText},
		}

		table.Body.Cells = append(table.Body.Cells, r)

	}
	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
}
