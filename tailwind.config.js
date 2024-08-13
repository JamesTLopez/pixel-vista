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
    themes: [
      "retro",
      {
        pixelvista: {
          primary: "#1EA0E6",
          secondary: "#1BCBE3",
          accent: "#B993E9",
          neutral: "#253D5B",
        },
      },
    ],
  },
};
