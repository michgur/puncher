/** @type {import('tailwindcss').Config} */

const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ["../templates/**/*.html", "../app/templ/**/*.templ"],
  theme: {
    extend: {
      colors: {
        'card': {
          '50': 'var(--card-50)',
          '100': 'var(--card-100)',
          '200': 'var(--card-200)',
          '300': 'var(--card-300)',
          '400': 'var(--card-400)',
          '500': 'var(--card-500)',
          '600': 'var(--card-600)',
          '700': 'var(--card-700)',
          '800': 'var(--card-800)',
          '900': 'var(--card-900)',
          '950': 'var(--card-950)',
        },
        'citron': {
          '50': '#fafbea',
          '100': '#f5f6c6',
          '200': '#e3e490',
          '300': '#d7d75b',
          '400': '#c8c832',
          '500': '#acac2f',
          '600': '#938d25',
          '700': '#746a20',
          '800': '#5f5721',
          '900': '#504921',
          '950': '#26220d',
        },
        'blueberry': {
          '50': '#f1f8fe',
          '100': '#e2f1fd',
          '200': '#a4d7f9',
          '300': '#90c7e9',
          '400': '#58aada',
          '500': '#318ec4',
          '600': '#1f71a8',
          '700': '#1a5b89',
          '800': '#1a4d70',
          '900': '#1d435e',
          '950': '#132a3e',
        },
        'gray': {
          '50': '#f9f9fb',
          '100': '#f3f2f7',
          '200': '#e9e8ed',
          '300': '#d6d4dd',
          '400': '#bcb9c6',
          '500': '#a19eae',
          '600': '#8b8797',
          '700': '#75727e',
          '800': '#625f67',
          '900': '#514f54',
          '950': '#34313a',
        },
        'peach': {
          '50': '#fdf5f3',
          '100': '#fbe9e5',
          '200': '#f6d0c7',
          '300': '#f1bcb0',
          '400': '#e79582',
          '500': '#da7259',
          '600': '#c5573d',
          '700': '#a54630',
          '800': '#893d2b',
          '900': '#733729',
          '950': '#3e1a11',
        },
        'brown': {
          '50': '#f5f3f1',
          '100': '#e6e0db',
          '200': '#cfc3b9',
          '300': '#b39f91',
          '400': '#9c8373',
          '500': '#8e7364',
          '600': '#795e55',
          '700': '#624b46',
          '800': '#54413f',
          '900': '#4a3a39',
          '950': '#291f1f',
        },
      },
    },
  },
  plugins: [
    plugin(({ addVariant, e }) => {
      function addSlotVariant(variant, selector) {
        addVariant(variant, ({ modifySelectors, separator }) => {
          modifySelectors(({ className }) => {
            return `.${variant}${e(`${separator}${className}`)}::slotted(*${selector || ''})`;
          })
        });
      }

      addSlotVariant('slot');
      addSlotVariant('slot-hover', ':hover');
    }),
    plugin(({ addVariant }) => {
      addVariant('slider-thumb', ['&::-webkit-slider-thumb', '&::-moz-range-thumb']);
      addVariant('slider-track', ['&::-webkit-slider-runnable-track', '&::-moz-range-track']);
    }),
    plugin(({ addUtilities, theme }) => {
      const colors = theme('colors');
      let newUtilities = {};

      Object.keys(colors).forEach(color => {
        if (color !== 'card' && typeof colors[color] === 'object') {
          newUtilities[`.card-${color}`] = Object.keys(colors[color]).reduce((acc, key) => {
            acc[`--card-${key}`] = colors[color][key];
            return acc;
          }, {});
        }
      });

      addUtilities(newUtilities, ['responsive', 'hover']);
    })
  ],
}
