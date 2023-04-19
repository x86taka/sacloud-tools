package main

import (
	"flag"
	"fmt"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	sacloudArchive "github.com/sacloud/iaas-service-go/archive"
	sacloudDisk "github.com/sacloud/iaas-service-go/disk"
	sacloudServer "github.com/sacloud/iaas-service-go/server"
	sacloudLocalSwitch "github.com/sacloud/iaas-service-go/swytch"
	darchive "github.com/x86taka/sacloud-tools/makehcl/data/archive"
	rdisk "github.com/x86taka/sacloud-tools/makehcl/resource/disk"
	rswytch "github.com/x86taka/sacloud-tools/makehcl/resource/swytch"
	"github.com/x86taka/sacloud-tools/makehcl/resource/utils"
	"github.com/x86taka/sacloud-tools/makehcl/resource/vm"
	"os"
)

var filePrefix = "output/"
var fileSuffix = ".tf"
var generatedID = map[string]int{} //生成済みかどうかのcheck用

// 全てのリソースのタグにつける tagたち
var requireTags = types.Tags{}

func main() {

	flag.Parse()
	args := flag.Args()

	// APIキー
	token := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	secret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")
	zone := os.Getenv("SAKURACLOUD_ZONE")

	// クライアントの作成
	client := iaas.NewClient(token, secret)

	// Service
	archiveService := sacloudArchive.New(client)
	diskService := sacloudDisk.New(client)
	serverService := sacloudServer.New(client)
	switchService := sacloudLocalSwitch.New(client)

	// Output
	dataOutput := ""
	diskOutput := ""
	serverOutput := ""
	swytchOutput := ""

	// FindServers
	serverFindRequest := sacloudServer.FindRequest{
		Zone: zone,
		Tags: args,
	}
	servers, err := serverService.Find(&serverFindRequest)
	if err != nil {
		panic(err)
	}

	for _, v := range servers {
		//Disks
		diskIDs := map[int]string{}
		for i := 0; i < len(v.GetDisks()); i++ {
			connectedDisk := v.GetDisks()[i]

			// Read sakuracloud_disk.
			readRequest := &sacloudDisk.ReadRequest{
				Zone: zone,
				ID:   connectedDisk.GetID(),
			}
			disk, err := diskService.Read(readRequest)
			if err != nil {
				panic(err)
			}

			if !disk.SourceArchiveID.IsEmpty() {
				// Read sakuracloud_archive
				readArchiveRequest := &sacloudArchive.ReadRequest{
					Zone: zone,
					ID:   disk.SourceArchiveID,
				}

				archive, err := archiveService.Read(readArchiveRequest)
				if err != nil {
					panic(err)
				}

				if !IsGenerated(disk.GetSourceArchiveID().String()) {
					// Create data.archive HCL
					dataArchive := &darchive.DataArchiveHCL{
						Name: archive.GetName(),
					}
					dataOutput += dataArchive.OutputHCL()
				}

				// Create resource.Disk HCL
				diskResource := &rdisk.DiskHCL{
					Name:            disk.GetName(),
					Size:            int64(disk.GetSizeGB()),
					Tags:            append(disk.GetTags(), requireTags...),
					SourceArchiveID: "data.sakuracloud_archive." + utils.FormatHCL(archive.GetName()) + ".id",
				}
				diskOutput += diskResource.OutputHCL()
			}
			diskIDs[i] = "sakuracloud_disk." + utils.FormatHCL(disk.GetName()) + ".id"
		}

		// Nics
		nicIDs := map[int]string{}
		for i := 0; i < len(v.GetInterfaces()); i++ {
			nic := v.GetInterfaces()[i]
			if !nic.GetSwitchID().IsEmpty() && nic.GetSwitchName() != "スイッチ" {
				readSwitchRequest := &sacloudLocalSwitch.ReadRequest{
					Zone: zone,
					ID:   nic.GetSwitchID(),
				}
				sw, err := switchService.Read(readSwitchRequest)
				if err != nil {
					panic(err)
				}
				if !IsGenerated(nic.SwitchID.String()) {
					swHCL := rswytch.LocalSwitchHCL{
						Name: sw.GetName(),
						Tags: append(sw.GetTags(), requireTags...),
					}
					swytchOutput += swHCL.OutputHCL()
				}
				nicIDs[i] = fmt.Sprintf("sakuracloud_switch.%s.id", utils.FormatHCL(sw.GetName()))
			} else {
				nicIDs[i] = "\"shared\""
			}

		}
		vm := &vm.VM{
			Name:  v.GetName(),
			Tags:  append(v.GetTags(), requireTags...),
			Cpus:  int64(v.GetCPU()),
			Mem:   int64(v.GetMemoryGB()),
			Disks: diskIDs,
			Nics:  nicIDs,
		}
		serverOutput += vm.OutputHCL()
	}
	// ファイル書き込み
	file, err := os.Create(filePrefix + "data" + fileSuffix)
	if err != nil {
		panic(err)
	}
	file.WriteString(dataOutput)
	file.Close()

	file, err = os.Create(filePrefix + "disk" + fileSuffix)
	if err != nil {
		panic(err)
	}
	file.WriteString(diskOutput)
	file.Close()

	file, err = os.Create(filePrefix + "server" + fileSuffix)
	if err != nil {
		panic(err)
	}
	file.WriteString(serverOutput)
	file.Close()

	file, err = os.Create(filePrefix + "switch" + fileSuffix)
	if err != nil {
		panic(err)
	}
	file.WriteString(swytchOutput)
	file.Close()
}

func IsGenerated(id string) bool {
	if _, ok := generatedID[id]; ok {
		return true
	}
	generatedID[id] = 1
	return false
}
