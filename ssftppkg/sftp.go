package ssftppkg

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/sgaunet/ssftp/pathh"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

const StrictHostKeyCheckingNameParameter = "StrictHostKeyChecking"

type SsftpClient struct {
	log             *logrus.Logger
	hostKeyChecking bool
}

func NewSsftpClient(log *logrus.Logger, sshOpts SshOptions) (*SsftpClient, error) {
	s := SsftpClient{
		log: log,
	}
	s.hostKeyChecking = true // Default behaviour

	for _, opt := range sshOpts {
		err := s.SetOption(opt)
		if err != nil {
			return &s, err
		}
	}
	return &s, nil
}

// SetOption is waiting for parameter with this format option=value
func (s *SsftpClient) SetOption(param string) error {
	splitted := strings.Split(param, "=")
	if len(splitted) != 2 {
		return errors.New("sshOption should have the format : param=value")
	}
	parameter := splitted[0]
	value := splitted[1]
	switch parameter {
	case StrictHostKeyCheckingNameParameter:
		if value == "yes" {
			s.log.Debugln("Option hostKeyChecking=true")
			s.hostKeyChecking = true
		} else {
			s.hostKeyChecking = false
		}
	default:
		return errors.New("Unknown option" + param)
	}
	return nil
}

func (s *SsftpClient) PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func (s *SsftpClient) ListFiles(client *sftp.Client, remoteDir string) (err error) {
	files, err := client.ReadDir(remoteDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		return
	}

	for _, f := range files {
		var name, modTime, size string

		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")
		size = fmt.Sprintf("%12d", f.Size())

		if f.IsDir() {
			name = name + "/"
			modTime = ""
			size = "N/A"
		}
		fmt.Fprintf(os.Stdout, "%19s %12s %s\n", modTime, size, name)
	}

	return
}

// Upload file to sftp server
func (s *SsftpClient) UploadFile(client *sftp.Client, localFile, remoteFile string) (err error) {
	s.log.Infof("Uploading %s to %s ...", localFile, remoteFile)
	s.log.Debugln("remoteFile=", remoteFile)
	srcFile, err := os.Open(localFile)
	if err != nil {
		return errors.New("Unable to open local file" + localFile + " : " + err.Error())
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := strings.ReplaceAll(filepath.Dir(remoteFile), "\\", "/")
	s.log.Debugln("parent=", parent)
	pathSeparator := string("/")
	dirs := strings.Split(parent, string(pathSeparator))
	remotepath := ""
	if len(dirs) == 1 {
		remotepath = parent
		s.log.Debugln("remotepath=", remotepath)
		client.Mkdir(remotepath) // should handle the error
	} else {
		for _, dir := range dirs {
			s.log.Debugln("dir=", dir)
			remotepath = strings.ReplaceAll(remotepath+pathSeparator+dir, "\\", "/")
			remotepath = strings.ReplaceAll(remotepath, "//", "/")
			s.log.Debugln("remotepath=", remotepath)
			client.Mkdir(remotepath) // should handle the error

			// log.Infoln("Create remote dir :", remotepath)
		}
	}

	// If remoteFile is a dir ...
	infoRemote, err := client.Stat(remoteFile)
	if err == nil {
		if infoRemote.IsDir() {
			remoteFile = remoteFile + "/" + filepath.Base(localFile)
		}
	}

	dstFile, err := client.OpenFile(remoteFile, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to open remote file: %v\n", err)
		return errors.New("unable to open remote file " + remoteFile + ":" + err.Error())
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to upload local file: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "(%d bytes copied)\n", bytes)

	return
}

// Download file from sftp server
func (s *SsftpClient) DownloadFile(client *sftp.Client, remoteFile, localFile string) (err error) {

	fmt.Fprintf(os.Stdout, "Downloading [%s] to [%s] ... ", remoteFile, localFile)

	srcFile, err := client.OpenFile(remoteFile, (os.O_RDONLY))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to open remote file: %v\n", err)
		return
	}
	defer srcFile.Close()

	// Check if dir exists
	dir := filepath.Dir(localFile)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			//fmt.Println("CREATE :", dir)
			os.Mkdir(dir, 0755)
		} else {
			// other error
			fmt.Fprint(os.Stderr, err.Error())
		}
	}

	dstFile, err := os.Create(localFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to open local file: %v\n", err)
		return
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to download remote file: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "(%d bytes copied)w\n", bytes)

	return
}

