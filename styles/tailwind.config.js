/** @type {import('tailwindcss').Config} */

const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ["../templates/**/*.html", "../app/templ/**/*.templ"],
  theme: {
    extend: {
      colors: {
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
    })
  ],
}
