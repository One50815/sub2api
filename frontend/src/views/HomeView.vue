<template>
  <div v-if="hasHomeContent" class="min-h-[100dvh]">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-[100dvh] w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-9 w-9 shrink-0 rounded-lg object-contain" />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer" class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500" :title="t('home.viewDocs')">
            <Icon name="book" size="md" />
          </a>
          <button type="button" class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500" :title="isDark ? t('home.switchToLight') : t('home.switchToDark')" @click="toggleTheme">
            <Icon :name="isDark ? 'sun' : 'moon'" size="md" />
          </button>
          <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white">
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>
    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img :src="siteLogo || '/logo.svg'" alt="Logo" class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain" />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-primary-600 px-5 py-2.5 text-sm font-medium text-white">
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>
    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 sm:px-6 dark:border-dark-800 dark:text-dark-400">&copy; {{ currentYear }} {{ siteName }}</footer>
  </div>

  <div v-else ref="homeRoot" :class="['omnio-home', { 'is-en': locale !== 'zh' }]">
    <header class="oh-header">
      <nav class="oh-nav oh-shell" :aria-label="t('homeFresh.primaryNavigation')">
        <router-link to="/home" class="oh-brand" @click="mobileOpen = false">
          <span class="oh-brand-mark">
            <img v-if="siteLogo" :src="siteLogo" :alt="siteName" />
            <OmnioMark v-else :label="siteName" />
          </span>
          <strong>{{ siteName }}</strong>
        </router-link>

        <div class="oh-nav-links">
          <a href="#workflow">{{ t('homeFresh.howItWorks') }}</a>
          <a href="#models">{{ t('homeFresh.models') }}</a>
          <a href="#benefits">{{ t('homeFresh.capabilities') }}</a>
          <a href="#faq">{{ t('homeFresh.faq') }}</a>
          <a :href="docUrl || '/docs/usage-guide/'">{{ t('home.docs') }}</a>
        </div>

        <div class="oh-nav-actions">
          <LocaleSwitcher class="oh-locale-switcher" />
          <button
            type="button"
            class="oh-nav-icon"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
          </button>
          <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="oh-nav-login">
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
          </router-link>
          <router-link
            v-if="!isAuthenticated"
            to="/register"
            class="oh-nav-primary"
          >
            {{ t('homeFresh.getStarted') }}
            <Icon name="arrowRight" size="xs" />
          </router-link>
          <button
            type="button"
            class="oh-menu-button"
            :aria-expanded="mobileOpen"
            :aria-label="t('homeFresh.toggleNavigation')"
            @click="mobileOpen = !mobileOpen"
          >
            <Icon :name="mobileOpen ? 'x' : 'menu'" size="sm" />
          </button>
        </div>
      </nav>

      <div v-if="mobileOpen" class="oh-mobile-menu">
        <a href="#workflow" @click="mobileOpen = false">{{ t('homeFresh.howItWorks') }}</a>
        <a href="#models" @click="mobileOpen = false">{{ t('homeFresh.models') }}</a>
        <a href="#benefits" @click="mobileOpen = false">{{ t('homeFresh.capabilities') }}</a>
        <a href="#faq" @click="mobileOpen = false">{{ t('homeFresh.faq') }}</a>
        <a :href="docUrl || '/docs/usage-guide/'" @click="mobileOpen = false">{{ t('home.docs') }}</a>
        <router-link :to="isAuthenticated ? dashboardPath : '/login'" @click="mobileOpen = false">
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
        <router-link v-if="!isAuthenticated" to="/register" class="is-primary" @click="mobileOpen = false">
          {{ t('homeFresh.getStarted') }}
        </router-link>
      </div>
    </header>

    <main>
      <section class="oh-hero">
        <div class="oh-shell oh-hero-stage">
          <div class="oh-hero-copy">
            <h1>
              <span>{{ t('homeFresh.heroLineOne') }}</span>
              <span>{{ t('homeFresh.heroLineTwo', { brand: siteName }) }}</span>
            </h1>
            <p>{{ t('homeFresh.heroDescription', { brand: siteName }) }}</p>
            <div class="oh-hero-actions">
              <router-link :to="isAuthenticated ? dashboardPath : '/register'" class="oh-button oh-button-dark">
                {{ isAuthenticated ? t('home.goToDashboard') : t('homeFresh.startFree') }}
                <Icon name="arrowRight" size="sm" />
              </router-link>
              <router-link to="/model-plaza" class="oh-button oh-button-light">
                {{ t('homeFresh.exploreModels') }}
                <Icon name="cube" size="sm" />
              </router-link>
            </div>
          </div>

          <aside class="oh-hero-proof" :aria-label="t('homeFresh.privacyProofLabel')">
            <div class="oh-proof-header">
              <span><Icon name="shield" size="md" /></span>
              <div>
                <small>{{ t('homeFresh.privacyProofLabel') }}</small>
                <strong>{{ t('homeFresh.privacyProofTitle') }}</strong>
              </div>
            </div>
            <div class="oh-proof-list">
              <div v-for="proof in heroProofs" :key="proof.label" class="oh-proof-item">
                <Icon :name="proof.icon" size="sm" />
                <div>
                  <small>{{ proof.label }}</small>
                  <strong>{{ proof.value }}</strong>
                </div>
              </div>
            </div>
          </aside>
        </div>

        <div class="oh-shell oh-launcher-wrap">
          <div class="oh-launcher-backdrop" aria-hidden="true"></div>
          <div class="oh-launcher">
            <div class="oh-launcher-topbar">
              <div class="oh-launcher-title">
                <span class="oh-product-mark"><OmnioMark :label="siteName" /></span>
                <div>
                  <strong>{{ t('homeFresh.launcherTitle', { brand: siteName }) }}</strong>
                  <small><i></i>{{ t('homeFresh.ready') }}</small>
                </div>
              </div>
              <span class="oh-launcher-account">{{ t('homeFresh.workspace') }}</span>
            </div>

            <div class="oh-launcher-body">
              <aside class="oh-task-nav" :aria-label="t('homeFresh.chooseTask')">
                <button
                  v-for="task in tasks"
                  :key="task.id"
                  type="button"
                  :class="{ 'is-active': selectedTask === task.id }"
                  @click="selectedTask = task.id"
                >
                  <Icon :name="task.icon" size="sm" />
                  <span>{{ task.label }}</span>
                  <Icon name="chevronRight" size="xs" />
                </button>
              </aside>

              <div class="oh-task-workspace">
                <div class="oh-workspace-heading">
                  <div>
                    <span>{{ t('homeFresh.currentTask') }}</span>
                    <h2>{{ activeTask.title }}</h2>
                  </div>
                  <span class="oh-model-chip">
                    <ModelIcon :model="activeTask.modelIcon" size="18px" />
                    {{ activeTask.model }}
                  </span>
                </div>

                <div class="oh-request-panel">
                  <span>{{ activeTask.promptLabel }}</span>
                  <p>{{ activeTask.prompt }}</p>
                  <router-link :to="activeTask.route" class="oh-run-task">
                    {{ activeTask.action }}
                    <Icon name="arrowRight" size="xs" />
                  </router-link>
                </div>

                <div class="oh-result-panel" aria-live="polite">
                  <div class="oh-result-heading">
                    <span class="oh-result-icon"><Icon :name="activeTask.resultIcon" size="sm" /></span>
                    <div>
                      <strong>{{ activeTask.resultTitle }}</strong>
                      <small>{{ activeTask.resultDescription }}</small>
                    </div>
                    <span class="oh-result-status"><i></i>{{ t('homeFresh.live') }}</span>
                  </div>
                  <div class="oh-result-facts">
                    <div v-for="fact in activeTask.facts" :key="fact.label">
                      <span>{{ fact.label }}</span>
                      <strong>{{ fact.value }}</strong>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="oh-model-ribbon" :aria-label="t('homeFresh.connectedTo')">
        <div class="oh-shell">
          <span>{{ t('homeFresh.connectedTo') }}</span>
          <div class="oh-model-list">
            <span v-for="model in models" :key="model.name">
              <ModelIcon :model="model.icon" size="22px" />
              {{ model.name }}
            </span>
          </div>
        </div>
      </section>

      <section id="workflow" class="oh-section oh-workflow">
        <div class="oh-shell oh-workflow-layout">
          <div class="oh-section-heading" data-home-reveal>
            <h2>{{ t('homeFresh.workflowTitle') }}</h2>
            <p>{{ t('homeFresh.workflowDescription', { brand: siteName }) }}</p>
          </div>

          <div class="oh-steps">
            <article class="oh-step oh-step-task" data-home-reveal>
              <div class="oh-step-copy">
                <h3>{{ t('homeFresh.stepOneTitle') }}</h3>
                <p>{{ t('homeFresh.stepOneDescription') }}</p>
                <router-link to="/model-plaza">
                  {{ t('homeFresh.seeAllTasks') }}
                  <Icon name="arrowRight" size="xs" />
                </router-link>
              </div>
              <div class="oh-task-board">
                <button
                  v-for="task in tasks"
                  :key="`flow-${task.id}`"
                  type="button"
                  :class="{ 'is-active': selectedTask === task.id }"
                  @click="selectedTask = task.id"
                >
                  <span><Icon :name="task.icon" size="md" /></span>
                  <strong>{{ task.label }}</strong>
                  <small>{{ task.short }}</small>
                </button>
              </div>
            </article>

            <article class="oh-step oh-step-connect" data-home-reveal>
              <div class="oh-code-window terminal-container">
                <div class="oh-code-toolbar">
                  <span>{{ t('homeFresh.sdkExample') }}</span>
                  <button type="button" :title="t('homeFresh.copyCode')" @click="copyCode">
                    <Icon :name="codeCopied ? 'check' : 'copy'" size="sm" />
                    {{ codeCopied ? t('homeFresh.copied') : t('homeFresh.copy') }}
                  </button>
                </div>
                <pre><code>{{ codeSample }}</code></pre>
                <div class="oh-code-footer">
                  <span><i></i>{{ t('homeFresh.endpointReady') }}</span>
                  <strong>{{ gatewayOrigin }}</strong>
                </div>
              </div>
              <div class="oh-step-copy">
                <h3>{{ t('homeFresh.stepTwoTitle') }}</h3>
                <p>{{ t('homeFresh.stepTwoDescription') }}</p>
                <router-link :to="isAuthenticated ? '/keys' : '/register'">
                  {{ t('homeFresh.createKey') }}
                  <Icon name="arrowRight" size="xs" />
                </router-link>
              </div>
            </article>

            <article class="oh-step oh-step-observe" data-home-reveal>
              <div class="oh-step-copy">
                <h3>{{ t('homeFresh.stepThreeTitle') }}</h3>
                <p>{{ t('homeFresh.stepThreeDescription') }}</p>
                <router-link :to="isAuthenticated ? '/usage' : '/register'">
                  {{ t('homeFresh.viewUsage') }}
                  <Icon name="arrowRight" size="xs" />
                </router-link>
              </div>
              <div class="oh-usage-surface">
                <div class="oh-usage-summary">
                  <div>
                    <span>{{ t('homeFresh.todaySpend') }}</span>
                    <strong>$4.28</strong>
                  </div>
                  <div>
                    <span>{{ t('homeFresh.successRate') }}</span>
                    <strong>99.8%</strong>
                  </div>
                  <span class="oh-date-chip">{{ t('homeFresh.exampleData') }}</span>
                </div>
                <div class="oh-request-list">
                  <div v-for="request in requestRows" :key="request.model">
                    <span class="oh-request-model"><ModelIcon :model="request.icon" size="19px" />{{ request.model }}</span>
                    <span>{{ request.tokens }}</span>
                    <span>{{ request.latency }}</span>
                    <strong>{{ request.cost }}</strong>
                  </div>
                </div>
                <div class="oh-usage-caption">
                  <Icon name="chart" size="sm" />
                  {{ t('homeFresh.requestLevelVisibility') }}
                </div>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section id="models" class="oh-section oh-model-section">
        <div class="oh-shell oh-model-layout">
          <div class="oh-model-copy" data-home-reveal>
            <span>{{ t('homeFresh.modelKicker') }}</span>
            <h2>{{ t('homeFresh.modelTitle') }}</h2>
            <p>{{ t('homeFresh.modelDescription') }}</p>
            <router-link to="/model-plaza" class="oh-button oh-button-dark">
              {{ t('homeFresh.openModelPlaza') }}
              <Icon name="arrowRight" size="sm" />
            </router-link>
          </div>

          <div class="oh-model-directory" data-home-reveal>
            <div v-for="model in modelDirectory" :key="model.name" class="oh-model-row">
              <span class="oh-directory-icon"><ModelIcon :model="model.icon" size="25px" /></span>
              <div>
                <strong>{{ model.name }}</strong>
                <small>{{ model.capability }}</small>
              </div>
              <span class="oh-directory-status"><i></i>{{ t('homeFresh.available') }}</span>
              <Icon name="chevronRight" size="sm" />
            </div>
          </div>
        </div>
      </section>

      <section id="benefits" class="oh-section oh-benefits">
        <div class="oh-shell">
          <div class="oh-section-heading is-left" data-home-reveal>
            <h2>{{ t('homeFresh.capabilityTitle') }}</h2>
            <p>{{ t('homeFresh.capabilityDescription', { brand: siteName }) }}</p>
          </div>

          <div class="oh-benefit-grid">
            <article class="oh-benefit-main" data-home-reveal>
              <div class="oh-benefit-visual">
                <div class="oh-policy-heading">
                  <span><Icon name="shield" size="sm" />{{ t('homeFresh.routingPolicy') }}</span>
                  <strong>{{ t('homeFresh.active') }}</strong>
                </div>
                <div class="oh-policy-route">
                  <span>
                    <small>{{ t('homeFresh.requestedModel') }}</small>
                    <strong><ModelIcon model="gpt-5" size="22px" />GPT-5</strong>
                  </span>
                  <Icon name="arrowRight" size="sm" />
                  <span>
                    <small>{{ t('homeFresh.resolvedModel') }}</small>
                    <strong><ModelIcon model="gpt-5" size="22px" />GPT-5</strong>
                  </span>
                </div>
                <div class="oh-policy-list">
                  <span><Icon name="checkCircle" size="sm" />{{ t('homeFresh.policyHealth') }}</span>
                  <span><Icon name="checkCircle" size="sm" />{{ t('homeFresh.policyBudget') }}</span>
                  <span><Icon name="checkCircle" size="sm" />{{ t('homeFresh.policyFallback') }}</span>
                </div>
              </div>
              <div>
                <span class="oh-benefit-icon"><Icon name="sync" size="md" /></span>
                <h3>{{ t('homeFresh.routingTitle') }}</h3>
                <p>{{ t('homeFresh.routingDescription') }}</p>
              </div>
            </article>

            <article v-for="benefit in benefits" :key="benefit.title" class="oh-benefit-item" data-home-reveal>
              <span class="oh-benefit-icon"><Icon :name="benefit.icon" size="md" /></span>
              <h3>{{ benefit.title }}</h3>
              <p>{{ benefit.description }}</p>
              <strong>{{ benefit.proof }}</strong>
            </article>
          </div>
        </div>
      </section>

      <section id="faq" class="oh-section oh-faq">
        <div class="oh-shell oh-faq-layout">
          <div class="oh-faq-copy" data-home-reveal>
            <h2>{{ t('homeFresh.faqTitle') }}</h2>
            <p>{{ t('homeFresh.faqDescription') }}</p>
            <a :href="docUrl || '/docs/usage-guide/'" class="oh-button oh-button-light">
              {{ t('homeFresh.readDocs') }}
              <Icon name="book" size="sm" />
            </a>
          </div>

          <div class="oh-faq-list" data-home-reveal>
            <article v-for="(item, index) in faqItems" :key="item.question" :class="{ 'is-open': openFaq === index }">
              <button
                type="button"
                :aria-expanded="openFaq === index"
                @click="openFaq = openFaq === index ? -1 : index"
              >
                <span>{{ item.question }}</span>
                <Icon :name="openFaq === index ? 'x' : 'plus'" size="sm" />
              </button>
              <div v-if="openFaq === index" class="oh-faq-answer">
                <p>{{ item.answer }}</p>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section class="oh-final-cta">
        <div class="oh-shell oh-final-inner" data-home-reveal>
          <div>
            <span>{{ t('homeFresh.finalKicker') }}</span>
            <h2>{{ t('homeFresh.finalTitle', { brand: siteName }) }}</h2>
            <p>{{ t('homeFresh.finalDescription') }}</p>
          </div>
          <router-link :to="isAuthenticated ? dashboardPath : '/register'" class="oh-button oh-button-dark">
            {{ isAuthenticated ? t('home.goToDashboard') : t('homeFresh.createAccount') }}
            <Icon name="arrowRight" size="sm" />
          </router-link>
        </div>
      </section>
    </main>

    <footer class="oh-footer">
      <div class="oh-shell oh-footer-main">
        <div class="oh-footer-brand">
          <router-link to="/home" class="oh-brand">
            <span class="oh-brand-mark">
              <img v-if="siteLogo" :src="siteLogo" :alt="siteName" />
              <OmnioMark v-else :label="siteName" />
            </span>
            <strong>{{ siteName }}</strong>
          </router-link>
          <p>{{ t('homeFresh.footerDescription') }}</p>
        </div>
        <div>
          <strong>{{ t('homeFresh.product') }}</strong>
          <router-link to="/model-plaza">{{ t('homeFresh.models') }}</router-link>
          <router-link :to="isAuthenticated ? '/usage' : '/login'">{{ t('homeFresh.usage') }}</router-link>
          <router-link :to="isAuthenticated ? '/omnio-pro' : '/login'">Omnio Pro</router-link>
        </div>
        <div>
          <strong>{{ t('homeFresh.resources') }}</strong>
          <a :href="docUrl || '/docs/usage-guide/'">{{ t('home.docs') }}</a>
          <router-link to="/key-usage">{{ t('homeFresh.keyUsage') }}</router-link>
          <router-link to="/legal/user-agreement">{{ t('homeFresh.userAgreement') }}</router-link>
        </div>
      </div>
      <div class="oh-shell oh-footer-bottom">
        <span>© {{ currentYear }} {{ siteName }}</span>
        <span>{{ t('homeFresh.footerNote') }}</span>
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
import '@/styles/home-v2.css'

