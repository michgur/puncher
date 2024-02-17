/** @type {import('tailwindcss').Config} */

const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ["../templates/**/*.html"],
  theme: {
    extend: {},
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
    })
  ],
  variants: {
    backgroundColor: ['responsive', 'dark', 'group-hover', 'focus-within', 'hover', 'focus', 'slot',]
  }
}