func (s *SsftpClient) SftpConnect(remote pathh.Path, port string, sshkeyFile string) (*sftp.Client, error) {
	privateKey, _ := os.ReadFile(sshkeyFile)
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("error setting up SSH config: %s", err)
	}

	authMethods := []ssh.AuthMethod{ssh.PublicKeys(signer)}
	// auth := s.PublicKeyFile(sshkeyFile)
	// if auth == nil {
	// 	panic("Key not found")
	// }
	sshConfig := ssh.ClientConfig{}
	sshConfig.SetDefaults()
	sshConfig.User = remote.GetUser()
	// sshConfig.Auth = []ssh.AuthMethod{
	// 	auth,
	// }
	sshConfig.Auth = authMethods
	if s.hostKeyChecking {
		s.log.Debugln("hostKeyChecking=true")
		hostKeyCallback, err := kh.New(fmt.Sprintf("%s/.ssh/known_hosts", os.Getenv("HOME")))
		if err != nil {
			log.Fatal("could not create hostkeycallback function: ", err)
		}
		sshConfig.HostKeyCallback = hostKeyCallback
	} else {
		s.log.Debugln("hostKeyChecking=false")
		sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	// sshConfig.HostKeyCallback = KeyPrint

	sshConfig.Timeout = 2 * time.Second
	conn, err := ssh.Dial("tcp", remote.GetServer()+":"+port, &sshConfig)
	if err != nil {
		return nil, errors.New("Failed to dial: " + err.Error())
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return client, err
}

func (s *SsftpClient) IsRemoteFileADir(client *sftp.Client, remoteFile string) (bool, error) {

	info, err := client.Stat(remoteFile)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return true, err
	}
	return false, err
}

func (s *SsftpClient) RecursiveDownload(client *sftp.Client, remoteFile string, localFile string) (err error) {
	files, err := client.ReadDir(remoteFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		err = errors.New("error during the recursive download")
		return
	}

	for _, f := range files {
		var name string
		name = f.Name()

		if f.IsDir() {
			name = name + "/"
			err2 := s.RecursiveDownload(client, remoteFile+"/"+name, localFile+string(os.PathSeparator)+name)
			if err2 != nil {
				err = errors.New("error during the recursive download")
			}
		} else {
			err = s.DownloadFile(client, remoteFile+"/"+name, localFile+string(os.PathSeparator)+name)
			if err != nil {
				err = errors.New("error during the recursive download")
			}
		}
	}

	return
}

// func fingerprintCallback(opts *ssh.Config, expectedFingerprint string) ssh.HostKeyCallback {
// 	if opts.SkipHostValidation {
// 		return nil
// 	}

// 	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
// 		switch len(expectedFingerprint) {
// 		case helpers.SHA1_FINGERPRINT_LENGTH:
// 			fingerprint := helpers.SHA1Fingerprint(key)
// 			if fingerprint != expectedFingerprint {
// 				return fmt.Errorf("Host key verification failed.\n\nThe fingerprint of the received key was %q.", fingerprint)
// 			}
// 		case helpers.MD5_FINGERPRINT_LENGTH:
// 			fingerprint := helpers.MD5Fingerprint(key)
// 			if fingerprint != expectedFingerprint {
// 				return fmt.Errorf("Host key verification failed.\n\nThe fingerprint of the received key was %q.", fingerprint)
// 			}
// 		case 0:
// 			fingerprint := helpers.MD5Fingerprint(key)
// 			return fmt.Errorf("Unable to verify identity of host.\n\nThe fingerprint of the received key was %q.", fingerprint)
// 		default:
// 			return errors.New("Unsupported host key fingerprint format")
// 		}
// 		return nil
// 	}
// }

func KeyPrint(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	fmt.Printf("%s %s %s\n", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	return nil
}
