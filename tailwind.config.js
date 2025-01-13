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
        },
      },
      container: {
        screens: {
          DEFAULT: "620px",
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
