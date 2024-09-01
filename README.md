# go-proxmox-backup-client
Library for accessing the Proxmox Backup API from Go

## Example

```
package main

import (
	"fmt"
	"log"
	bps "github.com/elbandi/go-proxmox-backup-client"
	"time"
)

const (
	repo        = "admin@1.2.3.4:storage"
	password    = "changeme"
	fingerprint = "11:22:33:44:55:66:77..."
)

func backup(id string, backupTime time.Time) {
	t := uint64(backupTime.Unix())
	client, err := bps.NewBackup(repo, "", id, t, password, fingerprint, "", "", false)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	err = client.AddConfig("test", []byte("test2"))
	log.Println(err)

	image, err := client.RegisterImage("test", 16)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(image.WriteAt([]byte("1234567890123456"), 0))
		fmt.Println(image.Close())
	}
	err = client.Finish()
	log.Println(err)
}

func restore(id string, backupTime time.Time) {
	t := uint64(backupTime.Unix())
	client, err := bps.NewRestore(repo, "", "vm", id, t, password, fingerprint, "", "")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	image, err := client.OpenImage("test.img.fidx")
	if err != nil {
		log.Println(err)
	} else {
		data := make([]byte, 10)
		fmt.Println(image.ReadAt(data, 10))
		fmt.Println(string(data))
	}
}

func main() {
	fmt.Println(bps.GetVersion())
	t := time.Now()
	backup("123", t)
	restore("123", t)
}
```