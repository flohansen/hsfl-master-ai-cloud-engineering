package loadtest

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/benchmark/config"
)

type LoadTest struct {
	config *config.LoadTestConfig
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
	log.Println("Starting load test")

	targets, err := ParseURLs(l.config.Targets)

	if err != nil {
		log.Println("Error parsing targets: ", err)
		return
	}

	if len(targets) == 0 {
		log.Println("No targets to test")
		return
	}

	users := l.config.Users
	specs := l.config.Specs
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

					url := PickRandomURL(targets)

					code, err := DoRequest(url)

					log.Println("User ", id, " made request with status code ", code, " and error ", err)

					time.Sleep(sleepDuration - time.Since(t))
				}
			}
		}(i)
	}

	wg.Wait()
}

func ParseURLs(urls []string) ([]*url.URL, error) {
	var parsedURLs []*url.URL

	for _, u := range urls {
		parsedURL, err := url.Parse(u)

		if err != nil {
			return nil, err
		}

		parsedURLs = append(parsedURLs, parsedURL)
	}

	return parsedURLs, nil
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

func PickRandomURL(targets []*url.URL) *url.URL {
	return targets[rand.Intn(len(targets))]
}
