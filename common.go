package go_proxmox_backup_client

/*
#cgo LDFLAGS: -lproxmox_backup_qemu
#include <stdio.h>
#include <stdlib.h>
#include <proxmox-backup-qemu.h>

*/
import "C"

func GetVersion() string {
	version := C.proxmox_backup_qemu_version()
	return C.GoString(version)
}