type IconName = InstanceType<typeof Icon>['$props']['name']

const { t, locale } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const homeRoot = ref<HTMLElement | null>(null)
const isDark = ref(document.documentElement.classList.contains('dark'))
const mobileOpen = ref(false)
const selectedTask = ref('key')
const codeCopied = ref(false)
const openFaq = ref(0)
let revealObserver: IntersectionObserver | undefined
let copyResetTimer: ReturnType<typeof setTimeout> | undefined

const siteName = computed(() => {
  const configuredName = appStore.cachedPublicSettings?.site_name?.trim()
  if (configuredName) return configuredName
  const injectedName = appStore.siteName?.trim()
  return injectedName && injectedName !== 'Sub2API' ? injectedName : 'Omnio'
})

const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', {
  allowRelative: true,
  allowDataUrl: true
}))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const isHomeContentUrl = computed(() => /^https?:\/\//i.test(homeContent.value.trim()))
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => authStore.isAdmin ? '/admin/dashboard' : '/dashboard')
const gatewayOrigin = computed(() => typeof window === 'undefined' ? 'https://api.omnio.ai' : window.location.origin)
const currentYear = new Date().getFullYear()

const heroProofs = computed(() => [
  { icon: 'shield' as IconName, label: t('homeFresh.proofPlaintext'), value: t('homeFresh.proofPlaintextDetail') },
  { icon: 'sync' as IconName, label: t('homeFresh.proofModel'), value: t('homeFresh.proofModelDetail') },
  { icon: 'dollar' as IconName, label: t('homeFresh.proofCost'), value: t('homeFresh.proofCostDetail') }
])

