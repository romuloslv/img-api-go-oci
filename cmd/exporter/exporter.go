package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"manager_oci/cmd/exporter/config"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
)

func UsageMsg() {
	fmt.Println(`flag needs an argument: -a
Usage of ./main:
  -a string
        Action to run: export/import img_id`)
}

func SetEnv(img_id string) {
	os.Setenv("OCI_EXPORT_IMAGEID", img_id)
	os.Setenv("OCI_EXPORT_OBJECTNAME", GetInfo(img_id))
}

func ReadFile() {
	_, err := os.ReadFile("~/.oci/config")
	if err != nil {
		log.Fatalln(err)
	}
}

func Rewrite() {
	file, err := os.OpenFile("~/.oci/config", os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	file.WriteAt([]byte("us-ashburn-1 "), 306)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func GetInfo(img_id string) string {
	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	req := core.GetImageRequest{
		ImageId: common.String(img_id),
	}

	resp, err := c.GetImage(context.Background(), req)
	helpers.FatalIfError(err)

	b, _ := json.Marshal(resp.Image.DisplayName)
	return strings.Split(string(b), "\"")[1]
}

func ExportImage() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	req := core.ExportImageRequest{
		ImageId: common.String(cfg.Export.ImageId),
		ExportImageDetails: core.ExportImageViaObjectStorageTupleDetails{
			BucketName:    common.String(cfg.Export.BucketName),
			NamespaceName: common.String(cfg.Export.Namespace),
			ObjectName:    common.String(cfg.Export.ObjectName),
			ExportFormat:  core.ExportImageDetailsExportFormatQcow2,
		},
	}

	resp, err := c.ExportImage(context.Background(), req)
	helpers.FatalIfError(err)
	fmt.Println("\n", resp.Image)
}

func ImportImage() {
	ReadFile()
	Rewrite()

	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	req := core.CreateImageRequest{
		CreateImageDetails: core.CreateImageDetails{
			ImageSourceDetails: core.ImageSourceViaObjectStorageTupleDetails{
				SourceImageType: core.ImageSourceDetailsSourceImageTypeQcow2,
				BucketName:      common.String(cfg.Export.BucketName),
				NamespaceName:   common.String(cfg.Export.Namespace),
				ObjectName:      common.String(cfg.Export.ObjectName),
			},
			CompartmentId: common.String(cfg.Export.CompartmentId),
			DisplayName:   common.String(cfg.Export.ObjectName)}}

	resp, err := c.CreateImage(context.Background(), req)
	helpers.FatalIfError(err)
	fmt.Println("\n", resp)
}
