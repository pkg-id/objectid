package objectid_test

import (
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