const tasks = computed(() => [
  {
    id: 'key',
    icon: 'key' as IconName,
    label: t('homeFresh.taskKey'),
    short: t('homeFresh.taskKeyShort'),
    title: t('homeFresh.keyTitle'),
    promptLabel: t('homeFresh.setupLabel'),
    prompt: t('homeFresh.keyPrompt'),
    action: t('homeFresh.createKey'),
    route: isAuthenticated.value ? '/keys' : '/register',
    model: t('homeFresh.unifiedGateway'),
    modelIcon: 'openrouter',
    resultIcon: 'shield' as IconName,
    resultTitle: t('homeFresh.keyResult'),
    resultDescription: t('homeFresh.keyResultDescription'),
    facts: [
      { label: t('homeFresh.protocol'), value: 'OpenAI' },
      { label: t('homeFresh.scope'), value: t('homeFresh.allModels') },
      { label: t('homeFresh.billing'), value: t('homeFresh.usageBased') }
    ]
  },
  {
    id: 'image',
    icon: 'sparkles' as IconName,
    label: t('homeFresh.taskImage'),
    short: t('homeFresh.taskImageShort'),
    title: t('homeFresh.imageTitle'),
    promptLabel: t('homeFresh.promptLabel'),
    prompt: t('homeFresh.imagePrompt'),
    action: t('homeFresh.openImageStudio'),
    route: isAuthenticated.value ? '/image-generation' : '/register',
    model: 'Gemini',
    modelIcon: 'gemini',
    resultIcon: 'sparkles' as IconName,
    resultTitle: t('homeFresh.imageResult'),
    resultDescription: t('homeFresh.imageResultDescription'),
    facts: [
      { label: t('homeFresh.mode'), value: t('homeFresh.textToImage') },
      { label: t('homeFresh.output'), value: '1:1' },
      { label: t('homeFresh.history'), value: t('homeFresh.saved') }
    ]
  },
  {
    id: 'models',
    icon: 'cube' as IconName,
    label: t('homeFresh.taskModels'),
    short: t('homeFresh.taskModelsShort'),
    title: t('homeFresh.modelsTitle'),
    promptLabel: t('homeFresh.goalLabel'),
    prompt: t('homeFresh.modelsPrompt'),
    action: t('homeFresh.openModelPlaza'),
    route: '/model-plaza',
    model: 'Claude',
    modelIcon: 'claude',
    resultIcon: 'cube' as IconName,
    resultTitle: t('homeFresh.modelMatch'),
    resultDescription: t('homeFresh.modelMatchDescription'),
    facts: [
      { label: t('homeFresh.strength'), value: t('homeFresh.longContext') },
      { label: t('homeFresh.status'), value: t('homeFresh.available') },
      { label: t('homeFresh.switching'), value: t('homeFresh.instant') }
    ]
  },
  {
    id: 'usage',
    icon: 'chart' as IconName,
    label: t('homeFresh.taskUsage'),
    short: t('homeFresh.taskUsageShort'),
    title: t('homeFresh.usageTitle'),
    promptLabel: t('homeFresh.periodLabel'),
    prompt: t('homeFresh.usagePrompt'),
    action: t('homeFresh.viewUsage'),
    route: isAuthenticated.value ? '/usage' : '/register',
    model: t('homeFresh.allModels'),
    modelIcon: 'gpt-5',
    resultIcon: 'chart' as IconName,
    resultTitle: t('homeFresh.usageResult'),
    resultDescription: t('homeFresh.usageResultDescription'),
    facts: [
      { label: t('homeFresh.requests'), value: '1,284' },
      { label: t('homeFresh.tokens'), value: '2.8M' },
      { label: t('homeFresh.cost'), value: '$4.28' }
    ]
  }
])

