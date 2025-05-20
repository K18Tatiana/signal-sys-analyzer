/**
 * plugins/vuetify.js
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'lightTheme',
    themes: {
      lightTheme: {
        dark: false,
        colors: {
          primary: '#1470AF', // $primary-color-light-mode
          secondary: '#053657', // $secondary-color-light-mode
          background: '#EAEAEE', // $background-color-light-mode
          surface: '#FFFFFF',
          error: '#FF5252',
          text: '#000000', // $text-color-light-mode
        },
        variables: {
          'font-family': 'Libre Baskerville, serif', // $font-primary
        },
      },
      darkTheme: {
        dark: true,
        colors: {
          primary: '#B7E2FF', // $primary-color-dark-mode
          secondary: '#8BD0FF', // $secondary-color-dark-mode
          background: '#343C41', // $background-color-dark-mode
          surface: '#1E1E1E',
          error: '#FF5252',
          text: '#ffffff', // $text-color-dark-mode
        },
        variables: {
          'font-family': 'Libre Baskerville, serif',
        },
      },
    },
  },
})
