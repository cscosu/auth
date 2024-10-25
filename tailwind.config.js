// Note that we don't get type completion from this import because we don't have a package.json, we use the single binary tailwindcss distribution
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["templates/**/*.html.tpl"],
  theme: {
    extend: {},
  },
  plugins: [],
};
