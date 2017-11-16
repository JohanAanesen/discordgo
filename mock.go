package discordgo

import (
	"io"
	"bytes"
	"encoding/json"
	"fmt"
)

func (s *Session) MockEvent(messageType int, message []byte) {

	var err error
	var reader io.Reader
	reader = bytes.NewBuffer(message)


	// Decode the event into an Event struct.
	var e *Event
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		fmt.Printf("error, %s", err)
		return
	}

	// Map event to registered event handlers and pass it along to any registered handlers.
	if eh, ok := registeredInterfaceProviders[e.Type]; ok {
	e.Struct = eh.New()

	// Attempt to unmarshal our event.
	if err = json.Unmarshal(e.RawData, e.Struct); err != nil {
	fmt.Printf("error unmarshalling %s event, %s", e.Type, err)
	}

	// Send event to any registered event handlers for it's type.
	// Because the above doesn't cancel this, in case of an error
	// the struct could be partially populated or at default values.
	// However, most errors are due to a single field and I feel
	// it's better to pass along what we received than nothing at all.
	// TODO: Think about that decision :)
	// Either way, READY events must fire, even with errors.
	s.handleEvent(e.Type, e.Struct)
	} else {
		fmt.Printf("unknown event: Op: %d, Seq: %d, Type: %s, Data: %s", e.Operation, e.Sequence, e.Type, string(e.RawData))
	}

	// For legacy reasons, we send the raw event also, this could be useful for handling unknown events.
	s.handleEvent(eventEventType, e)
}