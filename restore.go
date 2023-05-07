package go_proxmox_backup_client

/*
#cgo LDFLAGS: -lproxmox_backup_qemu
#include <stdio.h>
#include <stdlib.h>
#include <proxmox-backup-qemu.h>

ProxmoxRestoreHandle *restore_new(const char *repo,
								  const char *snapshot,
								  const char *password,
								  const char *fingerprint,
								  const char *key_file,
								  const char *key_password,
								  char **error) {
	return proxmox_restore_new(repo, snapshot,
		password, key_file, key_password,
		fingerprint,
		error);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type ProxmoxRestore struct {
	handle *C.ProxmoxRestoreHandle
}

func NewRestore(repo string, btype string, id string, backupTime uint64, password string, fingerprint string, keyFile string, keyPassword string) (*ProxmoxRestore, error) {
	cRepo := C.CString(repo)
	defer C.free(unsafe.Pointer(cRepo))

	cType := C.CString(btype)
	defer C.free(unsafe.Pointer(cType))

	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))

	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))

	cFingerprint := C.CString(fingerprint)
	defer C.free(unsafe.Pointer(cFingerprint))

	var cKeyFile *C.char
	if len(keyFile) > 0 {
		cKeyFile := C.CString(keyFile)
		defer C.free(unsafe.Pointer(cKeyFile))
	}
	var cKeyPassword *C.char
	if len(keyPassword) > 0 {
		cKeyPassword := C.CString(keyPassword)
		defer C.free(unsafe.Pointer(cKeyPassword))
	}

	var cErr *C.char

	snapshot := C.proxmox_backup_snapshot_string(cType, cId, C.long(int64(backupTime)), &cErr)
	if snapshot == nil {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}
	defer C.free(unsafe.Pointer(snapshot))

	Proxmox := new(ProxmoxRestore)

	Proxmox.handle = C.restore_new(cRepo, snapshot, cPassword, cFingerprint, cKeyFile, cKeyPassword, &cErr)

	if Proxmox.handle == nil {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}

	e := C.proxmox_restore_connect(Proxmox.handle, &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}

	return Proxmox, nil
}

func (pbs *ProxmoxRestore) Close() {
	C.proxmox_restore_disconnect(pbs.handle)
}

type RestoreImage struct {
	proxmox *ProxmoxRestore
	dev     C.uint8_t
}

func (pbs *ProxmoxRestore) OpenImage(name string) (*RestoreImage, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cErr *C.char

	e := C.proxmox_restore_open_image(pbs.handle, cName, &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return nil, errors.New(err)
	}
	return &RestoreImage{proxmox: pbs, dev: C.uint8_t(uint8(e))}, nil
}

func (image *RestoreImage) Size() (uint64, error) {
	var cErr *C.char

	e := C.proxmox_restore_get_image_length(image.proxmox.handle, image.dev, &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return 0, errors.New(err)
	}
	return uint64(e), nil
}

func (image *RestoreImage) ReadAt(p []byte, off int64) (int, error) {
	var cErr *C.char

	e := C.proxmox_restore_read_image_at(image.proxmox.handle, image.dev, (*C.uchar)(unsafe.Pointer(&p[0])), C.ulong(off), C.ulong(len(p)), &cErr)
	if e < 0 {
		err := C.GoString(cErr)
		C.proxmox_backup_free_error(cErr)
		return 0, errors.New(err)
	}
	return int(e), nil
}

