package common

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

// Number.MAX_SAFE_INT-15 to keep things easy
var bigM = new(big.Int).SetUint64(0x1FFFFFFFFFFFF0)

type TaskId [32]byte

func (t TaskId) String() string {
	return "0x" + hex.EncodeToString(t[:])
}

// taskid2Seed allows 2^64-1
// but for easy interop with existing tooling we reduce this
func (t TaskId) TaskId2Seed() uint64 {

	taskidInt := new(big.Int)

	//taskidInt.SetString(taskid, 10)
	taskidInt.SetBytes(t[:])

	return taskidInt.Mod(taskidInt, bigM).Uint64()
}

// space is the number of buckets to split the space into
// filter is the bucket to filter on
func (t TaskId) MatchFilter(filter int64, space uint64) bool {
	// If the filter is less than 0, always match
	if filter < 0 {
		return true
	}

	seed := t.TaskId2Seed()

	return seed%space == uint64(filter)

}

// space is the number of buckets to split the space into
// filter is the bucket to filter on
func (t TaskId) MatchFilterRange(filterFrom, fitlerTo int64, space uint64) bool {
	// If the filter is less than 0, always match
	if filterFrom < 0 {
		return true
	}

	seed := t.TaskId2Seed()

	modulus := seed % space

	return modulus >= uint64(filterFrom) && modulus <= uint64(fitlerTo)
}

// Convert a hexidecimal string to a TaskId
// Supports 0x prefixed and non-prefixed strings
// Returns an error if the input string is not a valid hexidecimal string
// or if the input string is not 32 bytes long
func ConvertTaskIdString2Bytes(hexString string) (TaskId, error) {

	// strip 0x prefix if present
	if hexString[:2] == "0x" {
		hexString = hexString[2:]
	}
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return TaskId{}, err
	}

	if len(decoded) != 32 {
		return TaskId{}, errors.New("invalid length, expected 32 bytes")
	}

	var arr [32]byte
	copy(arr[:], decoded)

	return TaskId(arr), nil
}

func (t *TaskId) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", src)
	}
	// strip 0x prefix if present
	if str[:2] == "0x" {
		str = str[2:]
	}
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}

	if len(bytes) != 32 {
		return errors.New("invalid length, expected 32 bytes")
	}

	copy(t[:], bytes)
	return nil
}

func (t TaskId) Value() (driver.Value, error) {
	return hex.EncodeToString(t[:]), nil
}
