/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./cmd/web/**/*.html",
    "./cmd/web/**/*.templ",
    "./web/**/*templ",
    "./web/**/*.jsx",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark"],
  },
};
