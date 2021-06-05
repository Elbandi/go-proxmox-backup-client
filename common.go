package go_proxmox_backup_client

/*
#cgo LDFLAGS: -lproxmox_backup_qemu
#include <stdio.h>
#include <stdlib.h>
#include <proxmox-backup-qemu.h>

uint64_t get_default_chunk_size() {
	return PROXMOX_BACKUP_DEFAULT_CHUNK_SIZE;
}

*/
import "C"

func GetVersion() string {
	version := C.proxmox_backup_qemu_version()
	return C.GoString(version)
}

func GetDefaultChunkSize() uint64 {
	return uint64(C.get_default_chunk_size())
}