const activeTask = computed(() => tasks.value.find((task) => task.id === selectedTask.value) || tasks.value[0])

const models = [
  { name: 'OpenAI', icon: 'gpt-5' },
  { name: 'Claude', icon: 'claude' },
  { name: 'Gemini', icon: 'gemini' },
  { name: 'Grok', icon: 'grok' },
  { name: 'Qwen', icon: 'qwen' },
  { name: 'DeepSeek', icon: 'deepseek' }
]

const modelDirectory = computed(() => [
  { name: 'GPT', icon: 'gpt-5', capability: t('homeFresh.reasoningAndTools') },
  { name: 'Claude', icon: 'claude', capability: t('homeFresh.longContextWork') },
  { name: 'Gemini', icon: 'gemini', capability: t('homeFresh.multimodalCreation') },
  { name: 'Grok', icon: 'grok', capability: t('homeFresh.realtimeResearch') },
  { name: 'Qwen', icon: 'qwen', capability: t('homeFresh.multilingualTasks') }
])

const requestRows = [
  { model: 'GPT-5', icon: 'gpt-5', tokens: '18.4K', latency: '1.2s', cost: '$0.82' },
  { model: 'Claude', icon: 'claude', tokens: '9.8K', latency: '0.9s', cost: '$0.44' },
  { model: 'Gemini', icon: 'gemini', tokens: '24.1K', latency: '1.6s', cost: '$0.36' }
]

