package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"

	"github.com/sgaunet/ssftp/pathh"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func initTrace(debugLevel string) {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{
	// 	DisableColors: true,
	// 	FullTimestamp: true,
	// })

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	switch debugLevel {
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}

var version string = "development"

func printVersion() {
	fmt.Println(version)
}

func usage() {
	fmt.Println("ssftp [-i sshkey] [-d debug] src dest")
	fmt.Println("    -i : ssh key")
	fmt.Println("    -d : debug mode")
	fmt.Println("    -v : print version and exit")
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
	// Parameters treatment (except src + dest)
	flag.StringVar(&sshkeyFile, "i", "", "SSH key File")
	flag.StringVar(&debugLevel, "d", "info", "Debug level (info,warn,debug)")
	flag.StringVar(&port, "p", "22", "Port number")
	flag.BoolVar(&vOption, "v", false, "Get version")
	flag.Parse()

	if vOption {
		printVersion()
		os.Exit(0)
	}

	if debugLevel != "info" && debugLevel != "warn" && debugLevel != "debug" {
		fmt.Fprintf(os.Stderr, "debuglevel should be info or warn or debug\n")
		usage()
		os.Exit(1)
	}
	initTrace(debugLevel)

	// src + dest are mandatory parameters
	if len(flag.Args()) != 2 {
		usage()
		os.Exit(0)
	}

	// Parameters treatment : src + dest
	args := flag.Args()
	src := pathh.New(args[0])
	dest := pathh.New(args[1])

	if src.IsRemote() && dest.IsRemote() {
		fmt.Fprintf(os.Stderr, "Cannot transfer from one server to the other\n")
		os.Exit(1)
	}
	if src.IsLocal() && dest.IsLocal() {
		fmt.Fprintf(os.Stderr, "Use cp instead\n")
		os.Exit(1)
	}

	// ssh key is mandatory for the first version
	if len(sshkeyFile) == 0 {
		fmt.Println("No SSH file")
		os.Exit(1)
	}

	// else {
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
		client, err = SftpConnect(src, port, sshkeyFile)
	}
	if dest.IsRemote() {
		client, err = SftpConnect(dest, port, sshkeyFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %s\n", err.Error())
		os.Exit(1)
	}

	// Close connection
	defer client.Close()
	// cwd, err := client.Getwd()
	// println("Current working directory:", cwd)

	if src.IsRemote() {
		log.Debugln("src is remote")
		is, err := IsRemoteFileADir(client, src.GetFilePath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed : %s\n", err.Error())
			os.Exit(1)
		}

		if is {
			err = recursiveDownload(client, src.GetFilePath(), dest.GetFilePath())
		} else {
			err = downloadFile(client, src.GetFilePath(), dest.GetFilePath())
		}
	}

	if dest.IsRemote() {
		log.Debugln("dest is remote")

		// Need info on localpath
		infoLocalFile, err := os.Stat(src.GetFilePath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed : %s\n", err.Error())
			os.Exit(1)
		}

		// If it's a directory, upload every files and keep the same tree
		if infoLocalFile.IsDir() {
			// walk throught the tree
			err := filepath.Walk(src.GetFilePath(),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return errors.New("problem with local file " + path + " : " + err.Error())
					}
					if !info.IsDir() {
						baseDirSrc := filepath.Base(src.GetFilePath()) // dirname of source
						completeRemotePath := filepath.Clean(dest.GetFilePath() + "/" + baseDirSrc + "/" + path[len(src.GetFilePath()):])
						completeRemotePath = strings.ReplaceAll(completeRemotePath, string(os.PathSeparator), "/")
						log.Infof("Upload to : %s (size %v)\n", completeRemotePath, info.Size())
						return uploadFile(client, path, completeRemotePath)
					}
					return nil
				})
			if err != nil {
				log.Errorln(err)
			}
		} else {
			// localpath is a simple file, upload it
			err = uploadFile(client, src.GetFilePath(), dest.GetFilePath())
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed : %s\n", err.Error())
		os.Exit(1)
	}
}
