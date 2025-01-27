function withOpacity(variableName) {
  return ({ opacityValue }) => {
    if (opacityValue !== undefined) {
      return `rgba(var(${variableName}), ${opacityValue})`
    }
    return `rgb(var(${variableName}))`
  }
}

module.exports = {
  content: [
    "./cmd/web/**/*.html", "./cmd/web/**/*.templ",
  ],
  theme: {
    extend: {
      textColor: {
        skin: {
          base: withOpacity('--color-text-base'),
          muted: withOpacity('--color-text-muted'),
          inverted: withOpacity('--color-text-inverted'),
        },
      },
      backgroundColor: {
        skin: {
          fill: withOpacity('--color-fill'),
          'fill-muted': withOpacity('--color-fill-muted'),
          'button-accent': withOpacity('--color-button-accent'),
          'button-accent-hover': withOpacity('--color-button-accent-hover'),
          'button-muted': withOpacity('--color-button-muted'),
          'button-muted-hover': withOpacity('--color-button-muted-hover'),
          'secondary-accent': withOpacity('--color-secondary-accent'),
          'secondary-accent-hover': withOpacity('--color-secondary-accent-hover'),
        },
      },
      borderColor: {
        skin: {
          base: withOpacity('--color-border'),
          muted: withOpacity('--color-border-muted'),
        },
      },
      gradientColorStops: {
        skin: {
          hue: withOpacity('--color-fill'),
          'secondary-accent': withOpacity('--color-secondary-accent'),
        },
      },
      boxShadow: {
        skin: {
          base: `0 1px 3px 0 rgba(var(--color-shadow))`,
          dark: `0 4px 6px -1px rgba(var(--color-shadow-dark))`,
        },
      },
      ringColor: {
        skin: {
          base: withOpacity('--color-border'),
          success: withOpacity('--color-success'),
          error: withOpacity('--color-error'),
          warning: withOpacity('--color-warning'),
        },
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
