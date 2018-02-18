# Sequence Diagram Generator

This package parses a string into a Sequence Diagram which can then be processed into a visual representation ([textdiagram](https://github.com/Laugusti/sequencediagram/tree/master/textdiagram))

## Usage

```go
sd, err := sequencediagram.ParseFromText(text)
if err != nil {
	log.Fatalf("error parsing sequence diagram: %v", err)
}
```

## Supported syntax
- Create a Title  
`title My Title`
- Define a Participant  
`participant Participant1`
- Message from A to B  
`A->B:Message`
- Message from A to B and response  
```
        A->B:Message
        B->A:Response
```
- Message from A to self  
`A->A:Message`
- Note  
`note right of A:Note`
