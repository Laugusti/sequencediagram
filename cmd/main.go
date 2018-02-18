package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Laugusti/sequencediagram"
	"github.com/Laugusti/sequencediagram/textdiagram"
)

type result struct {
	Diagram string
	Error   string
}

var mode = flag.String("mode", "web", "valid modes are cmd or web")

func main() {
	flag.Parse()

	if *mode != "web" && *mode != "cmd" {
		flag.Usage()
		os.Exit(1)
	}
	if *mode == "web" {
		webServer()
	} else {
		commandLine()
	}
}

func webServer() {
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "textSD.html")
	})
	http.HandleFunc("/creatediagram", func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}
		io.Copy(buf, r.Body) // NOTE: ignoring error
		sd, err := sequencediagram.ParseFromText(buf.String())
		if err != nil {
			log.Printf("failed to parse sequence diagram: %v", err)
			b, err := json.Marshal(result{"", err.Error()})
			if err != nil {
				log.Printf("failed to marshal result: %v", err)
				fmt.Fprintf(w, "failed to marshal result: %v", err)
				return
			}
			w.Write(b)
			return
		}
		buf.Reset()
		io.Copy(buf, textdiagram.Decode(sd))
		b, err := json.Marshal(result{buf.String(), ""})
		if err != nil {
			log.Printf("failed to marshal result: %v", err)
			fmt.Fprintf(w, "failed to marshal result: %v", err)
			return
		}
		log.Println("successfully parsed sequence diagram")
		w.Write(b)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func commandLine() {
	input := bufio.NewScanner(os.Stdin)
	var validLines string
	for input.Scan() {
		line := input.Text()
		sd, err := sequencediagram.ParseFromText(line)

		if err != nil {
			log.Println(err)
			continue
		}
		if len(validLines) > 0 {
			validLines += "\n"
		}
		validLines += line
		fmt.Println("\n")
		// NOTE: ignoring errors
		sd, _ = sequencediagram.ParseFromText(validLines)
		io.Copy(os.Stdout, textdiagram.Decode(sd))
	}
}
