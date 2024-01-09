# Benchmarking

This Go application is used to load test the API. A config file is used to specify the number of users, the duration of the test and the endpoints to test. The application will then send requests to the specified endpoints for the specified duration.

## Usage

```bash
go run main.go -config <path-to-config-file>
```

## Example Config

```json
{
  "users": 1,
  "startSleep": "1000ms",
  "specs": [
    {
      "targetSleep": "100ms",
      "duration": "30s"
    },
    {
      "targetSleep": "100ms",
      "duration": "30s"
    },
    {
      "targetSleep": "1000ms",
      "duration": "30s"
    }
  ],
  "targets": [
    "http://localhost:3001/"
  ]
}
```

Users is the number of users to simulate. The startSleep is the initial sleep time between requests. The specs array contains the specifications for each user. Each user will send requests to the specified endpoints for the specified duration. The targetSleep time is the time between requests that should be reached at the end of the corresponding phase. Within this time, linear interpolation takes place between the last sleep time and the target sleep time. The targets array contains the endpoints to be tested.

Using the example configuration above, the application will start with 1 user and a sleep time of 1000ms. After 30 seconds the sleep time will be 100ms (500ms after 15 seconds). After another 30 seconds the sleep time will be 100ms again. As there is no difference between the last sleep time and the target sleep time, the sleep time remains constant during this phase. After a further 30 seconds the sleep time between requests is 1000ms. The total duration of the stress test is 90 seconds.


