package dokernel

import (
	"fmt"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/franciscod/godo"
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

func main() {

	var args struct {
		List      bool `arg:"-l"`
		Set       bool `arg:"-s"`
		PowerOn   bool `arg:"-p"`
		Actions   bool `arg:"-a"`
		DropletID int  `arg:"positional"`
		KernelID  int  `arg:"positional"`
	}
	arg.MustParse(&args)

	if !(args.Set || args.List || args.Actions || args.PowerOn) {
		args.List = true
	}

	if !args.Set && args.KernelID != 0 {
		fmt.Println("You provided a kernel id but didn't use --set, exiting. (see --help)")
		return
	}

	client := clientFromToken(readTokenFromFile())

	if args.List {
		if args.DropletID == 0 {
			printAllDroplets(client)
		} else {
			droplet, _, err := client.Droplets.Get(args.DropletID)
			if err != nil {
				fmt.Println(err)
				return
			}

			printDroplet(client, droplet)
		}
	} else if args.Set {
		if args.DropletID == 0 || args.KernelID == 0 {
			fmt.Println("Both droplet and kernel ID are required. (see --help)")
			return
		}

		action, _, err := client.DropletActions.ChangeKernel(args.DropletID, args.KernelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Kernel change requested. Action ID is", action.ID, "and status is", action.Status)
		fmt.Println("Keep in mind the kernel will only change when the droplet is powered off!")

	} else if args.PowerOn {
		if args.DropletID == 0 {
			fmt.Println("Droplet ID is required. (see --help)")
			return
		}

		action, _, err := client.DropletActions.PowerOn(args.DropletID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Power on requested. Action ID is", action.ID, "and status is", action.Status)
	} else if args.Actions {
		as, _ := actionList(client)

		for _, a := range as {
			fmt.Println(a.ID, a.Type, a.Status, a.StartedAt, a.CompletedAt)
		}
	}

}
