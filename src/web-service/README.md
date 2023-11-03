# The web service component
[![Run tests (web service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-web-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-web-service.yml)

The web service component provides the frontend display of the application.

## Setup Frontend for development
1. Navigate in the `/frontend` folder of the web-service
2. Install dependencies: `npm ci`
3. For developing with hot module replacement use `npm run dev` and open up server provided by vite
4. For using webserver from golang use `npm run build`, run the `main.go` file and open localhost with port `:3000`.