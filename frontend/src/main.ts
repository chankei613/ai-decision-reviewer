import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import 'virtual:uno.css'
import '@unocss/reset/tailwind.css'
import './assets/globals.css'
import App from './App.vue'
import InboxView from './pages/InboxView.vue'
import ResolvedView from './pages/ResolvedView.vue'
import DashboardView from './pages/DashboardView.vue'
import HelpView from './pages/HelpView.vue'
import SettingsView from './pages/SettingsView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/inbox' },
    { path: '/inbox', component: InboxView },
    { path: '/resolved', component: ResolvedView },
    { path: '/dashboard', component: DashboardView },
    { path: '/help', component: HelpView },
    { path: '/settings', component: SettingsView },
  ],
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
