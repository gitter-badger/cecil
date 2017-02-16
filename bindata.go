// Code generated by go-bindata.
// sources:
// core/email-templates/account-verification-notification.txt
// core/email-templates/expiring-lease.txt
// core/email-templates/lease-approved.txt
// core/email-templates/lease-extended.txt
// core/email-templates/lease-resource-terminated.txt
// core/email-templates/misconfiguration-notice.txt
// core/email-templates/new-lease-no-owner-tag.txt
// core/email-templates/new-lease-owner-tag-not-whitelisted.txt
// core/email-templates/new-lease-valid-owner-tag-needs-approval.txt
// core/email-templates/new-lease-valid-owner-tag-no-approval-needed.txt
// core/email-templates/region-successfully-setup.txt
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _coreEmailTemplatesAccountVerificationNotificationTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x90\x5f\x4b\xc3\x30\x14\xc5\xdf\xf3\x29\xae\x79\xd9\xcb\xe6\xde\x47\x29\x48\x11\xa6\x88\x1b\xae\x82\x0f\x42\xb9\x4d\xef\xd6\xb8\x2e\xa9\xc9\xad\x5a\x4a\xbe\xbb\xad\xdd\xfc\x83\x20\xf8\x12\xb8\x27\xe7\xfc\x72\x72\x97\xd4\x42\xd7\x9d\xa3\x52\xb6\x31\x9c\x19\x3c\x50\x08\x53\x60\x0b\xca\x11\x32\x01\x1a\xb8\x58\x5f\xf5\xc2\x9e\xcc\x54\x78\x32\x05\x20\xac\x57\x9b\x14\x1c\x3d\x37\xe4\x79\xf0\x46\x79\xdc\x43\x5e\xc8\xe9\xad\x56\xc8\xda\x9a\x8c\xd1\xed\x88\xb3\xc6\x55\x21\x44\xf3\x3c\x86\x57\xcd\x25\x70\xa9\x3d\x5c\x6f\x56\xb7\x50\x63\x5b\x59\x2c\x16\x22\xca\x5d\x2c\x3e\xce\x4e\xfe\x24\x0c\x4f\xca\x85\xfc\x45\x1e\xf4\x10\x64\xf8\x4a\x8a\xe4\xfe\xee\x06\x2e\xdf\xf0\x50\x57\x74\x44\x8e\xb7\xaa\x2f\x00\x8f\xe3\x30\x5b\x82\x4c\xac\x61\x32\x3c\x4b\xdb\x9a\x16\x80\x75\x5d\x1d\xb1\xf3\x27\x6f\x8d\xfc\xb4\x3e\x8c\x7f\x3c\x8d\x05\x4c\xfe\xdb\x6e\x72\x0a\xff\xb5\x19\xf1\xad\x6b\x5a\xa2\xd9\x7b\xd8\x5a\x07\x8d\xd7\x66\x07\x09\x29\x5d\x9d\x89\xf7\x00\x00\x00\xff\xff\xcb\xcf\x60\x27\xa3\x01\x00\x00")

func coreEmailTemplatesAccountVerificationNotificationTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesAccountVerificationNotificationTxt,
		"core/email-templates/account-verification-notification.txt",
	)
}

