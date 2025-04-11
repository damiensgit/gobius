package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
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
func ExpRetry(logger zerolog.Logger, fn func() (any, error), tries int, base float64) (any, error) {
	var err error
	var result any
	totalNonceRetries := 1
	backoff := base

	for range tries {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		if strings.Contains(err.Error(), "solution already submitted") {
			return result, err
		} else if strings.Contains(err.Error(), "nonce too low") {
			logger.Error().Err(err).Int("retries", totalNonceRetries).Msg("nonce too low")
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

		logger.Warn().Dur("sleep_duration", sleepDuration).Err(err).Msg("retry request failed, retrying")

		//time.Sleep(time.Duration(seconds * float64(time.Second)))

		time.Sleep(sleepDuration)
		backoff *= 1.5 // Double the backoff time
	}

	logger.Error().Int("tries", tries).Msg("retry request failed after multiple attempts")
	return result, err
}

// ExpRetryWithContext is like ExpRetry but checks for context cancellation during the retry loop.
func ExpRetryWithContext(ctx context.Context, logger zerolog.Logger, fn func() (any, error), tries int, base float64) (any, error) {
	var err error
	var result any
	totalNonceRetries := 1 // Note: Nonce logic is copied but might be irrelevant here if only used by ExpRetryWithNonceContext
	backoff := base

	for range tries {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		// Basic error checking copied from original ExpRetry
		if strings.Contains(err.Error(), "solution already submitted") {
			return result, err
		} else if strings.Contains(err.Error(), "nonce too low") {
			logger.Error().Err(err).Int("retries", totalNonceRetries).Msg("nonce too low")
			time.Sleep(time.Duration(100) * time.Millisecond)
			tries += 1
			totalNonceRetries += 1
			if totalNonceRetries > 20 {
				break
			}
			continue
		}

		// Check if context was cancelled before sleeping
		if ctxErr := ctx.Err(); ctxErr != nil {
			logger.Warn().Err(ctxErr).Msg("context cancelled while retrying")
			// Return the context error preferentially, but include the latest function error for context
			return result, fmt.Errorf("context cancelled: %w (last error: %v)", ctxErr, err)
		}

		sleepDuration := time.Duration(backoff) * time.Millisecond
		logger.Warn().Dur("sleep_duration", sleepDuration).Err(err).Msg("retry request failed, retrying")
		time.Sleep(sleepDuration)
		backoff *= 1.5 // Double the backoff time
	}

	logger.Error().Int("tries", tries).Msg("retry request failed after multiple attempts")
	return result, err
}

func ExpRetryWithNonce(logger zerolog.Logger, fn func(nonce uint64) (any, error), tries int, base, backoffMultiplier float64) (any, error) {
	return ExpRetryWithNonceContext(context.Background(), logger, fn, tries, base, backoffMultiplier)
}

// As ExpRetry, but with a nonce handling:
// if we get "nonce too low" error we try to extract the expected nonce from the error message and retry with it
func ExpRetryWithNonceContext(ctx context.Context, logger zerolog.Logger, fn func(nonce uint64) (any, error), tries int, base, backoffMultiplier float64) (any, error) {
	var err error
	var result any
	totalNonceRetries := 1
	backoff := base

	nonce := uint64(0)

	for range tries {
		result, err = fn(nonce)
		if err == nil {
			return result, nil
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
				logger.Warn().Msg("state not found in error message for nonce adjustment")
			} else {

				stateStr := strings.Fields(parts[1])[0]
				state, stateErr := strconv.Atoi(strings.TrimSpace(stateStr))
				if stateErr != nil {
					logger.Error().Err(stateErr).Str("state_part", stateStr).Msg("failed to parse state for nonce adjustment")
				} else {
					logger.Info().Int("new_nonce", state).Msg("setting new nonce")
					nonce = uint64(state)
				}
			}
			duration := time.Duration(rand.Intn(30)) * time.Millisecond
			//duration := time.Duration(300+rand.Intn(250)+25*totalNonceRetries) * time.Millisecond
			time.Sleep(duration)
			logger.Warn().Err(err).Int("retries", totalNonceRetries).Dur("sleep_duration", duration).Msg("Nonce error, retrying")
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
		sleepDuration := time.Duration(backoff) * time.Millisecond

		logger.Warn().Err(err).Dur("sleep_duration", sleepDuration).Msg("retry request failed, retrying")

		time.Sleep(sleepDuration)
		backoff *= backoffMultiplier // Double the backoff time
	}

	logger.Error().Int("tries", tries).Msg("retry request failed after multiple attempts")
	return result, err
}

func Map[T, U any](data []T, f func(T) U) []U {

	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func ZerologConsoleFormatTimestamp(timeFormat string) zerolog.Formatter {
	return func(i interface{}) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.Parse(zerolog.TimeFieldFormat, tt)
			if err != nil {
				t = tt
			} else {
				t = ts.Format(timeFormat)
			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64 = i, 0
				switch zerolog.TimeFieldFormat {
				case zerolog.TimeFormatUnixMs:
					nsec = int64(time.Duration(i) * time.Millisecond)
					sec = 0
				case zerolog.TimeFormatUnixMicro:
					nsec = int64(time.Duration(i) * time.Microsecond)
					sec = 0
				}
				ts := time.Unix(sec, nsec).UTC()
				t = ts.Format(timeFormat)
			}
		}
		return t
	}
}
