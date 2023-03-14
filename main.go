package main

import (
	"flag"
	"fmt"
	"log"
	"manager_oci/cmd/exporter"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

func main() {
	_, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		log.Fatalln(err)
	}

	var action string
	flag.StringVar(&action, "a", "", "Action to run: export/import img_id")
	flag.Parse()

	now := time.Now()
	fmt.Println("Time:", now.Format(time.ANSIC))
	fmt.Println("Action:", action)
	fmt.Println("Args:", flag.Args())

	if action == "export" {
		exporter.SetEnv(flag.Arg(0))
		exporter.ExportImage()
	} else if action == "import" {
		exporter.SetEnv(flag.Arg(0))
		exporter.ImportImage()
	} else {
		exporter.UsageMsg()
		os.Exit(1)
	}
}
