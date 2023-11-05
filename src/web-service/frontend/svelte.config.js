import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/kit/vite';

const config = {
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			pages: 'static',
			assets: 'static',
			strict: false,
			trailingSlash: 'never',
		}),
	},
};

export default config;
