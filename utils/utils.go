package utils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// exponentially increasing retry
// expRetry is a function that implements exponential backoff retry logic.
// It takes a function `fn` that returns a result and an error, the number of `tries` to attempt,
// and the `base` value for exponential backoff calculation.
// It retries the function `fn` for the specified number of times, with increasing delays between retries,
// until either the function succeeds or the maximum number of tries is reached.
// If the function succeeds, it returns the result and a nil error.
// If the function fails after all retries, it returns the last result and the last error encountered.
// TODO: fix logging here
func ExpRetry(fn func() (interface{}, error), tries int, base float64) (interface{}, error) {
	var err error
	var result interface{}
	totalNonceRetries := 1
	backoff := base

	for retry := 0; retry < tries; retry++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		if strings.Contains(err.Error(), "solution already submitted") {
			return result, err
		} else if strings.Contains(err.Error(), "nonce too low") {
			log.Printf("Error: %v, retries: %d", err, totalNonceRetries)
			time.Sleep(time.Duration(100) * time.Millisecond)
			tries += 1
			totalNonceRetries += 1
			if totalNonceRetries > 20 {
				break
			}
			continue
			// } else if strings.Contains(err.Error(), "non existent commitment") { // non existent commitment means no commitment was found for this task yet

			// 	//continue
			// } else if strings.Contains(err.Error(), "commitment must be in past") { // we're submitting swithin 1 block from commitment
			// 	//log.Printf("Error: %v", err)
			// 	//continue
		}
		//non existent commitment
		//seconds := math.Pow(base, float64(retry))

		sleepDuration := time.Duration(backoff) * time.Millisecond

		log.Printf("retry request failed, retrying in %s seconds", sleepDuration)
		log.Printf("Error: %v", err)
		//time.Sleep(time.Duration(seconds * float64(time.Second)))

		time.Sleep(sleepDuration)
		backoff *= 1.5 // Double the backoff time
	}

	log.Printf("retry request failed %d times", tries)
	return result, err
}

func ExpRetryWithNonce(fn func(nonce uint64) (interface{}, error), tries int, base, backoffMultiplier float64) (interface{}, error) {
	return ExpRetryWithNonceContext(context.Background(), fn, tries, base, backoffMultiplier)
}

// As ExpRetry, but with a nonce handling:
// if we get "nonce too low" error we try to extract the expected nonce from the error message and retry with it
func ExpRetryWithNonceContext(ctx context.Context, fn func(nonce uint64) (interface{}, error), tries int, base, backoffMultiplier float64) (interface{}, error) {
	var err error
	var result interface{}
	totalNonceRetries := 1
	backoff := base

	nonce := uint64(0)

	for retry := 0; retry < tries; retry++ {
		result, err = fn(nonce)
		if err == nil {
			return result, nil
		}
		if ctx.Err() != nil {
			log.Printf("CONTEXT CANCELLED OR ERRORED: %s", ctx.Err().Error())
			return result, ctx.Err()
		}
		if strings.Contains(err.Error(), "solution already submitted") {
			return result, err
		} else if strings.Contains(err.Error(), "nonce too low") || strings.Contains(err.Error(), "nonce too high") {

			// on nova error msg about nonce too high is like this:
			//nonce too high: address 0xF141fBA5aaf8688724F29DfB2bBC6EE244537328, tx: 693061 state: 693059 693059-693061=-2
			// on nova error msg about nonce too low is like this:
			///nonce too low: address 0x6c3Db6ef57735B8b62D0bdDa32c94389933d2f5d, tx: 316308 state: 316309 316309-316308=1
			parts := strings.Split(err.Error(), "state: ")
			if len(parts) < 2 {
				log.Println("state not found in error message")
			} else {

				state, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err != nil {
					log.Println(err.Error())
				} else {
					fmt.Println("setting new nonce to:", state)
					nonce = uint64(state)
				}
			}
			duration := time.Duration(rand.Intn(30)) * time.Millisecond
			//duration := time.Duration(300+rand.Intn(250)+25*totalNonceRetries) * time.Millisecond
			time.Sleep(duration)
			log.Printf("Error: %v, retries: %d, sleep: %s", err, totalNonceRetries, duration)
			tries++
			totalNonceRetries++
			if totalNonceRetries > 25 {
				break
			}
			continue
			// } else if strings.Contains(err.Error(), "non existent commitment") {
			// 	// non existent commitment means no commitment was found for this task yet
			// 	log.Printf("Error: %v", err)
			// 	time.Sleep(time.Duration(100) * time.Millisecond)
			// 	continue
			// } else if strings.Contains(err.Error(), "commitment must be in past") {
			// 	// we're submitting a solution within 1 block from commitment
			// 	log.Printf("Error: %v", err)
			// 	continue
		}
		// seconds := math.Pow(base, float64(retry))
		// log.Printf("retry request failed, retrying in %f seconds", seconds)
		// log.Printf("Error: %v", err)
		// time.Sleep(time.Duration(seconds * float64(time.Second)))
		sleepDuration := time.Duration(backoff) * time.Millisecond

		log.Printf("retry request failed, retrying in %s", sleepDuration)
		log.Printf("Error: %v", err)
		//time.Sleep(time.Duration(seconds * float64(time.Second)))

		time.Sleep(sleepDuration)
		backoff *= backoffMultiplier // Double the backoff time
	}

	log.Printf("retry request failed %d times", tries)
	return result, err
}

func Map[T, U any](data []T, f func(T) U) []U {

	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

// var cache *lru.Cache

// func init() {
// 	var err error
// 	cache, err = lru.New(128) // Create a cache with a maximum size of 128 items
// 	if err != nil {
// 		panic(err)
// 	}
// }

//var nonceMu = sync.Mutex{}
