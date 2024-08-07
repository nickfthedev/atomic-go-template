/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.templ", "./web/**/*.jsx", "./web/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark"],
  },
};
