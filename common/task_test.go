package common

import (
	"crypto/rand"
	"fmt"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestTaskId2Seed(t *testing.T) {

	taskid, _ := ConvertTaskIdString2Bytes("0x172c4a80e56a0610fa3a5a93062b9a537c6aca2043eb6b85b052f4cd62d69e2c")
	expected := uint64(810646377409212)

	result := taskid.TaskId2Seed()

	if result != expected {
		t.Errorf("TaskId2Seed() returned incorrect result, got: %d, want: %d", result, expected)
	}
}

func TestConvertTaskIdString2Bytes_InvalidString(t *testing.T) {

	_, err := ConvertTaskIdString2Bytes("Hello World!")

	assert.Error(t, err, "error expected")
}

func TestConvertTaskIdString2Bytes_Oddlength(t *testing.T) {

	_, err := ConvertTaskIdString2Bytes("dc6a147f2cd937a")

	assert.Error(t, err, "error expected")
}

func TestConvertTaskIdString2Bytes_Toolength(t *testing.T) {

	_, err := ConvertTaskIdString2Bytes("dc6a147f2cd937a1")

	assert.Error(t, err, "error expected")
}
func TestTaskIdString(t *testing.T) {

	taskid, err := ConvertTaskIdString2Bytes("1c6a147f2cd937a1b290b0dc4eff49a084ad72468fcb36df6d8ddb00c5ff6f7b")
	assert.NoError(t, err, "error unexpected")

	assert.Equal(t, "0x1c6a147f2cd937a1b290b0dc4eff49a084ad72468fcb36df6d8ddb00c5ff6f7b", taskid.String(), "unexpected value")
}

func TestMatchFilter(t *testing.T) {
	taskId := TaskId(common.HexToHash("dc6a147f2cd937a1b290b0dc4eff49a084ad72468fcb36df6d8ddb00c5ff6f7b").Bytes())
	filter := int64(5)
	space := uint64(10)

	expected := true
	result := taskId.MatchFilter(filter, space)

	if result != expected {
		t.Errorf("IsTaskId() returned incorrect result, got: %v, want: %v", result, expected)
	}
}

func TestMatchFilter_FilterLessThanZero(t *testing.T) {
	taskId := TaskId(common.HexToHash("dc6a147f2cd937a1b290b0dc4eff49a084ad72468fcb36df6d8ddb00c5ff6f7b").Bytes())
	filter := int64(-1)
	space := uint64(10)

	expected := true

	result := taskId.MatchFilter(filter, space)

	if result != expected {
		t.Errorf("IsTaskId() returned incorrect result, got: %v, want: %v", result, expected)
	}
}

func TestMatchFilter_FilterNotMatch(t *testing.T) {
	taskId := TaskId(common.HexToHash("dc6a147f2cd937a1b290b0dc4eff49a084ad72468fcb36df6d8ddb00c5ff6f7b").Bytes())
	filter := int64(3)
	space := uint64(10)

	expected := false
	result := taskId.MatchFilter(filter, space)

	if result != expected {
		t.Errorf("IsTaskId() returned incorrect result, got: %v, want: %v", result, expected)
	}
}

const distributionSize = 1000000
const bucketSize = 100

func TestMatchFilter_Distribution(t *testing.T) {

	bucketsMatched := make([]int, bucketSize)

	// for i := 0; i < bucketSize; i++ {
	// 	buckets[i] = 1 + i
	// }

	for i := 0; i < distributionSize; i++ {
		var taskId TaskId
		_, err := rand.Read(taskId[:])
		if err != nil {
			t.Fatalf("Failed to generate random TaskId: %v", err)
		}
		var i int64
		for i = 0; i < bucketSize; i++ {
			result := taskId.MatchFilter(i, bucketSize)
			if result {
				bucketsMatched[i] += 1
			}
		}

	}

	sum := 0
	// display distribution of the buckets
	for i := 0; i < bucketSize; i++ {
		fmt.Printf("Bucket %d matched %d times\n", i, bucketsMatched[i])
		sum += bucketsMatched[i]
	}

	if sum != distributionSize {
		t.Errorf("Distribution of buckets is incorrect, got: %d, want: %d", sum, distributionSize)
	}

	expected := float64(distributionSize) / float64(bucketSize)
	var chiSquareStat float64
	for i := 0; i < bucketSize; i++ {
		observed := float64(bucketsMatched[i])
		chiSquareStat += math.Pow(observed-expected, 2) / expected
	}

	// Degrees of freedom is bucketSize - 1
	df := bucketSize - 1
	// Compare chiSquareStat to a critical value from the Chi-Square distribution
	criticalValue := 123.225 // from https://www.medcalc.org/manual/chi-square-table.php using df = 99 and p = 0.05
	if chiSquareStat > criticalValue {
		t.Errorf("Distribution is not as expected (Chi-Square = %.2f, df = %d)", chiSquareStat, df)
	}
}
