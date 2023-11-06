import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';

const config = {
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			pages: 'static',
			assets: 'static',
			strict: false,
			trailingSlash: 'never',
			fallback: 'index.html',
		}),
	},
};

export default config;