const benefits = computed(() => [
  {
    icon: 'eye' as IconName,
    title: t('homeFresh.visibilityTitle'),
    description: t('homeFresh.visibilityDescription'),
    proof: t('homeFresh.requestLevelData')
  },
  {
    icon: 'dollar' as IconName,
    title: t('homeFresh.costTitle'),
    description: t('homeFresh.costDescription'),
    proof: t('homeFresh.transparentRates')
  },
  {
    icon: 'lock' as IconName,
    title: t('homeFresh.controlTitle'),
    description: t('homeFresh.controlDescription'),
    proof: t('homeFresh.quotaAndAccess')
  }
])

const faqItems = computed(() => [
  { question: t('homeFresh.faqWhat'), answer: t('homeFresh.faqWhatAnswer', { brand: siteName.value }) },
  { question: t('homeFresh.faqSdk'), answer: t('homeFresh.faqSdkAnswer') },
  { question: t('homeFresh.faqChoose'), answer: t('homeFresh.faqChooseAnswer') },
  { question: t('homeFresh.faqBilling'), answer: t('homeFresh.faqBillingAnswer') },
  { question: t('homeFresh.faqData'), answer: t('homeFresh.faqDataAnswer') }
])

const codeSample = computed(() => `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "${gatewayOrigin.value}/v1",
  apiKey: process.env.OMNIO_API_KEY
});

const response = await client.chat.completions.create({
  model: "gpt-5",
  messages: [{ role: "user", content: "Hello" }]
});`)

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

async function copyCode() {
  try {
    await navigator.clipboard.writeText(codeSample.value)
    codeCopied.value = true
    if (copyResetTimer) clearTimeout(copyResetTimer)
    copyResetTimer = setTimeout(() => { codeCopied.value = false }, 1800)
  } catch {
    codeCopied.value = false
  }
}

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()

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
    }, { threshold: 0.14, rootMargin: '0px 0px -7% 0px' })

    targets.forEach((target) => revealObserver?.observe(target))
  })
})

onBeforeUnmount(() => {
  revealObserver?.disconnect()
  if (copyResetTimer) clearTimeout(copyResetTimer)
})
</script>
