package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestKillSleep(t *testing.T) {

	command := "sleep 10"

	closed := make(chan struct{})
	finished := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		runCommandPID(closed, &wg, command)
		close(finished)
	}()

	<-time.After(100 * time.Millisecond)

	close(closed)

	select {
	case <-time.After(time.Second):
		t.Error("time out waiting for kill")
	case <-finished:
	}
}

func TestPIDAccepts2(t *testing.T) {

	command := "./sleep.sh"

	closed := make(chan struct{})
	finished := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		runCommandPID(closed, &wg, command)
		close(finished)
	}()

	<-time.After(100 * time.Millisecond)

	close(closed)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("time out waiting for kill")
	case <-finished:
	}

}

func TestPIDIgnores2(t *testing.T) {

	command := "./sleepsoundly.sh"

	//"ffmpeg -f v4l2 -framerate 25 -video_size 640x480 -i /dev/video1 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 1000k -bf 0 - &"

	closed := make(chan struct{})
	finished := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		runCommandPID(closed, &wg, command)
		close(finished)
	}()

	<-time.After(100 * time.Millisecond)

	close(closed)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("time out waiting for kill")
	case <-finished:
	}

}

func TestPIDNoChildren(t *testing.T) {

	command := "./sleeploudly.sh"

	if _, err := os.Stat("./sleep.log"); err == nil {
		if err := os.Remove("./sleep.log"); err != nil {
			log.Fatalf("Error removing sleep.log: %v", err)
		}
	}

	closed := make(chan struct{})
	finished := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		runCommandPID(closed, &wg, command)
		close(finished)
	}()

	<-time.After(100 * time.Millisecond)

	close(closed)

	select {
	case <-time.After(time.Second):
		t.Error("time out waiting for kill")
	case <-finished:
	}

	//  wait for child proceses to finish, if still running
	<-time.After(3 * time.Second)
	n, err := countlines("./sleep.log")
	if err != nil {
		t.Fatalf("Error reading file was %v", err)
	}

	if n > 1 {
		t.Fatalf("Process was not killed, %d lines in file, wanted 1\n", n)
	}

}

func TestPIDWithChildren(t *testing.T) {

	command := "./sleeplots.sh"
	if _, err := os.Stat("./sleep.log"); err == nil {
		if err := os.Remove("./sleep.log"); err != nil {
			log.Fatalf("Error removing sleep.log: %v", err)
		}
	}

	closed := make(chan struct{})
	finished := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		runCommandPID(closed, &wg, command)
		close(finished)
	}()

	<-time.After(100 * time.Millisecond)

	close(closed)

	select {
	case <-time.After(time.Second):
		t.Error("time out waiting for kill")
	case <-finished:
	}

	//  wait for child proceses to finish, if still running
	<-time.After(3 * time.Second)
	n, err := countlines("./sleep.log")
	if err != nil {
		t.Fatalf("Error reading file was %v", err)
	}

	if n > 3 {
		t.Fatalf("Child processes were not killed: %d lines in file, wanted 3\n", n)
	}

}

func TestGIDWithChildren(t *testing.T) {

	command := "./sleeplots.sh"
	if _, err := os.Stat("./sleep.log"); err == nil {
		if err := os.Remove("./sleep.log"); err != nil {
			log.Fatalf("Error removing sleep.log: %v", err)
		}
	}

	closed := make(chan struct{})
	finished := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		runCommandGID(closed, &wg, command)
		close(finished)
	}()

	<-time.After(100 * time.Millisecond)

	close(closed)

	select {
	case <-time.After(time.Second):
		t.Error("time out waiting for kill")
	case <-finished:
	}

	//  wait for child proceses to finish, if still running
	<-time.After(3 * time.Second)
	n, err := countlines("./sleep.log")
	if err != nil {
		t.Fatalf("Error reading file was %v", err)
	}

	if n > 3 {
		t.Fatalf("Child processes were not killed: %d lines in file, wanted 3\n", n)
	}

}

// https://intelligentbee.com/2017/05/08/counting-lines-words-using-go/

func countlines(name string) (int, error) {

	f, err := os.Open(name)

	if err != nil {
		return -1, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanLines)

	// Count the lines.
	count := 0
	for scanner.Scan() {
		count++
	}

	return count, err

}
