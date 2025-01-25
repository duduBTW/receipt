/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./renderer/**/*.templ"],
  theme: {
    extend: {
      colors: {
        background: {
          900: "hsl(var(--background-900))",
          800: "hsl(var(--background-800))",
          700: "hsl(var(--background-700))",
          600: "hsl(var(--background-600))",
        },
      },
      container: {
        screens: {
          DEFAULT: "672px",
        },
      },
      borderRadius: {
        lg: `var(--radius)`,
        md: `calc(var(--radius) - 2px)`,
        sm: "calc(var(--radius) - 4px)",
      },
    },
  },
};