func coreEmailTemplatesAccountVerificationNotificationTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesAccountVerificationNotificationTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/account-verification-notification.txt", size: 419, mode: os.FileMode(436), modTime: time.Unix(1486825572, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesExpiringLeaseTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x92\xc1\xee\xa2\x30\x10\x87\xef\x3c\xc5\x2c\x97\xfd\x9b\x18\xbc\x1b\xe5\x62\x36\xd9\x3d\xef\xde\x9b\x02\x83\x4c\x2c\xad\x29\x25\x4a\x48\xdf\x7d\x87\x52\x74\x5d\x35\xf1\x54\x68\xbf\xe9\x7c\xfd\xb5\xe3\x48\x35\x68\xe3\x20\xeb\x9c\x2c\x4f\x42\xcb\x16\xc1\xfb\x04\xe0\x27\x0e\x30\x8e\x99\xb9\x68\xb4\x02\x5b\x49\xca\xfb\x35\x90\x66\x4e\x97\x08\xbb\x22\xe7\xd5\xe5\x57\x50\xe5\xfd\x6e\x53\xe4\x5c\x08\xf0\x65\x6a\x70\xc3\xf9\x09\x9a\xe6\x66\x6c\x1d\x38\xa3\x23\x61\xb1\x33\xbd\x65\xc2\xe2\x91\x8c\x9e\x99\x15\x50\x07\x78\x3d\x93\x25\x7d\xcc\x92\x71\x44\xcd\x4d\x92\x24\x28\x7f\xa8\xeb\x1a\x84\xda\x28\x65\x2e\xbc\x07\x84\x9a\xc7\x5d\x81\x0d\x6c\xbe\x0c\x3c\xfe\x0e\xcc\xb4\xef\x36\xca\xdd\x3b\xcd\x5e\xb1\x60\x06\xa9\x7a\xc4\x96\x1c\x02\x74\x53\x7e\xea\xf2\xcb\xc1\x85\x94\x9a\x45\xf0\x1e\x84\x43\xdb\x92\x96\x8e\x43\x10\x8e\x96\x8e\xf0\xc5\x4b\x0a\x65\x87\xa2\xea\x6d\x58\xf5\x1e\x64\xcd\x34\x90\xfb\xde\x41\x69\x31\xcc\xae\xb2\x37\xbd\x64\x64\xb0\x9a\x7a\xfd\x7b\x27\x71\x5a\x48\xe7\xfd\xab\xea\x3f\xd1\x08\x81\xda\x16\x2b\xe2\x2f\x35\x6c\xff\xe3\x78\x90\xd0\x58\xac\xf7\xe9\xcd\x74\x39\x0a\x8a\xde\xf2\x5d\xa4\xe0\xa4\x3d\xa2\xdb\xa7\xa2\x50\x52\x9f\xd2\xfc\xa0\x88\x03\x6c\x90\xcf\xef\xcc\x74\xfe\x5b\x45\x08\x70\x23\xf3\x17\x3a\x3f\xae\x8e\x33\x85\x62\x88\x81\x61\xf8\x17\xc5\x30\x27\xf5\x89\x58\x2c\xf9\xd4\x6a\xc6\xdf\x29\x71\x40\x0d\xd7\x75\xfc\xca\x2c\xf4\xdd\xf4\xca\x0e\x58\x92\xfa\x96\xfc\x0d\x00\x00\xff\xff\xa6\x5f\x75\x05\x5a\x03\x00\x00")

func coreEmailTemplatesExpiringLeaseTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesExpiringLeaseTxt,
		"core/email-templates/expiring-lease.txt",
	)
}

func coreEmailTemplatesExpiringLeaseTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesExpiringLeaseTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/expiring-lease.txt", size: 858, mode: os.FileMode(436), modTime: time.Unix(1487264642, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesLeaseApprovedTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x90\xc1\x6e\xf2\x30\x10\x84\xef\x79\x8a\xf9\x4f\x3f\x48\x88\xde\x11\xe2\xd2\x4b\xef\xed\x3d\x32\xc9\x06\x2c\xc2\x3a\x5a\x3b\xa5\xc8\xf2\xbb\x77\x63\x13\xda\xaa\xbd\xb4\x37\x7b\xf5\xed\xec\xcc\xc4\x68\x3b\xb0\x0b\x58\xfb\x60\x9a\x53\xcd\xe6\x4c\x48\xa9\x02\x9e\xe8\x8a\x18\xd7\xee\xc2\x24\x35\x9d\x8d\xed\x53\x5a\x21\x1c\x09\x3d\x19\x4f\x70\x1d\x2c\xeb\x12\x37\x84\xed\x7e\xa7\xe8\xfc\xad\x6d\x9b\xd2\xf6\x61\xbf\x53\x15\x60\xa1\x60\xb8\x0e\xdf\xa0\x69\x56\xb0\x55\xe6\x1c\xdf\x08\x21\xef\x46\x51\x42\xe8\x60\x1d\x17\x66\x89\xa3\xf1\xd8\x13\x31\xcc\x30\x88\x7b\xa5\x76\x5d\xc5\x48\xac\xa7\xaa\x2a\xa7\xf8\x4b\x82\xe9\xd3\xb9\xbe\x77\x17\xcb\x07\x64\x01\x35\xfc\x5b\x27\x9b\x2a\x27\x80\xae\xc9\xee\xf3\x53\xdf\xcf\x59\x73\x32\xb5\xb9\xa9\x7e\xd8\x2c\x82\xb7\xa5\x02\xda\xf6\x2b\x36\x57\x99\xa1\x7b\xde\xfb\xa5\xf9\xca\x8b\xe6\x68\x46\x11\xe2\x00\x7a\x1b\xac\x98\xa0\x86\x61\x7d\x66\x26\xb9\x3c\x25\x5f\x9b\x50\x04\xb1\xd0\x61\xee\xa1\x6e\xc7\x82\xa7\x04\xd3\x05\x12\xd8\xf0\xdf\xa3\x11\xca\xd3\xe5\x8f\xe7\x0c\x9f\xbc\x36\x27\x18\xfd\xd4\xdc\x23\x35\xb6\xff\x57\xbd\x07\x00\x00\xff\xff\x1d\x62\xde\x8c\x4e\x02\x00\x00")

func coreEmailTemplatesLeaseApprovedTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesLeaseApprovedTxt,
		"core/email-templates/lease-approved.txt",
	)
}

func coreEmailTemplatesLeaseApprovedTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesLeaseApprovedTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/lease-approved.txt", size: 590, mode: os.FileMode(436), modTime: time.Unix(1487264646, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesLeaseExtendedTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x90\xbf\x6e\x32\x31\x10\xc4\xfb\x7b\x8a\xf9\xaa\x0f\x24\x44\x7a\x84\x68\xd2\xa4\x4f\xfa\x93\xb9\xdb\x83\x15\xc7\x1a\xd9\x3e\x11\x64\xf9\xdd\xb3\xb6\x8f\xfc\x53\x9a\xa4\x5b\xaf\x7e\x33\x9e\xd9\x18\x79\x80\xd8\x80\xb5\x0f\xa6\x3b\xb5\x62\xce\x84\x94\x1a\xe0\x89\x6e\x88\x71\x6d\xaf\x42\xae\xa5\xb3\xe1\x31\xa5\x15\xc2\x91\x30\x92\xf1\x04\x3b\x80\x45\x45\xd2\x11\xae\x1c\x8e\xe0\x1e\xdb\xfd\x4e\x25\xf7\x75\xcb\x7d\x4a\xdb\x87\xfd\x4e\xdd\x80\x85\x0a\xc2\xed\x42\xdf\xa1\xbc\xab\xd8\xaa\x70\x56\x66\xc2\x91\xb7\x93\x53\xc2\xd1\x81\xad\x54\x66\x89\xa3\xf1\xd8\x13\x09\xe8\x35\x90\xf4\xd4\xaf\x9b\x18\x75\xd0\xd0\x4d\x69\xf3\x97\x26\xf9\x31\xd8\x71\xb4\x57\x96\x03\x8a\x81\x06\xfe\x6d\x92\x4d\x53\x1a\x40\x65\x6e\xf7\x79\xd4\xf9\xb9\x78\xe6\x50\x9b\xd9\xf5\x23\x66\x35\x9c\x45\x15\xe4\xfe\x2b\x76\x3f\x65\x81\xde\xfb\x66\xe3\x59\x76\xff\xe7\x45\x9b\x74\x93\x73\x24\x41\x73\x5d\xd8\x99\xa0\x91\xc1\xbe\x30\xd9\xb0\x6c\xc9\xb7\x26\x54\x4b\x2c\x74\x59\x2e\xd1\xf6\x53\xc5\x53\x82\x19\x02\x39\x70\xf8\xef\xd1\x39\x2a\xdb\xe5\x8f\xdf\x19\x39\x79\xbd\x9d\xc3\xe4\xf3\xed\x1e\xa9\xe3\xf1\x5f\xf3\x16\x00\x00\xff\xff\x4f\xb1\xac\xd0\x58\x02\x00\x00")

func coreEmailTemplatesLeaseExtendedTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesLeaseExtendedTxt,
		"core/email-templates/lease-extended.txt",
	)
}

func coreEmailTemplatesLeaseExtendedTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesLeaseExtendedTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/lease-extended.txt", size: 600, mode: os.FileMode(436), modTime: time.Unix(1487264651, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesLeaseResourceTerminatedTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xcc\x91\x31\x6f\xfa\x30\x10\xc5\x77\x7f\x8a\xf7\x9f\xfe\x20\x21\xba\x57\x88\xa5\x4b\xf7\x76\x8f\x4c\x72\x21\x27\xc2\x19\xd9\x46\x34\xb2\xfc\xdd\xeb\xd8\x86\xa8\x55\x87\x8e\xcd\xe4\x5c\x7e\xef\x5e\xde\x73\x08\xdc\x43\x8c\xc7\xd6\x79\xdd\x9e\x1a\xd1\x67\x42\x8c\x0a\x78\xa5\x09\x21\x6c\xcd\x4d\xc8\x36\x74\xd6\x3c\xc6\xb8\x01\x4b\xe2\xa4\x25\xdc\xd8\x0f\xe0\x0e\xbb\xc3\x3e\x51\xf7\x71\xc3\x5d\x8c\xbb\xa7\xc3\x3e\x2d\x98\x9f\x95\xe9\xe1\xa7\x0b\x7d\xc7\xe6\x59\x01\x37\x95\x34\x52\x19\x4b\xce\x5c\x6d\x62\x2c\x1d\xd9\x48\xa5\x40\x1f\x17\xb6\xd3\x82\xe5\x77\x72\x8d\xf6\x85\x58\x63\xd0\x0e\x07\x22\x81\x27\x7b\x66\xd1\x9e\x3a\x68\x5f\xd7\x17\xd1\xf2\xe5\xa1\xc3\x2a\xcd\x47\xd2\x8e\x9a\xee\x6a\xb5\xcf\x96\xd0\x7d\x42\xc1\xfe\xbf\x43\x6b\x29\x4f\xd7\x2a\x04\x92\x14\x4f\xa9\x5c\xda\x2f\x0b\xf3\x03\xa1\x37\xe3\x68\x6e\x2c\x47\x64\x0d\x56\x7f\x20\x8c\x2a\x8b\xec\x7e\x39\xa4\xd3\x5b\xfe\xbf\x39\xd3\x73\x75\x59\x52\x16\x8b\x2a\x29\x20\x77\x5f\xb1\xfb\xed\x67\xe8\x51\xd7\x8f\x4e\xc0\xfb\xa0\xe5\xe4\x52\x39\x16\x57\x37\x97\xf3\x42\x2d\x8f\xff\xd4\x67\x00\x00\x00\xff\xff\xf0\xb8\x3d\x4f\x93\x02\x00\x00")

