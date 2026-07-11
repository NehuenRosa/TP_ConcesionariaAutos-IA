/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      colors: {
        brand: {
          50: '#EEEDFE',
          100: '#DDD9FC',
          200: '#BBB4F9',
          300: '#9990F5',
          400: '#766BF2',
          500: '#534AB7',
          600: '#423A94',
          700: '#312B71',
          800: '#211D4E',
          900: '#100E2B',
        },
        surface: '#E8E9F0',
        'surface-text': '#4A4D66',
        'hero-bg': '#16182B',
        'text-primary': '#1B1D2A',
        'text-secondary': '#6B6E85',
        'text-placeholder': '#9C9FBB',
        'border-subtle': '#E4E5EE',
        'accent-light': '#EEEDFE',
        'accent-text': '#3C3489',
        'badge-nuevo-bg': '#E1F5EE',
        'badge-nuevo-text': '#085041',
        'badge-usado-bg': '#FAEEDA',
        'badge-usado-text': '#854F0B',
      },
      animation: {
        'fade-in': 'fadeIn 0.4s ease-out',
        'slide-up': 'slideUp 0.4s ease-out',
        'scale-in': 'scaleIn 0.25s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(16px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.96)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
      },
    },
  },
  plugins: [],
}
