package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"

	"github.com/sgaunet/ssftp/pkg/sftpclient"
	"github.com/sgaunet/ssftp/pkg/sftppath"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var version string = "development"

func printVersion() {
	fmt.Println(version)
}

func usage() {
	fmt.Println("ssftp [-i sshkey] [-d debug] src dest")
	fmt.Println("    -i : ssh key")
	fmt.Println("    -d : debug mode")
	fmt.Println("    -v : print version and exit")
	fmt.Println("    -p : port")
	fmt.Println("    src: local file/folder or distant sftp file/dir")
	fmt.Println("    dest: same")
	fmt.Println("\nOrder of paramters matters.")
}

func main() {
	var err error
	var client *sftp.Client
	// Arguments
	var sshkeyFile string
	var debugLevel string
	var port string
	var vOption bool
	var sshOpts sftpclient.SshOptions
	// Parameters treatment (except src + dest)
	flag.StringVar(&sshkeyFile, "i", "", "SSH key File")
	flag.StringVar(&debugLevel, "d", "info", "Debug level (info,warn,debug)")
	flag.Var(&sshOpts, "o", "Options (Ex: StrictHostKeyChecking=no) ")
	flag.StringVar(&port, "p", "22", "Port number")
	flag.BoolVar(&vOption, "v", false, "Get version")
	flag.Parse()

	if vOption {
		printVersion()
		os.Exit(0)
	}

	if debugLevel != "info" && debugLevel != "warn" && debugLevel != "debug" {
		log.Errorf("debuglevel should be info or warn or debug\n")
		usage()
		os.Exit(1)
	}
	initTrace(debugLevel)
	s, err := sftpclient.NewSsftpClient(log, sshOpts)
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}

	// src + dest are mandatory parameters
	if len(flag.Args()) != 2 {
		usage()
		os.Exit(0)
	}

	// Parameters treatment : src + dest
	args := flag.Args()
	src := sftppath.New(args[0])
	dest := sftppath.New(args[1])

	if src.IsRemote() && dest.IsRemote() {
		log.Errorf("Cannot transfer from one server to the other\n")
		os.Exit(1)
	}
	if src.IsLocal() && dest.IsLocal() {
		log.Errorf("Use cp instead\n")
		os.Exit(1)
	}

	// ssh key is mandatory for the first version
	if len(sshkeyFile) == 0 {
		log.Errorln("No SSH file")
		os.Exit(1)
	}

	// 	sshConfig = ssh.ClientConfig{
	// 		User: "vagrant",
	// 		Auth: []ssh.AuthMethod{
	// 			ssh.Password("vagrant"),
	// 		},
	// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// 		Timeout:         2 * time.Second,
	// 	}
	// }

	if src.IsRemote() {
		client, err = s.SftpConnect(src, port, sshkeyFile)
	}
	if dest.IsRemote() {
		client, err = s.SftpConnect(dest, port, sshkeyFile)
	}

	if err != nil {
		log.Errorf("Failed to connect: %s\n", err.Error())
		os.Exit(1)
	}

	// Close connection
	defer client.Close()
	// cwd, err := client.Getwd()
	// println("Current working directory:", cwd)

	if src.IsRemote() {
		var isRemote bool
		log.Debugln("src is remote")
		isRemote, err = s.IsRemoteFileADir(client, src.GetFilePath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed : %s\n", err.Error())
			os.Exit(1)
		}

		if isRemote {
			err = s.RecursiveDownload(client, src.GetFilePath(), dest.GetFilePath())
			if err != nil {
				log.Errorf("Failed to recursive download : %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			err = s.DownloadFile(client, src.GetFilePath(), dest.GetFilePath())
			if err != nil {
				log.Errorf("Failed to download file %s : %s\n", src.GetFilePath(), err.Error())
				os.Exit(1)
			}
		}
	}

	if dest.IsRemote() {
		log.Debugln("dest is remote")
		// Need info on localpath
		log.Debugln("src.GetFilePath()=", src.GetFilePath())

		// If it's a directory, upload every files and keep the same tree
		if isDirExists(src.GetFilePath()) {
			// walk throught the tree
			err = filepath.Walk(src.GetFilePath(),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return errors.New("problem with local file " + path + " : " + err.Error())
					}
					if !isDirExists(path) {
						// ssftp ...  localdir  user@server:/abso
						// ssftp ...  ./localdir  user@server:/abso
						// baseDirSrc := filepath.Base(src.GetFilePath()) // dirname of source
						completeRemotePath := filepath.Clean(dest.GetFilePath() + "/" + filepath.Base(src.GetFilePath()) + "/" + path[len(src.GetFilePath()):])
						completeRemotePath = filepath.ToSlash(completeRemotePath)
						// log.Infof("Upload to : %s (size %v)\n", completeRemotePath, info.Size())
						return s.UploadFile(client, path, completeRemotePath)
					}
					return nil
				})
			if err != nil {
				log.Errorln(err)
			}
		} else {
			// localpath is a simple file, upload it
			err = s.UploadFile(client, src.GetFilePath(), dest.GetFilePath())
			if err != nil {
				log.Errorf("Failed to upload file %s : %s\n", src.GetFilePath(), err.Error())
				os.Exit(1)
			}
		}
	}
}

func isDirExists(dir string) bool {
	log.Debugf("isDirExists(%v)", dir)
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	i, _ := os.Stat(dir)
	return i.IsDir()
}

func isFileExists(file string) bool {
	log.Debugf("isFileExists(%v)", file)
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	i, _ := os.Stat(file)
	return !i.IsDir()
}

func isThereAFileOrDir(file string) bool {
	log.Debugf("isThereAFileOrDir(%v)", file)
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	return true
}
