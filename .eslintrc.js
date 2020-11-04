module.exports = {
  root: true,

  env: {
    node: true
  },

  rules: {
    semi: 'off',
    'no-console': 'off',
    'no-debugger': 'off',
    'generator-star-spacing': 'off',
    'eol-last': 'off',
    'object-curly-newline': 'off'
  },

  parserOptions: {
    parser: '@typescript-eslint/parser'
  },

  extends: ['plugin:vue/essential', '@vue/standard', '@vue/typescript']
}
