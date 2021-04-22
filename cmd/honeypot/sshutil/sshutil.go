package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/util"
	gossh "golang.org/x/crypto/ssh"
	"io/fs"
	"io/ioutil"
	"os"
)

const (
	HostKeyCannotRead        = "unable to read HostKey from file"
	HostKeyNotParsable       = "HostKey cannot be parsed"
	HostKeyCannotBeGenerated = "HostKey cannot be generated"
	NoHostKeyMarker          = "HostKey is not provided and generation is skipped"
)

//HostKey generates or reads host key file, used to identify server
func HostKey() (gossh.Signer, error) {
	var hostKeyFile string
	if config.GetAppConfig().HostKey != "" {
		hostKeyFile = config.GetAppConfig().HostKey
	} else if config.GetAppConfig().GenerateHostKey {
		hostKeyFile = os.TempDir() + string(os.PathSeparator) + "honeypot.id_rsa"
		if !util.IsFileExists(hostKeyFile) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return nil, errors.New(HostKeyCannotBeGenerated + " Failed to generate key")
			}
			validateErr := privateKey.Validate()
			if validateErr != nil {
				return nil, errors.New(HostKeyCannotBeGenerated + " Failed to validate generated key")
			}
			privateKeyPem := pem.EncodeToMemory(
				&pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
				},
			)
			writeErr := ioutil.WriteFile(hostKeyFile, privateKeyPem, fs.FileMode(0700))
			if writeErr != nil {
				return nil, errors.New(HostKeyCannotBeGenerated + " Failed to write generated key to file")
			}
		}
	} else {
		hostKeyFile = ""
	}

	if hostKeyFile == "" {
		return nil, errors.New(NoHostKeyMarker)
	}
	hostKeyBytes, readError := ioutil.ReadFile(hostKeyFile)
	if readError != nil {
		return nil, errors.New(HostKeyCannotRead)
	}
	hostKey, parseError := gossh.ParsePrivateKey(hostKeyBytes)
	if parseError != nil {
		return nil, errors.New(HostKeyNotParsable)
	}
	return hostKey, nil
}
