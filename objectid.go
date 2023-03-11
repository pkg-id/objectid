// Package objectid provides a way to generate MongoDB-style ObjectID.
// ref: https://www.mongodb.com/docs/manual/reference/method/ObjectId
package objectid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	timestampSize = 4
	processSize   = 5
	counterSize   = 3
)

var counter Counter
var counterOnce sync.Once

var machineProcessID MachineProcessID
var machineProcessIDOnce sync.Once

// SetCounter sets the global counter value for generating ObjectIDs.
func SetCounter(c Counter) { counter = c }

// SetMachineAndProcessID sets the machine and process ID for generating ObjectIDs.
func SetMachineAndProcessID(pid MachineProcessID) { machineProcessID = pid }

func init() {
	counterOnce.Do(func() {
		c, err := NewSecureCounter(rand.Reader)
		if err != nil {
			panic(err)
		}
		SetCounter(c)
	})

	machineProcessIDOnce.Do(func() {
		pid, err := NewMachineProcessID(rand.Reader)
		if err != nil {
			panic(err)
		}
		SetMachineAndProcessID(pid)
	})
}

// ID is the implementation of MongoDB ObjectID.
type ID [timestampSize + processSize + counterSize]byte

// Nil is a zero value of the ID.
var Nil ID

// New generates a new ID using the current time epochs, machine and process ids, and the global counter value.
func New() ID {
	epochs := time.Now().Unix()
	return NewEpochs(epochs)
}

// NewEpochs same as New but with given epochs.
func NewEpochs(epochs int64) ID {
	var id ID
	binary.BigEndian.PutUint32(id[:timestampSize], uint32(epochs))
	copy(id[timestampSize:timestampSize+processSize], machineProcessID[:])
	putBigEndianUint24(id[timestampSize+processSize:], counter.Next())
	return id
}

// String returns a string representation of the ID.
func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

// Timestamp returns the timestamp portion of the ID as a time.Time object.
func (id ID) Timestamp() time.Time {
	epochs := binary.BigEndian.Uint32(id[:timestampSize])
	return time.Unix(int64(epochs), 0).UTC()
}

// Count returns the counter value portion of the ID.
func (id ID) Count() uint32 {
	counterBits := make([]byte, 4)
	copy(counterBits[1:], id[timestampSize+processSize:])
	return binary.BigEndian.Uint32(counterBits)
}

// IsZero returns true if the ObjectID is the Nil value.
func (id ID) IsZero() bool { return id == Nil }

// MarshalText implements the encoding.TextMarshaler interface.
// This is useful when using the ID as a map key during JSON marshalling.
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// This is useful when using the ID as a map key during JSON unmarshalling.
func (id *ID) UnmarshalText(b []byte) error {
	decoded, err := Decode(string(b))
	if err != nil {
		return err
	}
	*id = decoded
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (id *ID) UnmarshalJSON(data []byte) error {
	// remove the surrounding quotes from the JSON string
	str := strings.Trim(string(data), "\"")
	decoded, err := Decode(str)
	if err != nil {
		return err
	}
	*id = decoded
	return nil
}

// Decode decodes the string representation and returns the corresponding ID.
func Decode(s string) (ID, error) {
	if len(s) != 24 {
		return Nil, errors.New("length is not 24 bytes")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return Nil, fmt.Errorf("decode hex: %w", err)
	}

	var id ID
	copy(id[:], b)
	return id, nil
}

// Counter is the implementation of the counter in ObjectID.
type Counter uint32

// NewSecureCounter generates a new secure counter value for generating ID.
func NewSecureCounter(reader io.Reader) (Counter, error) {
	var buf [4]byte // ensure for 32-byte.
	_, err := io.ReadFull(reader, buf[:])
	if err != nil {
		return 0, fmt.Errorf("generate initial counter: %w", err)
	}
	n := binary.BigEndian.Uint32(buf[:])
	return Counter(n), nil
}

// Next returns the next value of the counter.
func (c *Counter) Next() uint32 { return atomic.AddUint32((*uint32)(c), 1) }

// MachineProcessID is the implementation of the machine and process id portion for the ObjectID.
type MachineProcessID [processSize]byte

// NewMachineProcessID generates a new machine and process ids for generating ID.
func NewMachineProcessID(reader io.Reader) (MachineProcessID, error) {
	var process MachineProcessID
	_, err := io.ReadFull(reader, process[:])
	if err != nil {
		return process, fmt.Errorf("generate machine and process id: %w", err)
	}
	return process, nil
}

// putBigEndianUint24 converts an uint32 to a big-endian byte slice with 24 bits.
func putBigEndianUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
