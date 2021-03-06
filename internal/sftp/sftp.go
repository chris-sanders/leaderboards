package sftp

import (
	"fmt"
	"github.com/chris-sanders/leaderboards/internal/cfg"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
)

var client *sftp.Client
var config *cfg.Config

func InitClient(iconfig *cfg.Config) {
	config = iconfig
	sshConfig := &ssh.ClientConfig{
		User:            config.Sync.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Sync.Password),
		},
	}

	host := fmt.Sprintf("%v:%v", config.Sync.Server, config.Sync.Port)
	ssh_client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		log.Panic("Failed to dial: " + err.Error())
	}
	log.Info("Successfully connected to ssh server.")

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(ssh_client)
	if err != nil {
		log.Fatal(err)
	}
	client = sftp
}

func CloseClient() {
	client.Close()
}

func GetRemoteFiles() {
	db_path := fmt.Sprintf("%v/*-db.dat", config.Sync.Folder)
	db_files, _ := client.Glob(db_path)
	for _, file := range db_files {
		base := filepath.Base(file)
		if strings.HasPrefix(base, config.Global.Account) {
			log.Tracef("Not downloading remote copy of my own database: %v\n", base)
			continue // do not download our own db file
		}
		log.Info("Downloading database: %v \n", base)
		remote_file, _ := client.Open(file)
		defer remote_file.Close()
		local_file, _ := os.Create(filepath.Base(file))
		defer local_file.Close()
		remote_file.WriteTo(local_file)
	}
}

func UploadDb() {
	local_path := fmt.Sprintf("%v-db.dat", config.Global.Account)
	local_file, _ := os.Open(local_path)
	defer local_file.Close()
	db_path := filepath.Join(config.Sync.Folder, local_path)
	db_file, err := client.Create(filepath.ToSlash(db_path))
	if err != nil {
		log.Error("Failed to create file: ", err.Error())
	}
	defer db_file.Close()
	log.Infof("Uploading file %v", db_path)
	db_file.ReadFrom(local_file)
}
