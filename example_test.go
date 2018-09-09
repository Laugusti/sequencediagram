package sequencediagram_test

import (
	"fmt"
	"log"

	"github.com/Laugusti/sequencediagram"
)

func ExampleParseFromText() {
	text := `title Example Tattler
participant Dad
Brother 1->Brother 2:secret
note left of Sister:*eavesdrop*
Sister->Dad:tattle`
	sd, err := sequencediagram.ParseFromText(text)
	if err != nil {
		log.Fatalf("error parsing sequence diagram: %v", err)
	}
	fmt.Println(sd)
	// Output:
	// title Example Tattler
	// participant Dad
	// Brother 1->Brother 2:secret
	// note left of Sister:*eavesdrop*
	// Sister->Dad:tattle
}
