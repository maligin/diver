package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {
	ucpNodesGet.Flags().StringVar(&id, "id", "", "ID of the Docker Node")

	ucpNodes.AddCommand(ucpNodesList)
	ucpNodes.AddCommand(ucpNodesGet)

	// Add nodes to UCP root commands
	UCPRoot.AddCommand(ucpNodes)

}

var ucpNodes = &cobra.Command{
	Use:   "nodes",
	Short: "Interact with Docker Nodes",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpNodesList = &cobra.Command{
	Use:   "list",
	Short: "Retrieve all Docker Nodes",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		nodes, err := client.ListAllNodes()
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Debugf("Found [%d] nodes", len(nodes))
		if len(nodes) == 0 {
			log.Fatalf("No Nodes found")
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Name\tID\tRole\tVersion\tPlatform")
		for i := range nodes {
			// Combine OS/Arch
			platform := nodes[i].Description.Platform.OS + "/" + nodes[i].Description.Platform.Architecture
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", nodes[i].Description.Hostname, nodes[i].ID, nodes[i].Spec.Role, nodes[i].Description.Engine.EngineVersion, platform)
		}
		w.Flush()
	},
}

var ucpNodesGet = &cobra.Command{
	Use:   "get",
	Short: "Get information about a particular Docker Node",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		node, err := client.GetNode(id)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Debugf("Retrieved information about [%s]", node.Description.Hostname)
		for k, v := range node.Spec.Labels {
			fmt.Printf("%s / %s\n", k, v)
		}
	},
}
