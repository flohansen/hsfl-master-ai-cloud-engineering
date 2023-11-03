import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';

const config = {
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			pages: 'static',
			assets: 'static',
			strict: false,
		}),
		prerender: {
			entries: ['/shopping-lists/:id'],
		},
	},
};

export default config;
