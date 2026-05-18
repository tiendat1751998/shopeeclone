import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#ee4d2d', hover: '#d73211', light: '#fff0ed' },
        secondary: { DEFAULT: '#00bfa5', light: '#e0f7f5' },
        accent: '#f5a623',
        surface: '#ffffff',
        background: '#f5f5f5',
        text: { primary: '#222222', secondary: '#757575', disabled: '#bdbdbd' },
        border: '#e8e8e8',
        success: '#4caf50',
        error: '#f44336',
        warning: '#ff9800',
        info: '#2196f3',
      },
      borderRadius: { sm: '4px', md: '8px', lg: '12px', xl: '16px' },
      boxShadow: {
        sm: '0 1px 2px rgba(0,0,0,0.08)',
        md: '0 2px 8px rgba(0,0,0,0.12)',
        lg: '0 4px 16px rgba(0,0,0,0.16)',
      },
      maxWidth: { container: '1200px' },
      animation: { shimmer: 'shimmer 1.5s infinite', fadeIn: 'fadeIn 0.3s ease-in' },
      keyframes: {
        shimmer: { '0%': { backgroundPosition: '200% 0' }, '100%': { backgroundPosition: '-200% 0' } },
        fadeIn: { from: { opacity: '0' }, to: { opacity: '1' } },
      },
    },
  },
  plugins: [],
};

export default config;
