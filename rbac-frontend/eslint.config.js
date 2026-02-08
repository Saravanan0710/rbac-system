export default {
  files: ['**/*.js', '**/*.jsx'],
  ignores: ['dist', 'node_modules'],
  languageOptions: {
    ecmaVersion: 2021,
    sourceType: 'module',
    globals: {
      console: 'readonly',
      process: 'readonly'
    }
  },
  rules: {
    'no-unused-vars': 'warn',
    'no-console': 'off'
  }
}
