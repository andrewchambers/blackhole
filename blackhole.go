package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting the current user: %s\n", err)
		os.Exit(1)
	}

	hookDir := filepath.Join(u.HomeDir, ".blackhole_hooks")
	if len(os.Args) == 2 {
		hookDir = os.Args[1]
	}

	temp, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating temporary file: %s\n", err)
		os.Exit(1)
	}
	defer temp.Close()
	defer os.Remove(temp.Name())

	_, err = io.Copy(temp, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during upload: %s\n", err)
		os.Exit(1)
	}

	err = temp.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finishing upload: %s\n", err)
		os.Exit(1)
	}

	hooks, err := ioutil.ReadDir(hookDir)
	for _, hook := range hooks {
		hookPath := filepath.Join(hookDir, hook.Name())
		cmd := exec.Command(hookPath, temp.Name())
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			// The hook can provide it's own error message if it wants.
			os.Exit(1)
		}
	}
}
