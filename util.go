package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ErrCardFull = errors.New("card full")

// Run a command
func runForeground(cmd string, args []string) error {
	fmt.Println(cmd, strings.Join(args, " "))
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	return c.Run()
}
