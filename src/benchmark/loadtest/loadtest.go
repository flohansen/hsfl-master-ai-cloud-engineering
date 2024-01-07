package loadtest

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
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

func CalculateDuration(specs []config.Spec) time.Duration {
	var duration time.Duration
	for _, spec := range specs {
		duration += spec.Duration
	}
	return duration
}

func (l *LoadTest) Run() {
	log.Println("Starting load test")

	users := l.config.Users
	specs := l.config.Specs
	targets := l.config.Targets
	startSleep := l.config.StartSleep

	targetSleep := specs[0].Sleep
	overallDuration := CalculateOverallDuration(specs)

	var wg sync.WaitGroup

	for i := 0; i < users; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			stop := time.NewTimer(overallDuration)
			start := time.Now()

			lastStage := 0
			lastSleepDuration := startSleep
			targetSleepDuration := targetSleep

			for {
				select {
				case <-stop.C:
					return
				default:
					log.Print("Time since start: ", time.Since(start))

					t := time.Now()

					currentStage := GetCurrentStage(specs, time.Since(start))
					elapsedStageTime := CalculateStageElapsedTime(specs, currentStage, time.Since(start))
					currentStageDuration := specs[currentStage].Duration

					if currentStage != lastStage {
						lastSleepDuration = targetSleepDuration
						targetSleepDuration = specs[currentStage].Sleep

						lastStage = currentStage
					}

					sleepDuration := CalculateSleepDuration(elapsedStageTime, currentStageDuration, lastSleepDuration, targetSleepDuration)

					url, err := PickRandomURL(targets)

					if err != nil {
						log.Println("Error picking random URL: ", err)
						continue
					}

					code, err := DoRequest(url)

					log.Println("User ", id, " made request with status code ", code, " and error ", err)

					time.Sleep(sleepDuration - time.Since(t))
				}
			}
		}(i)
	}

	wg.Wait()
}

func DoRequest(target *url.URL) (uint64, error) {
	conn, err := net.Dial("tcp", target.Host)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	_, err = fmt.Fprintf(conn, "GET %s HTTP/1.1\r\nHost: go\r\n\r\n", target.Path)
	if err != nil {
		return 0, err
	}

	code := uint64(200)

	return code, nil
}

func GetCurrentStage(specs []config.Spec, elapsed time.Duration) int {
	var stage int
	var stageDuration time.Duration

	for i, spec := range specs {
		stageDuration += spec.Duration

		if elapsed < stageDuration {
			stage = i
			break
		}
	}
	return stage
}

func CalculateOverallDuration(specs []config.Spec) time.Duration {
	var duration time.Duration

	for _, spec := range specs {
		duration += spec.Duration
	}

	return duration
}

func CalculateStageElapsedTime(specs []config.Spec, currentStage int, elapsed time.Duration) time.Duration {
	var stageDuration time.Duration

	if currentStage == 0 {
		return elapsed
	}

	for i := 0; i < currentStage; i++ {
		stageDuration += specs[i].Duration
	}

	return elapsed - stageDuration
}

func CalculateSleepDuration(elapsed time.Duration, duration time.Duration, lastSleep time.Duration, targetSleep time.Duration) time.Duration {
	if lastSleep < targetSleep {
		interval := float64(targetSleep - lastSleep)
		remainingPortion := (float64(elapsed) / float64(duration))

		return lastSleep + time.Duration(remainingPortion*interval)
	} else {
		interval := float64(lastSleep - targetSleep)
		remainingPortion := 1.0 - (float64(elapsed) / float64(duration))

		return targetSleep + time.Duration(remainingPortion*interval)
	}
}

func PrintAverageResponseTimes(responseTimesByTarget map[string]responseTimeEntry) {
	for target, data := range responseTimesByTarget {
		average := data.total / time.Duration(data.count)
		log.Println("Average response time for ", target, " was ", average)
	}
}

func PickRandomURL(targets []string) (*url.URL, error) {
	url, err := url.Parse(targets[rand.Intn(len(targets))])

	if err != nil {
		return nil, err
	}

	return url, nil
}

func MakeRequest(httpClient *http.Client, target string, id int, responseTimesByTarget map[string]responseTimeEntry, mu *sync.Mutex) error {
	requestStart := time.Now()

	resp, err := httpClient.Get(target)
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
