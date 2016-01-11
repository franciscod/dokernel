package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/docopt/docopt-go"
)

func printKernel(kernel *godo.Kernel) {
	fmt.Printf(`
    id: %d
    name: %s
    version: %s
`, kernel.ID, kernel.Name, kernel.Version)
}

func printDroplet(client *godo.Client, droplet *godo.Droplet) {
	fmt.Printf("Droplet %d (%s) current kernel:", droplet.ID, droplet.Name)
	kbase := strings.Replace(droplet.Kernel.Name, droplet.Kernel.Version, "", 1)
	printKernel(droplet.Kernel)
	fmt.Println()
	fmt.Printf("Listing (not older) kernels matching '%s'...\n", kbase)

	ks, _ := kernelList(client, droplet.ID)
	for _, k := range ks {
		if k.ID < droplet.Kernel.ID {
			continue
		}

		if !strings.Contains(k.Name, kbase) {
			continue
		}

		printKernel(&k)
	}
}

func printAllDroplets(client *godo.Client) {
	ds, _ := dropletList(client)

	for _, d := range ds {
		printDroplet(client, &d)
	}
}

func getIntArgument(arguments map[string]interface{}, name string) int {
	if arguments[name] == nil {
		return 0
	}

	val, err := strconv.Atoi(arguments[name].(string))

	if err != nil {
		fmt.Println(name + " should be an integer")
		os.Exit(1)
	}

	return val
}

func main() {

	usage := `
dokernel - a tool for changing digitalocean.com kernels

Usage:
  dokernel list [DROPLETID]
  dokernel set DROPLETID KERNELID
  dokernel poweron DROPLETID
  dokernel actions
`

	arguments, err := docopt.Parse(usage, nil, true, "dokernel 0.0.0", false)

	if err != nil {
		fmt.Println(err)
	}

	client := clientFromToken(readTokenFromFile())

	list := arguments["list"].(bool)
	set := arguments["set"].(bool)
	poweron := arguments["poweron"].(bool)
	actions := arguments["actions"].(bool)

	dropletID := getIntArgument(arguments, "DROPLETID")
	kernelID := getIntArgument(arguments, "KERNELID")

	if list {

		if dropletID == 0 {

			printAllDroplets(client)

		} else {

			droplet, _, err := client.Droplets.Get(dropletID)

			if err != nil {
				fmt.Println(err)
				return
			}

			printDroplet(client, droplet)

		}

	} else if set {

		action, _, err := client.DropletActions.ChangeKernel(dropletID, kernelID)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Kernel change requested. Action ID is", action.ID, "and status is", action.Status)
		fmt.Println("Keep in mind the kernel will only change when the droplet is powered off!")

	} else if poweron {

		action, _, err := client.DropletActions.PowerOn(dropletID)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Power on requested. Action ID is", action.ID, "and status is", action.Status)

	} else if actions {

		as, _ := actionList(client)

		for _, a := range as {
			fmt.Println(a.ID, a.Type, a.Status, a.StartedAt, a.CompletedAt)
		}
	}
}
