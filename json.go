package swagger2

import (
	"encoding/json"
)

// LoadJson parses the incoming byte array as Swagger 2 JSON data
func LoadJson(in []byte) (*Swagger, error) {
	var s Swagger
	err := json.Unmarshal(in, &s)
	if err != nil {
		return nil, err
	}
	return &s, err
}

// Json serializes the Swagger 2 document into JSON format
func (s *Swagger) Json() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
