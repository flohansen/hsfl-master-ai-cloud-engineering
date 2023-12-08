package loadtest

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/benchmark/config"
)

type LoadTest struct {
	config *config.LoadTestConfig
}

type responseTimeEntry struct {
	total time.Duration
	count int
}

const MinimumSleepDuration = 10 * time.Millisecond
const MaximumSleepDuration = 1 * time.Second
const DecreaseInterval = float64(MaximumSleepDuration - MinimumSleepDuration)

func NewLoadTest(config *config.LoadTestConfig) *LoadTest {
	return &LoadTest{
		config: config,
	}
}

func (l *LoadTest) Run() {
	rampupDuration := time.Duration(l.config.RampUp) * time.Second
	duration := time.Duration(l.config.Duration) * time.Second
	overallDuration := rampupDuration + duration
	targets := l.config.Targets
	users := l.config.Users

	log.Println("Starting load test with ", users, " users for ", overallDuration, " seconds")

	var wg sync.WaitGroup

	responseTimesByTarget := make(map[string]responseTimeEntry)
	var mu sync.Mutex

	for i := 0; i < users; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			stop := time.NewTimer(overallDuration)
			start := time.Now()

			for {
				select {
				case <-stop.C:
					return
				default:
					target := PickRandomTarget(targets)

					go func() {
						err := MakeRequest(target, id, responseTimesByTarget, &mu)
						if err != nil {
							log.Println("Error making request: ", err)
						}
					}()

					elapsed := time.Since(start)
					sleepDuration := CalculateSleepDuration(elapsed, rampupDuration)
					time.Sleep(sleepDuration)
				}
			}
		}(i)
	}

	wg.Wait()

	log.Println("Finished load test")

	PrintAverageResponseTimes(responseTimesByTarget)
}

func CalculateSleepDuration(elapsed time.Duration, rampupDuration time.Duration) time.Duration {
	if elapsed < rampupDuration {
		// linearly decrease sleep duration from maximum duration to minimumSleepDuration
		remainingPortion := 1.0 - (float64(elapsed) / float64(rampupDuration))
		return MinimumSleepDuration + time.Duration(remainingPortion*DecreaseInterval)
	} else {
		return MinimumSleepDuration
	}
}

func PrintAverageResponseTimes(responseTimesByTarget map[string]responseTimeEntry) {
	for target, data := range responseTimesByTarget {
		average := data.total / time.Duration(data.count)
		log.Println("Average response time for ", target, " was ", average)
	}
}

func PickRandomTarget(targets []string) string {
	return targets[rand.Intn(len(targets))]
}

func MakeRequest(target string, id int, responseTimesByTarget map[string]responseTimeEntry, mu *sync.Mutex) error {
	requestStart := time.Now()

	resp, err := http.Get(target)
	if err != nil {
		return err
	}

	resp.Body.Close()

	log.Println("User ", id, " made request to ", target, " with status code ", resp.StatusCode)

	UpdateResponseTimes(target, requestStart, responseTimesByTarget, mu)

	return nil
}

func UpdateResponseTimes(target string, requestStart time.Time, responseTimesByTarget map[string]responseTimeEntry, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	data := responseTimesByTarget[target]
	data.total += time.Since(requestStart)
	data.count++
	responseTimesByTarget[target] = data
}
