import antfu from '@antfu/eslint-config'

export default antfu(
  {
    formatters: true,
    unocss: true,
    vue: true,
  },
  {
    rules: {
      'no-console': 'warn',
      'unused-imports/no-unused-vars': 'warn',
      'style/brace-style': ['warn', '1tbs', { allowSingleLine: true }],
    },
  },
)
