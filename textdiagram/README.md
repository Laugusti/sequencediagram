# Text Sequence Diagram

Creates a textual representation of a sequence diagram

## Usage

```go
sd, err := sequencediagram.ParseFromText(text)
if err != nil {
	log.Fatalf("error parsing sequence diagram: %v", err)
}
r := textdiagram.Decode(sd)
td, err := ioutil.ReadAll(r)
if err != nil {
	log.Fatalf("error reading text diagram: %v", err)
}
fmt.Println(string(td))
```

## Example

Sequence Diagram
```
title Making a request
Client->Server:Request
Server->Server:Redirect
Server->Database:Query
Database-->Server:Result
note right of Server:Do Stuff
Server->Client:Response
```

Text Diagram
```
             Making a request             

┌────────┐       ┌────────┐    ┌──────────┐
│ Client │       │ Server │    │ Database │
└────────┘       └────────┘    └──────────┘
     │  ┌─────────┐   │              │
      ──┤ Request ├──▶
     │  └─────────┘   │              │
                       ────┐
     │                │    │Redirect │
                       ◀───┘
     │                │  ┌───────┐   │
                       ──┤ Query ├──▶
     │                │  └───────┘   │
                         ┌────────┐
     │                │◀-┤ Result ├--│
                         └────────┘
     │                │ ┌──────────╗ │
                        │ Do Stuff │
     │                │ └──────────┘ │
        ┌──────────┐
     │◀─┤ Response ├──│              │
        └──────────┘
     │                │              │
┌────────┐       ┌────────┐    ┌──────────┐
│ Client │       │ Server │    │ Database │
└────────┘       └────────┘    └──────────┘
```
