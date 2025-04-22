/** @type {require('tailwindcss').Config} */
module.exports = {
  content: ["./emails/**/*.tsx"],
  presets: [require("tailwindcss-preset-email")],
  important: false,
};
