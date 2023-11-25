export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    colors: {
      'blue': {
        'light': '#BFDAD9',
        'dark': '#305F5D',
      },
      'green': {
        'light': '#70AA67',
        'dark': '#317050',
      },
      'gray': {
        'light': '#F4F4F9',
        'dark': '#8F8F8F',
      },
      'red': '#BB3D41',
      'white': '#FFFFFF',
      'black': '#262626',
    },
    fontFamily: {
      'manrope': ['Manrope', 'sans-serif'],
    },
    fontSize: {
      'xxs': ['0.625rem', '1'], // 10px
      'xs': ['0.75rem', '1.2'], // 12px
      'sm': ['0.875rem', '1.4'], // 14px
      'base': ['1rem', '1.5'], // 16px
      'lg': ['1.125rem', '1.55'], // 18px
      'xl': ['1.25rem', '1.55'], // 20px
      '2xl': ['1.5rem', '1.35'], // 24px
      '3xl': ['1.75rem', '1.43'], // 28px
      '4xl': ['2.25rem', '1.33'], // 36px
      '5xl': ['3rem', '1.25'], // 48px
    },
    extend: {
      dropShadow: {
        'navBar': 'drop-shadow(0px -4px 100px rgba(38, 38, 38, 1))',
      }
    },
  },
  plugins: [],
}

