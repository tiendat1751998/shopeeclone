//packages/design-tokens/src/colors.ts
export const colors = {
  primary: {
    50: '#e5f4ff',
    100: '#b8e0ff',
    200: '#8acbff',
    300: '#5cb6ff',
    400: '#2ea1ff',
    500: '#189eff',
    600: '#0a7ecc',
    700: '#005e99',
    800: '#003e66',
    900: '#001e33',
  },
  secondary: {
    50: '#fff3ef',
    100: '#ffd9cc',
    200: '#ffbfaa',
    300: '#ffa588',
    400: '#ff8b66',
    500: '#F15E2C',
    600: '#d44a1a',
    700: '#a83a14',
    800: '#7c2a0e',
    900: '#501a08',
  },
  success: { 50: '#ecfdf5', 500: '#22c55e', 600: '#16a34a', 700: '#15803d' },
  warning: { 50: '#fffbeb', 500: '#f59e0b', 600: '#d97706', 700: '#b45309' },
  danger: { 50: '#fef2f2', 500: '#ef4444', 600: '#dc2626', 700: '#b91c1c' },
  gray: {
    50: '#fafafa', 100: '#f5f5f5', 200: '#e8e8e8', 300: '#d4d4d4',
    400: '#a3a3a3', 500: '#757575', 600: '#525252', 700: '#404040',
    800: '#262626', 900: '#171717',
  },
  background: {
    primary: '#ffffff',
    secondary: '#f5f5f5',
    tertiary: '#fafafa',
    dark: '#0B1E33',
  },
  text: {
    primary: '#222222',
    secondary: '#757575',
    tertiary: '#a3a3a3',
    inverse: '#ffffff',
  },
  border: {
    default: '#e8e8e8',
    hover: '#d4d4d4',
    focus: '#189eff',
  },
} as const;

//packages/design-tokens/src/typography.ts
export const typography = {
  fontFamily: {
    sans: ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'Roboto', '"Helvetica Neue"', 'Arial', 'sans-serif'],
    mono: ['ui-monospace', 'SFMono-Regular', '"SF Mono"', 'Menlo', 'Consolas', 'monospace'],
  },
  fontSize: {
    '2xs': '10px',
    xs: '11px',
    sm: '13px',
    base: '14px',
    md: '15px',
    lg: '18px',
    xl: '22px',
    '2xl': '28px',
    '3xl': '36px',
  },
  fontWeight: {
    regular: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
    extrabold: 800,
  },
  lineHeight: {
    tight: 1.2,
    normal: 1.5,
    relaxed: 1.75,
  },
} as const;

//packages/design-tokens/src/spacing.ts
export const spacing = {
  0: '0px', 0.5: '2px', 1: '4px', 1.5: '6px', 2: '8px',
  3: '12px', 4: '16px', 5: '20px', 6: '24px', 8: '32px',
  10: '40px', 12: '48px', 16: '64px', 20: '80px', 24: '96px',
} as const;

//packages/design-tokens/src/shadows.ts
export const shadows = {
  sm: '0 1px 2px rgba(0,0,0,0.05)',
  md: '0 4px 6px rgba(0,0,0,0.07)',
  lg: '0 10px 15px rgba(0,0,0,0.1)',
  xl: '0 20px 25px rgba(0,0,0,0.1)',
} as const;

//packages/design-tokens/src/index.ts
export { colors } from './colors';
export { typography } from './typography';
export { spacing } from './spacing';
export { shadows } from './shadows';

export const breakpoints = {
  sm: '640px',
  md: '768px',
  lg: '1024px',
  xl: '1280px',
  '2xl': '1536px',
} as const;

export const borderRadius = {
  none: '0px',
  sm: '4px',
  md: '8px',
  lg: '12px',
  xl: '16px',
  full: '9999px',
} as const;

export const zIndex = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modal: 1040,
  popover: 1050,
  tooltip: 1060,
  toast: 1070,
} as const;

export const transition = {
  fast: '150ms ease',
  normal: '250ms ease',
  slow: '350ms ease',
} as const;