func coreEmailTemplatesLeaseResourceTerminatedTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesLeaseResourceTerminatedTxt,
		"core/email-templates/lease-resource-terminated.txt",
	)
}

func coreEmailTemplatesLeaseResourceTerminatedTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesLeaseResourceTerminatedTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/lease-resource-terminated.txt", size: 659, mode: os.FileMode(436), modTime: time.Unix(1487264655, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesMisconfigurationNoticeTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xf2\x48\xad\x54\xc8\x2c\x51\x48\x2c\x28\x48\x4d\x2c\x2a\x56\x28\xc9\x48\x2c\x51\x70\x4e\x4d\xce\xcc\x51\xc8\x2c\x56\xc8\xcd\x2c\xd6\x4d\xce\xcf\x4b\xcb\x4c\x2f\x2d\x4a\x4d\xd1\xe3\xb2\x49\x2a\xb2\x83\x10\xae\x45\x45\xf9\x45\x56\x10\x76\x75\xb5\x5e\x6a\x51\x51\x6d\x2d\x17\x20\x00\x00\xff\xff\xbb\x11\xc5\x62\x4c\x00\x00\x00")

func coreEmailTemplatesMisconfigurationNoticeTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesMisconfigurationNoticeTxt,
		"core/email-templates/misconfiguration-notice.txt",
	)
}

func coreEmailTemplatesMisconfigurationNoticeTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesMisconfigurationNoticeTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/misconfiguration-notice.txt", size: 76, mode: os.FileMode(436), modTime: time.Unix(1486825687, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesNewLeaseNoOwnerTagTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x52\xc1\x8a\xdb\x30\x10\xbd\xfb\x2b\xa6\xb9\x34\x81\x90\xbd\x2f\x59\x43\xd9\x4b\xf7\xd4\x43\xf7\x6e\x14\x6b\x6c\x0f\x2b\x4b\x41\x92\x63\x8c\xd1\xbf\x77\x24\xd9\x6e\xb3\x94\x12\xe8\xc1\x60\x66\xde\x9b\x79\xef\x8d\xe6\x99\x1a\xd0\xc6\xc3\xc9\x79\x51\x7f\x54\x5a\xf4\x08\x21\x14\x00\xdf\x71\x82\x79\x3e\x99\x51\xa3\xad\xb0\x17\xa4\x42\x38\x82\x33\x3d\x1a\x8d\x50\x5b\x14\x1e\x25\x08\xd0\x38\x02\x69\x66\xeb\x1a\x99\xb6\x27\x09\xe7\x4b\xc9\xcc\xb5\x58\x91\x0c\xe1\xfc\x74\x29\x8f\x60\x1a\xf0\xd3\x15\x3f\x03\x62\x6d\x81\xf0\x04\xa3\x97\xbe\x45\x67\x06\xcb\x7d\x8b\x2d\x19\x9d\x11\x87\x13\x77\x6d\x19\xbf\x62\x9e\x51\xf3\xec\xa2\x48\x2e\xfe\xcb\x41\xe2\x46\xf9\x8f\x6f\x67\xf4\xcf\xc8\x82\xb8\xf1\x79\x61\xfd\xd6\x90\x09\x09\xb9\x02\x49\xde\xc3\xd6\x60\xee\xcd\xbc\x79\x90\x06\x5d\xba\x4a\x27\x6e\xc8\x0a\x6f\x42\x71\xac\xaf\x58\x93\xfa\x11\xdd\x80\x17\x6d\x74\x02\x23\x77\x9d\xa3\x56\xb3\x11\xf2\xe0\x0d\x4c\x66\x80\xbd\xef\xb8\x2e\x7b\xd2\x87\x53\x51\xa4\xe9\x59\xf1\x5b\xbe\xb5\xb8\x5e\xad\xb9\xa1\x3c\x46\xce\x48\x4a\xc1\x05\xc1\xa3\x65\x42\x8e\xc4\x2f\x32\xd7\x1a\xfb\xaf\x3c\xad\x9e\x60\xcf\x2d\x85\xc2\x61\x25\x07\x9b\xba\x21\x80\x68\x18\xcd\x03\xbf\xba\x1c\x2d\x57\x3f\x6d\x7f\x5f\x37\x00\xf5\x3d\x4a\xe2\x3f\x35\x3d\xff\x81\x38\x0b\xe8\x2c\x36\x2f\xbb\x6d\xfe\x26\xaa\x1a\x2c\x5f\x6f\xc7\xc6\x6d\x8b\xfe\x65\x57\x5d\x94\xd0\x1f\xbb\xf2\x55\x11\x07\xdb\xa1\xc5\x68\x9e\x55\x6f\x8c\x14\xec\x93\x28\xef\x24\x7c\xcb\xce\x61\x1f\x63\xda\x9c\x73\x5a\xe9\x91\x1c\xfe\x2d\x66\x89\xed\x61\x29\x0b\xfe\x6f\x42\xde\x3b\xa6\x38\x68\x8c\x85\xc1\x91\x6e\xf3\x6d\xbf\x14\xbf\x02\x00\x00\xff\xff\xe5\x21\x7d\x2b\x8f\x03\x00\x00")

