package go_proxmox_backup_client

/*
#cgo LDFLAGS: -lproxmox_backup_qemu
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <proxmox-backup-qemu.h>

ProxmoxBackupHandle *backup_new(const char *repo,
								const char *backup_id,
								uint64_t backup_time,
								const char *password,
								const char *fingerprint,
								const char *key_file,
								const char *key_password,
								char **error) {
	return proxmox_backup_new(repo, backup_id, backup_time, PROXMOX_BACKUP_DEFAULT_CHUNK_SIZE,
		password, key_file, key_password, NULL,
		false, key_file != NULL && strlen(key_file) > 0,
		fingerprint,
		error);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type ProxmoxBackup struct {
	handle *C.ProxmoxBackupHandle
}

func NewBackup(repo string, id string, backupTime uint64, password string, fingerprint string, keyFile string, keyPassword string) (*ProxmoxBackup, error) {
	cRepo := C.CString(repo)
	defer C.free(unsafe.Pointer(cRepo))

	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))

	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))

	cFingerprint := C.CString(fingerprint)
	defer C.free(unsafe.Pointer(cFingerprint))

	var cKeyFile *C.char
	if len(keyFile) > 0 {
		cKeyFile = C.CString(keyFile)
		defer C.free(unsafe.Pointer(cKeyFile))
	}
	var cKeyPassword *C.char
	if len(keyPassword) > 0 {
		cKeyPassword = C.CString(keyPassword)
		defer C.free(unsafe.Pointer(cKeyPassword))
	}

	var cErr *C.char

	Proxmox := new(ProxmoxBackup)

	Proxmox.handle = C.backup_new(cRepo, cId, C.ulong(backupTime), cPassword, cFingerprint, cKeyFile, cKeyPassword, &cErr)

	if Proxmox.handle == nil {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}

	e := C.proxmox_backup_connect(Proxmox.handle, &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}

	return Proxmox, nil
}

func (pbs *ProxmoxBackup) Close() {
	C.proxmox_backup_disconnect(pbs.handle)
}

func (pbs *ProxmoxBackup) Abort(reason string) {
	cReason := C.CString(reason)
	defer C.free(unsafe.Pointer(cReason))

	C.proxmox_backup_abort(pbs.handle, cReason)
}

func (pbs *ProxmoxBackup) AddConfig(name string, data []byte) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cErr *C.char

	e := C.proxmox_backup_add_config(pbs.handle, cName, (*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(len(data)), &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return errors.New(err)
	}
	return nil
}

func (pbs *ProxmoxBackup) Finish() error {
	var cErr *C.char

	e := C.proxmox_backup_finish(pbs.handle, &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return errors.New(err)
	}
	return nil
}

type BackupImage struct {
	proxmox *ProxmoxBackup
	dev     C.uint8_t
}

func (pbs *ProxmoxBackup) RegisterImage(name string, size uint64) (*BackupImage, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cErr *C.char

	e := C.proxmox_backup_register_image(pbs.handle, cName, C.ulong(size), C.bool(false), &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}
	return &BackupImage{proxmox: pbs, dev: C.uint8_t(uint8(e))}, nil
}

func (image *BackupImage) Write(data []byte, offset uint64) (int, error) {
	var cErr *C.char

	e := C.proxmox_backup_write_data(image.proxmox.handle, image.dev, (*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(offset), C.ulong(len(data)), &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return 0, errors.New(err)
	}
	return int(e), nil
}

func (image *BackupImage) Close() error {
	var cErr *C.char

	e := C.proxmox_backup_close_image(image.proxmox.handle, image.dev, &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return errors.New(err)
	}
	return nil
}
