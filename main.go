package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

func main() {
	//do nothing
}

func runCommandPID(closed <-chan struct{}, wg *sync.WaitGroup, command string) {
	defer wg.Done()

	tokens := strings.Split(command, " ")
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return
	}

	finished := make(chan struct{})

	go func() {
		_ = cmd.Wait()
		close(finished)
	}()

	for {
		select {

		case <-closed:

			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				log.Fatalf("failed to kill process: %v", err)
			}
			return

		case <-finished:
			return
		}

	}

}

func runCommandGID(closed <-chan struct{}, wg *sync.WaitGroup, command string) {
	defer wg.Done()

	tokens := strings.Split(command, " ")
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Stdout = os.Stdout
	// https://varunksaini.com/posts/kiling-processes-in-go/
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return
	}

	finished := make(chan struct{})

	go func() {
		_ = cmd.Wait()
		close(finished)
	}()

	for {
		select {

		case <-closed:

			pgid, err := syscall.Getpgid(cmd.Process.Pid)

			if err == nil {
				if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
					log.Fatalf("failed to kill process: %v ", err)
				}
			} else {
				log.Fatalf("failed to get pgid because %v", err)
			}

		case <-finished:
			return
		}

	}

}
