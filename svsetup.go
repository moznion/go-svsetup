package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
)

import "github.com/jessevdk/go-flags"

type cmdOptions struct {
	OptHelp bool `short:"h" long:"help" description:"show this help message and exit"`
}

func main() {
	var err error
	var st int

	defer func() { os.Exit(st) }()

	opts := &cmdOptions{}

	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.Parse()
	if err != nil || len(args) != 1 {
		p.WriteHelp(os.Stdout)
		st = 1
		return
	}

	if opts.OptHelp {
		p.WriteHelp(os.Stdout)
		return
	}

	appName := args[0]
	fmt.Println(appName)

	oldMask := syscall.Umask(0)
	err = os.Mkdir(appName, os.ModeDir)
	if err != nil {
		st = 1
		return
	}
	err = os.Chmod(appName, 0777|os.ModeSticky)
	if err != nil {
		st = 1
		return
	}

	logDir := filepath.Join(appName, "log")
	logMainDir := filepath.Join(logDir, "main")
	err = os.MkdirAll(logMainDir, 0777)
	if err != nil {
		st = 1
		return
	}
	syscall.Umask(oldMask)

	var out *os.File
	runFile := filepath.Join(appName, "run")
	out, _ = os.Create(runFile)

	user, err := user.Current()
	if err != nil {
		st = 1
		return
	}
	out.WriteString(`#!/bin/sh

# HOME="<your home>"

exec 2>&1
# TODO write your own processings
# exec setuidgid ` + user.Username + " <some commands>",
	)

	out.Close()
	os.Chmod(runFile, 0777)

	multilogPath, _ := exec.LookPath("multilog")
	if multilogPath == "" {
		multilogPath = "</path/to/multilog>"
	}
	logSize := "1024"
	rotate := "5"

	logRunFile := filepath.Join(logDir, "run")
	out, _ = os.Create(logRunFile)
	out.WriteString("#!/bin/sh\n" +
		"\n" +
		"USER=" + user.Username + "\n" +
		"MULTILOG=" + multilogPath + "\n" +
		"LOG_SIZE=" + logSize + "\n" +
		"ROTATE=" + rotate + "\n" +
		"\n" +
		"exec 2>&1\n" +
		"exec setuidgid $USER \\\n" +
		"    $MULTILOG t s$LOG_SIZE n$ROTATE ./main",
	)
	os.Chmod(logRunFile, 0777)
	out.Close()
}
