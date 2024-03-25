const tailwindcss = require("tailwindcss");
const autoprefixer = require("autoprefixer");
const advancedVariables = require('postcss-advanced-variables')

const config = {
  plugins: [
    //Some plugins, like tailwindcss/nesting, need to run before Tailwind,
    tailwindcss(),
    //But others, like autoprefixer, need to run after,
    autoprefixer,
    advancedVariables(),
  ],
};

module.exports = config;
