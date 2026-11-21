import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./index.html', './src/**/*.{vue,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        alabaster: '#F7F7F3',
        sage: '#E8EFE6',
        'sage-active': '#DCE8D9',
        forest: '#1A241B',
        'forest-elevated': '#253227',
        charcoal: '#141A15',
        'slate-sage': '#5A655C',
        emerald: {
          brand: '#4ADE80',
          banner: '#34D361',
          accent: '#22C55E',
        },
        border: {
          light: '#E2E6E0',
          dark: '#2C3A2E',
        },
      },
      fontFamily: {
        sans: ['"Plus Jakarta Sans"', '"Inter"', 'system-ui', 'sans-serif'],
        body: ['"Inter"', 'system-ui', 'sans-serif'],
        serif: ['"Instrument Serif"', '"Playfair Display"', 'Georgia', 'serif'],
      },
      keyframes: {
        marquee: {
          '0%': { transform: 'translateX(0%)' },
          '100%': { transform: 'translateX(-50%)' },
        },
        'ken-burns': {
          '0%': { transform: 'scale(1) translate3d(0, 0, 0)' },
          '50%': { transform: 'scale(1.06) translate3d(-0.8%, -0.5%, 0)' },
          '100%': { transform: 'scale(1) translate3d(0, 0, 0)' },
        },
      },
      animation: {
        marquee: 'marquee 28s linear infinite',
        'ken-burns': 'ken-burns 22s ease-in-out infinite',
      },
    },
  },
  plugins: [],
};

export default config;
