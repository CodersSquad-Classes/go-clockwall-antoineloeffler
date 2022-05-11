package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type clock struct {
	name, host string
}

func (c *clock) watch(w io.Writer, r io.Reader) {
	a := bufio.NewScanner(r)
	for a.Scan() {
		fmt.Fprintf(w, "%s: %s\n", c.name, a.Text())
	}
	fmt.Println(c.name, "done")
	if a.Err() != nil {
		log.Printf("can't read from %s: %s", c.name, a.Err())
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: clockwall NAME=HOST ...")
		os.Exit(1)
	}
	clocks := make([]*clock, 0)
	for _, a := range os.Args[1:] {
		fields := strings.Split(a, "=")
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "bad arg: %s\n", a)
			os.Exit(1)
		}
		clocks = append(clocks, &clock{fields[0], fields[1]})
	}
	for _, c := range clocks {
		conn, err := net.Dial("tcp", c.host)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		go c.watch(os.Stdout, conn)
	}

	for {
		time.Sleep(time.Minute)
	}
}
