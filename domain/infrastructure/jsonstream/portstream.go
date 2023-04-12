package jsonstream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Structure used to stream the Port data.
type PortStream struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	Key   string
	Error error
	Data  PortStream
}

// Stream helps transmit each streams within a channel.
type Stream struct {
	stream chan Entry
}

// NewPortStream returns a new `Stream` type.
func NewPortStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

// Watch watches JSON streams. Each stream entry will either have an error or a
// PortSteam object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) Watch() <-chan Entry {
	return s.stream
}

// Start starts streaming JSON file line by line. If an error occurs, the channel
// will be closed.
func (s Stream) Start(file io.Reader) {
	log.Println("Start Port Stream")

	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	decoder := json.NewDecoder(file)

	// Read opening delimiter. `{`
	openingDelimiter, err := decoder.Token()
	if err != nil {
		// #todo replace errors by struct in order to help the asserts at tests.
		errorMessage := fmt.Errorf("Error decoding opening delimiter: %w", err)
		log.Println(errorMessage)
		s.stream <- Entry{Error: errorMessage}

		return
	}

	if openingDelimiter != json.Delim('{') {
		errorMessage := fmt.Errorf("Opening delimiter is wrong. Expected { - Found %v", openingDelimiter)
		log.Println(errorMessage)
		s.stream <- Entry{
			Error: errorMessage,
		}

		return
	}

	log.Printf("Opening delimiter read %v.\n", openingDelimiter)

	// Read file content as long as there is something.
	line := 1

	for decoder.More() {
		// Reading key
		token, err := decoder.Token()

		if err != nil {
			errorMessage := fmt.Errorf("Error decoding key. Line %d - Error: %w", line, err)
			log.Println(errorMessage)
			s.stream <- Entry{Error: errorMessage}

			return
		}

		key, ok := token.(string)
		if !ok {
			errorMessage := fmt.Errorf("Error type asserting the key. Line %d - Error: %w", line, err)
			log.Println(errorMessage)
			s.stream <- Entry{Error: errorMessage}
		}

		log.Printf("Key %s decoded.\n", key)

		// Reading port
		var port PortStream
		if err := decoder.Decode(&port); err != nil {
			errorMessage := fmt.Errorf("Error decoding port. Key %v - Line %d - Error: %w", key, line, err)
			log.Println(errorMessage)
			s.stream <- Entry{
				Key:   key,
				Error: errorMessage}
		} else {
			s.stream <- Entry{
				Key:  key,
				Data: port,
			}
			log.Printf("Port ID'd by %s decoded.\n", key)
		}

		line++
	}

	// Read closing delimiter. `}`
	closingDelimiter, err := decoder.Token()
	if err != nil {
		errorMessage := fmt.Errorf("Error decoding closing delimiter: %w", err)
		log.Println(errorMessage)
		s.stream <- Entry{Error: errorMessage}

		return
	}

	log.Printf("Closing delimiter read %v.\n", closingDelimiter)
}
