import tailwindcss from '@tailwindcss/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  css: ['~/assets/css/tailwind.css'],
  modules: ['shadcn-nuxt', '@nuxt/icon'],

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: '',
    componentDir: './components/ui'
  },

  icon: {
    clientBundle: {
      scan: true,
    },
    customCollections: [
      {
        prefix: 'my-icon',
        dir: './assets/my-icons'
      },
    ],
  },

  runtimeConfig: {
    public: {
      serverUrl: ''
    }
  },

  app: {
    baseURL: '/home/',
    head: {
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/home/favicon.ico' }
      ]
    }
  }
})
