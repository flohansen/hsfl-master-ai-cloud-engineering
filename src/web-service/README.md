# The web service component
[![Run tests (web service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-web-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-web-service.yml)

The web service component provides the frontend display of the application.

## Frontend Technology Stack
The frontend of this project is built on technologies, ensuring efficient development and a polished user interface.

### JavaScript Framework: [SvelteKit](https://kit.svelte.dev/)
We leverage the power of SvelteKit, a modern JavaScript framework with TypeScript support.
SvelteKit, in addition to its role as a framework, handles routing and takes care of static HTML generation.

### CSS Framework: [Tailwind CSS](https://tailwindcss.com/)
Styling is handled with Tailwind CSS, a utility-first CSS framework. Tailwind CSS allows us to create a visually appealing and responsive user interface.

### Build Tool: Vite [Vite](https://vitejs.dev/)
The project's static files are generated using Vite, a fast build tool that enhances the development experience. Vite's speed and extensibility make it a valuable asset in our toolchain, ensuring efficient bundling and optimization.

## Setup Frontend for development
1. Navigate in the `/frontend` folder of the web-service
2. Install dependencies: `npm ci`
3. For developing with hot module replacement use `npm run dev` and open up server provided by vite
4. For using webserver from golang use `npm run build`, run the `main.go` file and open localhost with port `:8080` or use docker and call localhost.