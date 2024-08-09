/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./view/*/*.templ",
    "./**/*.templ",
    "./**/*.html",
    "./**/*.templ",
    "./**/*.go",
  ],
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["retro"],
  },
};
