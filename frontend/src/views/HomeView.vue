<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <div v-else ref="homeRoot" class="omnio-v2-page">
    <header :class="['ov2-header', { 'is-scrolled': headerScrolled }]">
      <nav class="ov2-nav ov2-shell" :aria-label="t('homeV2.primaryNavigation')">
        <router-link to="/home" class="ov2-brand" @click="mobileOpen = false">
          <span class="ov2-brand-mark"><OmnioMark :label="siteName" /></span>
          <strong>{{ siteName }}</strong>
        </router-link>

        <div class="ov2-nav-links">
          <a href="#network">{{ t('homeV2.models') }}</a>
          <a href="#use-cases">{{ t('homeV2.forEveryone') }}</a>
          <a href="#developers">{{ t('homeV2.forDevelopers') }}</a>
          <a href="#pricing">{{ t('homeV2.pricing') }}</a>
          <a :href="docUrl || '/docs/usage-guide/'">{{ t('home.docs') }}</a>
        </div>

        <div class="ov2-nav-actions">
          <LocaleSwitcher />
          <button
            type="button"
            class="ov2-icon-button"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="ov2-nav-login"
          >
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
          </router-link>
          <router-link
            v-if="!isAuthenticated"
            to="/register"
            class="ov2-nav-cta"
          >
            {{ t('homeV2.startNow') }}
          </router-link>
          <button
            type="button"
            class="ov2-mobile-trigger"
            :aria-expanded="mobileOpen"
            :aria-label="t('homeV2.toggleNavigation')"
            @click="mobileOpen = !mobileOpen"
          >
            <Icon :name="mobileOpen ? 'x' : 'menu'" size="sm" />
          </button>
        </div>
      </nav>

      <div v-if="mobileOpen" class="ov2-mobile-menu">
        <a href="#network" @click="mobileOpen = false">{{ t('homeV2.models') }}</a>
        <a href="#use-cases" @click="mobileOpen = false">{{ t('homeV2.forEveryone') }}</a>
        <a href="#developers" @click="mobileOpen = false">{{ t('homeV2.forDevelopers') }}</a>
        <a href="#pricing" @click="mobileOpen = false">{{ t('homeV2.pricing') }}</a>
        <a :href="docUrl || '/docs/usage-guide/'" @click="mobileOpen = false">{{ t('home.docs') }}</a>
        <router-link :to="isAuthenticated ? dashboardPath : '/login'" @click="mobileOpen = false">
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </header>

    <main>
      <section id="top" class="ov2-hero">
        <div class="ov2-hero-grid ov2-shell">
          <div class="ov2-hero-copy">
            <span class="ov2-kicker ov2-hero-kicker"><i></i>{{ t('homeV2.heroKicker') }}</span>
            <h1>
              {{ t('homeV2.heroLineOne') }}
              <span>{{ t('homeV2.heroLineTwo') }}</span>
            </h1>
            <p>{{ t('homeV2.heroDescription') }}</p>

            <div class="ov2-hero-actions">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/register'"
                class="ov2-button ov2-button-primary"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('homeV2.startUsing') }}
                <Icon name="arrowRight" size="xs" />
              </router-link>
              <a href="#network" class="ov2-button ov2-button-secondary">
                {{ t('homeV2.viewModels') }}
              </a>
            </div>

            <div class="ov2-hero-signals">
              <span><i></i>{{ t('homeV2.providersOnline') }}</span>
              <span>{{ t('homeV2.oneKey') }}</span>
              <span>{{ t('homeV2.transparentPricing') }}</span>
            </div>
          </div>

          <div class="ov2-hero-visual">
            <HeroNetwork :site-name="siteName" />
          </div>
        </div>

        <div class="ov2-provider-strip">
          <div class="ov2-shell">
            <span>{{ t('homeV2.connectedModels') }}</span>
            <div>
              <span><ModelIcon model="gpt-5" size="17px" />GPT</span>
              <i></i>
              <span><ModelIcon model="claude" size="17px" />Claude</span>
              <i></i>
              <span><ModelIcon model="gemini" size="17px" />Gemini</span>
              <i></i>
              <span><ModelIcon model="grok" size="17px" />Grok</span>
            </div>
          </div>
        </div>
      </section>

      <UseCaseShowcase />
      <ModelNetwork :site-name="siteName" />
      <DeveloperWorkspace
        :gateway-origin="gatewayOrigin"
        :doc-url="docUrl"
        :is-authenticated="isAuthenticated"
        :dashboard-path="dashboardPath"
      />
      <PricingPaths />

      <section class="ov2-final-cta" data-home-reveal>
        <div class="ov2-shell">
          <div class="ov2-final-cta-inner">
            <div class="ov2-final-portal" aria-hidden="true">
              <span class="ov2-final-portal-gate">
                <i></i><i></i>
                <OmnioMark :label="siteName" />
              </span>
              <span class="ov2-final-model-pulse is-gpt"><ModelIcon model="gpt-5" size="18px" /></span>
              <span class="ov2-final-model-pulse is-claude"><ModelIcon model="claude" size="18px" /></span>
              <span class="ov2-final-model-pulse is-gemini"><ModelIcon model="gemini" size="18px" /></span>
              <small>PORTAL // READY</small>
            </div>
            <div>
              <span class="ov2-kicker"><i></i>{{ t('homeV2.readyKicker') }}</span>
              <h2>{{ t('homeV2.readyTitle') }}</h2>
              <p>{{ t('homeV2.readyDescription') }}</p>
            </div>
            <router-link
              :to="isAuthenticated ? dashboardPath : '/register'"
              class="ov2-button ov2-button-light"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : t('homeV2.createAccount') }}
              <Icon name="arrowRight" size="xs" />
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <footer class="ov2-footer">
      <div class="ov2-shell ov2-footer-grid">
        <div>
          <router-link to="/home" class="ov2-brand">
            <span class="ov2-brand-mark"><OmnioMark :label="siteName" /></span>
            <strong>{{ siteName }}</strong>
          </router-link>
          <p>{{ t('homeV2.footerDescription') }}</p>
        </div>
        <div>
          <strong>{{ t('homeV2.product') }}</strong>
          <a href="#network">{{ t('homeV2.models') }}</a>
          <router-link to="/available-channels">{{ t('homeV2.modelStatus') }}</router-link>
          <a href="/docs/pricing/">{{ t('homeV2.pricing') }}</a>
        </div>
        <div>
          <strong>{{ t('homeV2.resources') }}</strong>
          <a :href="docUrl || '/docs/usage-guide/'">{{ t('home.docs') }}</a>
          <a href="/docs/about-omnio/">{{ t('homeV2.about') }}</a>
          <router-link to="/legal/user-agreement">{{ t('homeV2.legal') }}</router-link>
        </div>
      </div>
      <div class="ov2-shell ov2-footer-bottom">
        <span>© {{ currentYear }} {{ siteName }}</span>
        <span>{{ t('homeV2.footerNote') }}</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'
