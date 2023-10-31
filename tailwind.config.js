const plugin = require("tailwindcss/plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.html", "./static/*.js"],
  theme: {
    extend: {},
  },
  plugins: [
    plugin(function ({ addComponents }) {
      const indicator = {
        ".indicator": {
          opacity: "0",
          height: "0",
        },
        ".indicator.htmx-request": {
          opacity: "1",
          height: "100%",
        },
      };
      addComponents(indicator);
    }),
  ],
};
