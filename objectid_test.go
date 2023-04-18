package objectid_test

import (
	"bytes"
	"encoding/json"
	"github.com/pkg-id/objectid"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	id1 := objectid.New()
	id2 := objectid.New()
	if id1 == id2 {
		t.Fatalf("expect id are differs for each generation")
	}
}

func TestID_Timestamp(t *testing.T) {
	epochs := time.Now().Unix()
	id1 := objectid.NewEpochs(epochs)
	id2 := objectid.NewEpochs(epochs)

	if id1 == id2 {
		t.Fatalf("expect id are differs for each generation")
	}

	if id1.Timestamp().Unix() != id2.Timestamp().Unix() {
		t.Errorf("expect the timestamp are equal")
	}
}

func TestID_Count(t *testing.T) {
	epochs := time.Now().Unix()
	id1 := objectid.NewEpochs(epochs)
	id2 := objectid.NewEpochs(epochs)

	if id1 == id2 {
		t.Fatalf("expect id are differs for each generation")
	}

	if id1.Count() != id2.Count()-1 {
		t.Errorf("expect the counter of the next ID is one apart from the counter of previous ID")
	}
}

func TestID_String_Decode(t *testing.T) {
	epochs := time.Now().Unix()
	id1 := objectid.NewEpochs(epochs)
	id2, err := objectid.Decode(id1.String())
	if err != nil {
		t.Fatalf("expect no error; got error %v", err)
	}

	if id1 != id2 {
		t.Errorf("expect the ids are equal")
	}
}

func TestDecode(t *testing.T) {
	id, err := objectid.Decode(strings.Repeat("a", 23))
	if err == nil {
		t.Errorf("expect error, since the length is less than 24")
	}

	if !id.IsZero() {
		t.Fatalf("expec id must be a zero value")
	}

	id, err = objectid.Decode(strings.Repeat("z", 24))
	if err == nil {
		t.Errorf("expect error, since the string is not a valid hex")
	}

	if !id.IsZero() {
		t.Fatalf("expec id must be a zero value")
	}
}

func TestID_MarshalText(t *testing.T) {
	// Test valid input
	id := objectid.New()
	expectedOutput := []byte(id.String())
	output, err := id.MarshalText()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(expectedOutput, output) {
		t.Errorf("Expected %v but got %v", expectedOutput, output)
	}

	// Test invalid input
	id = objectid.ID{}
	expectedOutput = []byte(id.String())
	output, err = id.MarshalText()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(expectedOutput, output) {
		t.Errorf("Expected %v but got %v", expectedOutput, output)
	}
}

func TestID_UnmarshalText(t *testing.T) {
	// Test valid input
	id := objectid.New()
	b := []byte(id.String())
	expectedOutput := &id
	err := expectedOutput.UnmarshalText(b)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if *expectedOutput != id {
		t.Errorf("Expected %v but got %v", id, *expectedOutput)
	}

	// Test invalid input
	b = []byte("invalid")
	expectedOutput = &id
	err = expectedOutput.UnmarshalText(b)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestID_MarshalJSON(t *testing.T) {
	// Test valid input
	id := objectid.New()
	expectedOutput, err := json.Marshal(id)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	output, err := id.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(expectedOutput, output) {
		t.Errorf("Expected %v but got %v", expectedOutput, output)
	}

	// Test invalid input
	id = objectid.ID{}
	expectedOutput, err = json.Marshal(id)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	output, err = id.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(expectedOutput, output) {
		t.Errorf("Expected %v but got %v", expectedOutput, output)
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	type Data struct {
		ID objectid.ID `json:"id"`
	}

	rawJSON := `{"id":"640c5fe5d243553cda8dde1b"}`

	var obj Data
	err := json.Unmarshal([]byte(rawJSON), &obj)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedID := "640c5fe5d243553cda8dde1b"
	if obj.ID.String() != expectedID {
		t.Errorf("Unexpected ID value: got %s, expected %s", obj.ID, expectedID)
	}
}

func TestID_Value(t *testing.T) {
	var id objectid.ID
	value, err := id.Value()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	vid := value.(string)
	if vid != id.String() {
		t.Errorf("Unexpected ID value: got %s, expected %s", vid, id)
	}
}

func TestID_Scan(t *testing.T) {

	t.Run("invalid src type", func(t *testing.T) {
		var id objectid.ID
		err := id.Scan(1)
		if err == nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if id != objectid.Nil {
			t.Errorf("Unexpected ID value: got %s, expected %s", id, objectid.Nil)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		var id objectid.ID
		err := id.Scan("xxx-xxx-xx")
		if err == nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if id != objectid.Nil {
			t.Errorf("Unexpected ID value: got %s, expected %s", id, objectid.Nil)
		}
	})

	t.Run("from bytes", func(t *testing.T) {
		var id objectid.ID

		sid := objectid.New()
		src := []byte(sid.String())
		err := id.Scan(src)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if id != sid {
			t.Errorf("Unexpected ID value: got %s, expected %s", id, sid)
		}
	})

	t.Run("from string", func(t *testing.T) {
		var id objectid.ID

		sid := objectid.New()
		err := id.Scan(sid.String())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if id != sid {
			t.Errorf("Unexpected ID value: got %s, expected %s", id, sid)
		}
	})
}
