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
	rc := 0
	func() {
		u, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting the current user: %s\n", err)
			rc = 1
			return
		}

		hookDir := filepath.Join(u.HomeDir, ".blackhole_hooks")
		if len(os.Args) == 2 {
			hookDir = os.Args[1]
		}

		temp, err := ioutil.TempFile("", "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating temporary file: %s\n", err)
			rc = 1
			return
		}
		defer temp.Close()
		defer os.Remove(temp.Name())

		_, err = io.Copy(temp, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error during upload: %s\n", err)
			rc = 1
			return
		}

		err = temp.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error finishing upload: %s\n", err)
			rc = 1
			return
		}

		hooks, err := ioutil.ReadDir(hookDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading hook dir: %s\n", err)
			rc = 1
			return
		}

		if len(hooks) == 0 {
			fmt.Fprintf(os.Stderr, "at least one hook is required.\n")
			rc = 1
			return
		}

		for _, hook := range hooks {
			hookPath := filepath.Join(hookDir, hook.Name())
			cmd := exec.Command(hookPath, temp.Name())
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			if err != nil {
				// The hook can provide it's own error message if it wants.
				rc = 1
				return
			}
		}
	}()

	os.Exit(rc)
}
