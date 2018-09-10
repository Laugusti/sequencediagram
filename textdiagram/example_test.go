package textdiagram_test

import (
	"io"
	"log"
	"os"

	"github.com/Laugusti/sequencediagram"
	"github.com/Laugusti/sequencediagram/textdiagram"
)

func ExampleEncode() {
	text := `title Example Tattler
participant Dad
Brother 1->Brother 2:secret
note left of Sister:*eavesdrop*
Sister->Dad:tattle`
	sd, err := sequencediagram.ParseFromText(text)
	if err != nil {
		log.Fatalf("error parsing sequence diagram: %v", err)
	}

	if _, err := io.Copy(os.Stdout, textdiagram.Encode(sd)); err != nil {
		log.Fatalf("error copying text diagram to stdout: %v", err)
	}
	//                  Example Tattler
	//
	// ┌─────┐┌───────────┐  ┌───────────┐      ┌────────┐
	// │ Dad ││ Brother 1 │  │ Brother 2 │      │ Sister │
	// └─────┘└───────────┘  └───────────┘      └────────┘
	//    │         │  ┌────────┐  │                 │
	//               ──┤ secret ├─▶
	//    │         │  └────────┘  │                 │
	//                               ┌─────────────╗
	//    │         │              │ │ *eavesdrop* │ │
	//                               └─────────────┘
	//    │                              ┌────────┐  │
	//     ◀─────────────────────────────┤ tattle ├──
	//    │                              └────────┘  │
	// ┌─────┐┌───────────┐  ┌───────────┐      ┌────────┐
	// │ Dad ││ Brother 1 │  │ Brother 2 │      │ Sister │
	// └─────┘└───────────┘  └───────────┘      └────────┘
}