func coreEmailTemplatesNewLeaseNoOwnerTagTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesNewLeaseNoOwnerTagTxt,
		"core/email-templates/new-lease-no-owner-tag.txt",
	)
}

func coreEmailTemplatesNewLeaseNoOwnerTagTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesNewLeaseNoOwnerTagTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/new-lease-no-owner-tag.txt", size: 911, mode: os.FileMode(436), modTime: time.Unix(1487264658, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x52\xc1\x8a\xdb\x30\x10\xbd\xfb\x2b\x5e\x73\x69\x02\x21\x7b\x5f\xb2\x86\xb2\x97\xf6\xd4\x43\x73\x37\x4a\x3c\xb6\x87\x95\xa5\x20\xc9\x35\xc6\xf8\xdf\x3b\xb2\x6c\xb7\x59\x4a\x59\xe8\xc1\x60\x66\xde\x9b\x79\xef\x8d\xc6\x91\x2b\x18\x1b\x70\xf2\x41\xdd\xde\x0a\xa3\x5a\xc2\x34\x65\xc0\x57\x1a\x30\x8e\x27\xdb\x1b\x72\x05\xb5\x8a\xf5\x34\x1d\xe1\x6d\x4b\xd6\x10\x6e\x8e\x54\xa0\x12\x0a\x86\x7a\xb0\x11\xb6\xb9\x91\xd0\xf6\x5c\xe2\x7c\xcd\x85\xb9\x16\x0b\x2e\xa7\xe9\xfc\x74\xcd\x8f\xb0\x15\xc2\x70\xa7\xf7\x80\x58\x5b\x20\x32\xc1\x9a\xa5\xef\xc8\xdb\xce\x49\xdf\x51\xcd\xd6\x24\xc4\xe1\x24\x5d\x97\xc7\x2f\x1b\x47\x32\x32\x3b\xcb\x66\x17\xff\xe5\x60\xe6\x46\xf9\x1f\xdf\x2e\xe8\x1f\x91\x85\xb8\xf1\x79\x61\xfd\xd6\x90\x08\x33\x72\x05\x72\xf9\x08\x5b\x83\x79\x34\x73\x69\x08\xaf\x74\x63\xfd\x3d\x0a\x47\x50\x75\x8c\x24\x34\xec\xb7\x9c\x21\xff\xf1\x68\x1c\xeb\x84\xbe\xe1\x40\x9a\x7d\x88\xee\xd0\x13\x94\xf7\x5c\x1b\x31\xc7\x01\xc1\x62\xb0\x1d\xf6\x11\xa8\xca\x96\xcd\xe1\x94\x65\xf3\xc6\xe4\xe2\x5b\xba\xbf\xba\xdf\x9d\xfd\x49\xe5\x31\x72\x7a\xd6\x1a\x57\x42\x20\x27\x84\x14\x53\x58\xa4\xaf\x35\xc9\xa4\x08\xbc\xfa\xc4\x5e\x5a\x9a\x94\xa7\xa2\xec\xdc\xdc\x9d\x26\xa8\x4a\xd0\x32\xf0\xb3\x4f\x71\x4b\xf5\xdd\xf6\xcb\xba\x01\xdc\xb6\x54\xb2\xfc\xe9\xe1\xf9\x0f\xc4\x59\xa1\x71\x54\xbd\xec\xb6\xf9\x9b\xa8\xa2\x73\x72\xd1\x9d\x24\xe4\x6a\x0a\x2f\xbb\xe2\xaa\x95\x79\xdb\xe5\xaf\x9a\x25\xec\x86\x1c\x45\xf3\xa2\x7a\x63\xcc\x61\x3f\xa9\xfc\x41\xc2\x97\xe4\x1c\xfb\x18\xd3\xe6\x5c\xd2\x9a\x1f\xce\xe1\xdf\x62\x96\xd8\x3e\x2c\x65\xc1\xff\x4d\xc8\xa5\x11\x8a\x47\x65\x1d\x3a\xcf\xa6\x4e\x8f\xe0\x53\xf6\x2b\x00\x00\xff\xff\x97\x28\x70\x5e\xa3\x03\x00\x00")

func coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxt,
		"core/email-templates/new-lease-owner-tag-not-whitelisted.txt",
	)
}

func coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/new-lease-owner-tag-not-whitelisted.txt", size: 931, mode: os.FileMode(436), modTime: time.Unix(1487264662, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x53\xc1\x6e\xdb\x30\x0c\xbd\xfb\x2b\xb8\x5c\x96\x02\x41\x72\x2f\x5a\x03\x45\x2f\xbb\x6d\xc0\x76\x17\x14\x9b\x8e\x89\xc8\x52\x20\x29\x31\x0c\xc3\xff\x3e\x8a\x52\xb2\x35\xe8\xb0\x5e\x7a\x08\x14\x53\xef\x91\xef\x3d\xda\xf3\x4c\x1d\x58\x17\x61\x1b\xa2\x6e\x8e\xca\xea\x01\x61\x59\x2a\x80\x6f\x38\xc1\x3c\x6f\xdd\x68\xd1\x2b\x1c\x34\x99\x65\xd9\xc0\xe4\xce\xb0\x76\x1e\x82\x1b\xd0\x59\x04\x34\x01\xe1\x1c\xc8\x1e\xd2\x95\x87\x57\x6c\xc8\x7c\x4f\x1c\x88\xfa\x00\x3b\x38\xe2\x94\x7a\x3e\x40\xe3\x51\x47\x6c\x41\x83\xc5\x11\xc8\xf2\x3c\xdb\x20\x0f\x02\x58\x53\x0b\x4f\xfb\x9a\xa7\x5d\xcb\x8a\xda\x65\x79\xda\xed\xeb\x0d\xb8\x0e\xe2\x74\xc2\x7b\x40\xaa\x15\x88\xf4\x70\xb6\x20\x3c\x06\x16\xc2\x08\x8f\x07\x72\x36\x63\x1e\xb6\x7c\xeb\xeb\xf4\xab\xe6\x19\x2d\x77\xaf\x2a\xf1\xfe\xc9\xbe\x93\x2e\x6d\x6f\x76\xc5\x4d\x8f\x20\x43\xef\x33\x91\x62\x0e\xe4\x7f\x6e\x2a\x86\x89\x99\x72\xf0\xf9\x33\xb1\x21\xcd\x7c\x2c\xdc\x3f\xce\x32\xad\x10\x32\x90\xda\xb7\xb0\x6b\xe0\x6f\x23\x02\x78\x89\x22\x38\xd2\x20\xe2\x47\x4f\x31\x99\x8e\x3d\x05\x90\x74\x72\x36\xbd\xbe\x60\x8a\xcd\x2a\xd7\x29\xdd\x44\xba\xa0\x32\xa8\x03\x86\x65\x81\xfc\x2c\xce\x72\x6d\xc3\x39\xc2\x88\xec\x9a\xcd\x4b\x7e\xfa\x74\xf2\xee\xa2\x0d\x74\x1c\xb2\x34\xe7\x94\xff\x5a\x1a\x93\x7f\x08\x17\x1a\x43\x2c\x9f\x03\x5a\xbd\x08\x07\x57\x10\x5d\xe1\x63\xa6\x0a\xb0\xbc\x16\x2c\xde\x8f\xc4\x3c\x8a\x30\x92\x31\xb0\x67\x10\xfa\x81\x6c\x4e\x3e\x96\x14\xae\x35\x0e\x59\x25\xb3\x39\x0d\x58\xf3\x95\xb4\x53\xed\xd9\xcb\x6d\xf2\xd3\x31\x9a\x1b\x7e\x0d\x79\x83\x5c\x7d\x77\x25\x45\xe0\xe3\xdd\x0d\x1f\x1a\x7a\x8f\xdd\xf3\xea\xd6\xbd\xe8\x57\x67\xcf\xaf\x1b\x3b\xd2\xfe\x80\xf1\x79\xa5\xf6\x46\xdb\xe3\xaa\x7e\x15\xd3\x6c\x05\x93\x59\x56\x5c\xf0\xb2\xb1\x9d\xae\xdf\x19\xfe\xeb\x6a\x12\x68\x18\xb0\x25\xfe\x67\xa6\x8f\x48\xb9\xa5\xf3\x61\x31\x37\xc6\xbf\xe4\xb0\x9a\x9e\xa9\x41\xb6\x9b\xbf\x1a\xf9\x60\xbe\x54\xbf\x03\x00\x00\xff\xff\xa9\x3d\x8e\x35\x80\x04\x00\x00")

func coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxt,
		"core/email-templates/new-lease-valid-owner-tag-needs-approval.txt",
	)
}

func coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/new-lease-valid-owner-tag-needs-approval.txt", size: 1152, mode: os.FileMode(436), modTime: time.Unix(1487264666, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x52\x4d\x8b\xdb\x30\x10\xbd\xfb\x57\xbc\xcd\xa5\x09\x2c\xc9\x7d\xc9\x1a\xca\xb2\xd0\x4b\xd9\xc3\x2e\x94\x9e\x8c\x1c\x8f\xe3\x21\xb2\x54\x64\x85\x60\x8c\xfe\x7b\x47\xf2\x47\x5a\xd3\x63\xf7\x24\x31\x7a\x33\xef\x43\x33\x0c\x5c\xc3\x58\x8f\x7d\xe7\xd5\xe9\x52\x18\xd5\x12\x42\xc8\x80\x6f\xd4\x63\x18\xf6\xf6\x66\xc8\x15\xd4\x2a\xd6\x21\x3c\xa2\xb7\x57\x6c\xad\x43\x67\x5b\xb2\x86\x40\xba\x23\x5c\x3b\x36\xe7\xf8\xe4\xf0\x42\x27\xd6\x6f\xb1\x07\x5e\x9d\x81\x03\x2e\xd4\xc7\xa1\x3b\x9c\x1c\x29\x4f\x15\x14\x0c\xdd\xc0\x46\x08\xcd\x89\x84\x09\xd8\x72\x85\x63\x99\x0b\xdd\x5c\x2e\xb8\x0a\xe1\x78\x28\xf3\x47\xd8\x1a\xbe\xff\x45\x6b\x40\xac\x4d\x90\x34\xc3\x9a\x09\xe1\xa8\x13\x25\x82\x70\x74\x66\x6b\x46\xcc\x6e\x8f\x8f\x46\xf9\x2f\x1d\xbe\xfe\x78\x7d\x7f\xfb\xfe\xfa\x90\xa5\xae\x63\xe9\xf2\xfb\x65\x18\xc8\x08\x6f\x96\xa5\x58\x3e\x27\x92\x25\x91\xa8\x58\x99\x25\x88\xe4\xb3\x21\x24\xd2\x75\x5a\xa9\x38\x46\xf5\x3f\x7c\xca\xed\x3d\x4e\x44\xd4\xf1\x34\xcd\xbb\xbb\x1d\x47\x4d\x2d\x23\x90\xab\xbf\x61\xf3\xf7\x4c\xa0\x3f\x8f\x25\x44\xe0\x67\xf4\xbf\x18\xbc\xb1\xd6\x28\x09\x9e\x5c\xcb\x66\x74\xe7\xa7\xa9\x73\x4d\x8c\x14\x9e\x67\x09\xd8\xca\x93\x26\xd5\x51\x51\x5d\x5d\x7a\x0d\x01\xaa\x16\x34\x38\x9a\x4c\x29\x49\x75\xb7\xcf\x56\x2a\xe4\xfc\x98\x79\xc0\x6d\x4b\x15\xcb\x4d\xf7\x4f\x2b\x9c\x1c\x0a\x8d\xa3\xfa\x79\xb3\x70\x2d\x02\x8b\xab\x93\x2f\xde\xc8\xbf\xb9\x33\xf9\xe7\x4d\x51\x6a\x65\x2e\x9b\xfc\x45\xb3\x64\xd2\x90\x13\x33\x36\x3a\x58\x3a\x52\x26\x07\x95\xff\x4b\x4e\x23\xbd\x1d\x6a\x59\x95\x71\x3b\xd2\x62\x3c\x64\xbf\x03\x00\x00\xff\xff\x0c\xd2\x3a\x83\x83\x03\x00\x00")

func coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxt,
		"core/email-templates/new-lease-valid-owner-tag-no-approval-needed.txt",
	)
}

func coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/new-lease-valid-owner-tag-no-approval-needed.txt", size: 899, mode: os.FileMode(436), modTime: time.Unix(1487264675, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _coreEmailTemplatesRegionSuccessfullySetupTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x4c\x8e\xd1\x6a\x86\x30\x0c\x85\xef\x7d\x8a\xe3\xbd\xb8\x17\x10\x2f\xc7\x1e\x60\xf7\x92\x6a\xa6\x61\x6d\x3a\x9a\x16\x11\xe9\xbb\xaf\x4c\x06\xff\x4d\x20\xf9\xc8\xf9\xce\x07\x5f\xb8\xef\x91\xb6\x20\xba\x70\x20\xf1\xb5\x0e\xc8\x07\x23\xf1\x2e\x51\x31\xb9\xb9\xf1\x67\x59\x94\x02\xd7\x3a\xbd\xb9\x19\x07\x19\x1c\xb3\xc2\xca\xba\xb2\xd9\x57\xf1\xfe\x82\x71\x2e\x3f\x7d\x37\xb9\x34\x3f\xe3\x3d\xc5\x00\x8d\x27\xa2\x0e\x20\xef\x21\x6a\x99\xb4\x7d\xb4\x4b\xf3\x88\xfd\x8b\x4e\x69\xd4\x31\x42\x54\xc9\x31\xf1\x06\x77\xfd\x15\xd9\x0b\xa5\x4d\x48\xc7\x97\xd8\xcf\x83\xf4\xdb\xfa\xee\x37\x00\x00\xff\xff\xec\x7d\x16\x9b\xbf\x00\x00\x00")

func coreEmailTemplatesRegionSuccessfullySetupTxtBytes() ([]byte, error) {
	return bindataRead(
		_coreEmailTemplatesRegionSuccessfullySetupTxt,
		"core/email-templates/region-successfully-setup.txt",
	)
}

func coreEmailTemplatesRegionSuccessfullySetupTxt() (*asset, error) {
	bytes, err := coreEmailTemplatesRegionSuccessfullySetupTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "core/email-templates/region-successfully-setup.txt", size: 191, mode: os.FileMode(436), modTime: time.Unix(1486826381, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"core/email-templates/account-verification-notification.txt": coreEmailTemplatesAccountVerificationNotificationTxt,
	"core/email-templates/expiring-lease.txt": coreEmailTemplatesExpiringLeaseTxt,
	"core/email-templates/lease-approved.txt": coreEmailTemplatesLeaseApprovedTxt,
	"core/email-templates/lease-extended.txt": coreEmailTemplatesLeaseExtendedTxt,
	"core/email-templates/lease-resource-terminated.txt": coreEmailTemplatesLeaseResourceTerminatedTxt,
	"core/email-templates/misconfiguration-notice.txt": coreEmailTemplatesMisconfigurationNoticeTxt,
	"core/email-templates/new-lease-no-owner-tag.txt": coreEmailTemplatesNewLeaseNoOwnerTagTxt,
	"core/email-templates/new-lease-owner-tag-not-whitelisted.txt": coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxt,
	"core/email-templates/new-lease-valid-owner-tag-needs-approval.txt": coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxt,
	"core/email-templates/new-lease-valid-owner-tag-no-approval-needed.txt": coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxt,
	"core/email-templates/region-successfully-setup.txt": coreEmailTemplatesRegionSuccessfullySetupTxt,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"core": &bintree{nil, map[string]*bintree{
		"email-templates": &bintree{nil, map[string]*bintree{
			"account-verification-notification.txt": &bintree{coreEmailTemplatesAccountVerificationNotificationTxt, map[string]*bintree{}},
			"expiring-lease.txt": &bintree{coreEmailTemplatesExpiringLeaseTxt, map[string]*bintree{}},
			"lease-approved.txt": &bintree{coreEmailTemplatesLeaseApprovedTxt, map[string]*bintree{}},
			"lease-extended.txt": &bintree{coreEmailTemplatesLeaseExtendedTxt, map[string]*bintree{}},
			"lease-resource-terminated.txt": &bintree{coreEmailTemplatesLeaseResourceTerminatedTxt, map[string]*bintree{}},
			"misconfiguration-notice.txt": &bintree{coreEmailTemplatesMisconfigurationNoticeTxt, map[string]*bintree{}},
			"new-lease-no-owner-tag.txt": &bintree{coreEmailTemplatesNewLeaseNoOwnerTagTxt, map[string]*bintree{}},
			"new-lease-owner-tag-not-whitelisted.txt": &bintree{coreEmailTemplatesNewLeaseOwnerTagNotWhitelistedTxt, map[string]*bintree{}},
			"new-lease-valid-owner-tag-needs-approval.txt": &bintree{coreEmailTemplatesNewLeaseValidOwnerTagNeedsApprovalTxt, map[string]*bintree{}},
			"new-lease-valid-owner-tag-no-approval-needed.txt": &bintree{coreEmailTemplatesNewLeaseValidOwnerTagNoApprovalNeededTxt, map[string]*bintree{}},
			"region-successfully-setup.txt": &bintree{coreEmailTemplatesRegionSuccessfullySetupTxt, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