import Icon from '@/components/icons/Icon.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import OmnioMark from '@/components/home-v2/OmnioMark.vue'
import HeroNetwork from '@/components/home-v2/HeroNetwork.vue'
import ModelNetwork from '@/components/home-v2/ModelNetwork.vue'
import UseCaseShowcase from '@/components/home-v2/UseCaseShowcase.vue'
import DeveloperWorkspace from '@/components/home-v2/DeveloperWorkspace.vue'
import PricingPaths from '@/components/home-v2/PricingPaths.vue'
import '@/styles/home-v2.css'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const homeRoot = ref<HTMLElement | null>(null)
const isDark = ref(document.documentElement.classList.contains('dark'))
const headerScrolled = ref(false)
const mobileOpen = ref(false)
let revealObserver: IntersectionObserver | undefined

const siteName = computed(() => {
  const configuredName = appStore.cachedPublicSettings?.site_name?.trim()
  if (configuredName) return configuredName
  const injectedName = appStore.siteName?.trim()
  return injectedName && injectedName !== 'Sub2API' ? injectedName : 'Omnio'
})

const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isHomeContentUrl = computed(() => /^https?:\/\//i.test(homeContent.value.trim()))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => authStore.isAdmin ? '/admin/dashboard' : '/dashboard')
const gatewayOrigin = computed(() => typeof window === 'undefined' ? 'https://api.omnio.ai' : window.location.origin)
const currentYear = new Date().getFullYear()

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function updateHeader() {
  headerScrolled.value = window.scrollY > 20
}

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
  window.addEventListener('scroll', updateHeader, { passive: true })
  updateHeader()

  void nextTick(() => {
    if (!homeRoot.value) return
    const targets = Array.from(homeRoot.value.querySelectorAll<HTMLElement>('[data-home-reveal]'))
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    if (reducedMotion || !('IntersectionObserver' in window)) {
      targets.forEach((target) => target.classList.add('is-in-view'))
      return
    }

    homeRoot.value.classList.add('has-home-motion')
    revealObserver = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return
        entry.target.classList.add('is-in-view')
        revealObserver?.unobserve(entry.target)
      })
    }, {
      threshold: 0.12,
      rootMargin: '0px 0px -8% 0px'
    })

    targets.forEach((target) => revealObserver?.observe(target))
  })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateHeader)
  revealObserver?.disconnect()
})
</script>
