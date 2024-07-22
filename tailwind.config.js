/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./view/*/*.templ",
    "./**/*.templ",
    "./**/*.html",
    "./**/*.templ",
    "./**/*.go",
  ],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["retro"],
  },
};
