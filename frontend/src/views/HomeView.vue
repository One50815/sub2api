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

  <div v-else class="omnio-public-shell dark">
    <header :class="['omnio-public-header', { 'is-scrolled': headerScrolled }]">
      <nav class="omnio-public-nav" :aria-label="t('homeReplica.primaryNavigation')">
        <router-link to="/home" class="omnio-public-brand">
          <img :src="siteLogo || '/assets/brand/omnio-mark.svg?v=3'" alt="" />
          <span>{{ siteName }}</span>
        </router-link>

        <div class="omnio-public-links">
          <router-link to="/home">{{ t('nav.home') }}</router-link>
          <router-link :to="dashboardPath">{{ t('nav.console') }}</router-link>
          <a href="/docs/pricing/">{{ t('homeReplica.pricing') }}</a>
          <a href="/docs/usage-guide/">{{ t('homeReplica.tutorial') }}</a>
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a>
          <router-link v-else to="/key-usage">{{ t('home.docs') }}</router-link>
          <a href="/docs/about-omnio/">{{ t('homeReplica.about') }}</a>
        </div>

        <div class="omnio-public-actions">
          <LocaleSwitcher />
          <button type="button" class="omnio-header-icon" :title="isDark ? t('home.switchToLight') : t('home.switchToDark')" @click="toggleTheme">
            <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
          </button>
          <button v-if="isAuthenticated" type="button" class="omnio-header-icon" :title="t('homeReplica.notifications')">
            <Icon name="bell" size="sm" />
          </button>
          <router-link v-if="isAuthenticated" :to="dashboardPath" class="omnio-profile-avatar" :aria-label="t('home.goToDashboard')">
            {{ profileInitial }}
          </router-link>
          <router-link v-else to="/login" class="omnio-header-cta">
            {{ t('home.login') }}
          </router-link>
          <button type="button" class="omnio-mobile-menu-button" :aria-expanded="mobileOpen" :aria-label="t('homeReplica.toggleNavigation')" @click="mobileOpen = !mobileOpen"><Icon :name="mobileOpen ? 'x' : 'menu'" size="sm" /></button>
        </div>
      </nav>
      <div v-if="mobileOpen" class="omnio-mobile-menu">
        <router-link to="/home" @click="mobileOpen = false">{{ t('nav.home') }}</router-link>
        <router-link :to="dashboardPath" @click="mobileOpen = false">{{ t('nav.console') }}</router-link>
        <a href="/docs/pricing/" @click="mobileOpen = false">{{ t('homeReplica.pricing') }}</a>
        <a href="/docs/usage-guide/" @click="mobileOpen = false">{{ t('homeReplica.tutorial') }}</a>
        <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a>
        <router-link v-else to="/key-usage" @click="mobileOpen = false">{{ t('home.docs') }}</router-link>
        <a href="/docs/about-omnio/" @click="mobileOpen = false">{{ t('homeReplica.about') }}</a>
      </div>
    </header>

    <main class="suno-replica-shell">
      <section class="suno-replica-hero">
        <div class="suno-replica-hero-aura"></div>
        <div class="suno-replica-grain"></div>

        <article class="suno-hero-side-card suno-hero-side-card-left">
          <div class="suno-model-art-card suno-model-art-cyan">
            <img :src="models[0].artwork" alt="" class="suno-model-art-image" />
            <div class="suno-model-art-overlay"></div>
            <div class="suno-model-art-content">
              <div class="suno-model-art-meta"><span class="suno-model-provider"><ModelIcon :model="models[0].name" size="18px" />{{ models[0].provider }}</span><span>● LIVE</span></div>
              <div><small>{{ models[0].capability }}</small><h3>{{ models[0].name }}</h3><p><span>{{ models[0].route }}</span><span>{{ models[0].latency }}</span></p></div>
            </div>
          </div>
        </article>
        <article class="suno-hero-side-card suno-hero-side-card-right">
          <div class="suno-model-art-card suno-model-art-violet">
            <img :src="models[2].artwork" alt="" class="suno-model-art-image" />
            <div class="suno-model-art-overlay"></div>
            <div class="suno-model-art-content">
              <div class="suno-model-art-meta"><span class="suno-model-provider"><ModelIcon :model="models[2].name" size="18px" />{{ models[2].provider }}</span><span>● LIVE</span></div>
              <div><small>{{ models[2].capability }}</small><h3>{{ models[2].name }}</h3><p><span>{{ models[2].route }}</span><span>{{ models[2].latency }}</span></p></div>
            </div>
          </div>
        </article>

        <div class="suno-hero-content">
          <div class="suno-headline-frame">
            <h1 class="suno-replica-hero-title" :key="headlineIndex">
              {{ headlines[headlineIndex] }}<span class="suno-type-caret" aria-hidden="true"></span>
            </h1>
          </div>
          <p class="suno-hero-copy">{{ t('homeReplica.heroCopy', { brand: siteName }) }}</p>

          <div class="suno-gateway-quickstart">
            <div class="suno-quickstart-header">
              <div><span class="suno-quickstart-eyebrow">{{ t('homeReplica.quickstart') }}</span><strong>{{ siteName }} API</strong></div>
              <span class="suno-quickstart-status">● {{ t('homeReplica.gatewayOnline') }}</span>
            </div>
            <div class="suno-quickstart-endpoint">
              <code>{{ gatewayOrigin }}/v1</code>
              <button type="button" @click="copyText(`${gatewayOrigin}/v1`)"><Icon name="copy" size="xs" /><span>{{ t('common.copy') }}</span></button>
            </div>
            <div class="suno-quickstart-tabs" role="tablist" :aria-label="t('homeReplica.sdk')">
              <button v-for="snippet in snippets" :key="snippet.id" type="button" :data-active="activeSnippet === snippet.id" @click="activeSnippet = snippet.id">{{ snippet.label }}</button>
            </div>
            <div class="suno-quickstart-code">
              <pre>{{ selectedSnippet.code }}</pre>
              <button type="button" @click="copyText(selectedSnippet.code)"><Icon name="copy" size="xs" /><span>{{ t('common.copy') }}</span></button>
            </div>
            <div class="suno-quickstart-actions">
              <router-link to="/key-usage"><Icon name="book" size="xs" />{{ t('homeReplica.readDocs') }}</router-link>
              <router-link :to="isAuthenticated ? dashboardPath : '/register'" data-primary="true">{{ isAuthenticated ? t('home.goToDashboard') : t('homeReplica.getApiKey') }}<Icon name="arrowRight" size="xs" /></router-link>
            </div>
          </div>
        </div>

        <div class="suno-provider-rail" :aria-label="t('home.providers.title')">
          <div class="suno-provider-rail-track">
            <div v-for="(model, index) in railModels" :key="`${model.name}-${index}`" class="suno-provider-model">
              <span class="suno-provider-model-icon"><ModelIcon :model="model.name" size="28px" /></span>
              <span>{{ model.name }}</span>
            </div>
          </div>
        </div>
      </section>

      <section class="suno-showcase-section">
        <div class="suno-section-heading">
          <h2 class="suno-replica-section-title">{{ t('homeReplica.everyModel') }}</h2>
          <p class="suno-replica-section-copy">{{ t('homeReplica.everyModelCopy') }}</p>
        </div>
        <div class="suno-showcase-viewport">
          <div class="suno-model-showcase-track">
            <article v-for="model in showcaseModels" :key="model.key" :class="['suno-model-art-card', `suno-model-art-${model.accent}`, 'suno-showcase-model-card']">
              <img :src="model.artwork" alt="" class="suno-model-art-image" loading="lazy" />
              <div class="suno-model-art-overlay"></div>
              <div class="suno-model-art-content">
                <div class="suno-model-art-meta"><span class="suno-model-provider"><ModelIcon :model="model.name" size="18px" />{{ model.provider }}</span><span>● LIVE</span></div>
                <div><small>{{ model.capability }}</small><h3>{{ model.name }}</h3><p><span>{{ model.route }}</span><span>{{ model.latency }}</span></p></div>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section class="suno-feature-section">
        <div class="suno-section-heading suno-section-heading-left">
          <h2 class="suno-replica-section-title">{{ t('homeReplica.everything') }}</h2>
        </div>
        <div class="suno-feature-grid">
          <article v-for="(feature, index) in features" :key="feature.title" class="suno-feature-card">
            <div><Icon :name="feature.icon" size="md" class="text-white/70" /><h3>{{ feature.title }}</h3><p>{{ feature.description }}</p></div>
            <div v-if="index === 0" class="suno-feature-visual suno-feature-network">
              <div v-for="model in models.slice(0, 4)" :key="model.name" class="suno-feature-model-row"><span class="suno-feature-model-icon"><ModelIcon :model="model.name" size="18px" /></span><span>{{ model.name }}</span><span class="ml-auto text-white/36">online</span></div>
            </div>
            <div v-else-if="index === 1" class="suno-feature-visual suno-feature-bars"><span v-for="height in featureBars" :key="height" class="suno-feature-bar" :style="{ height: `${height}%` }"></span></div>
            <div v-else-if="index === 2" class="suno-feature-visual suno-feature-network"><div v-for="row in policyRows" :key="row[0]" class="suno-feature-policy-row"><span>{{ row[0] }}</span><span class="text-emerald-300/75">{{ row[1] }}</span></div></div>
            <div v-else-if="index === 3" class="suno-feature-visual suno-feature-failover"><p><span>primary_01</span><span>degraded</span></p><hr /><p><span>fallback_02</span><span>active · 38ms</span></p></div>
            <div v-else-if="index === 4" class="suno-feature-visual suno-feature-cost"><div><small>{{ t('homeReplica.averageRequest') }}</small><strong>$0.0018</strong></div><span>-19.1%</span></div>
            <div v-else class="suno-feature-visual suno-feature-routes"><span v-for="route in protocolRoutes" :key="route" class="suno-feature-route">{{ route }}</span></div>
          </article>
        </div>
      </section>

      <section class="suno-pricing-section">
        <div class="suno-section-heading">
          <h2 class="suno-replica-section-title">{{ t('homeReplica.flexiblePricing') }}</h2>
          <p class="suno-replica-section-copy">{{ t('homeReplica.flexiblePricingCopy') }}</p>
          <div class="suno-pricing-toggle"><button type="button" :data-active="billingMode === 'usage'" @click="billingMode = 'usage'">{{ t('homeReplica.usageBilling') }}</button><button type="button" :data-active="billingMode === 'pro'" @click="billingMode = 'pro'">Omnio Pro</button></div>
        </div>
        <div class="suno-pricing-grid">
          <article v-for="(plan, index) in plans" :key="plan.title" class="suno-pricing-card" :data-featured="(billingMode === 'usage' && index === 0) || (billingMode === 'pro' && index === 1)">
            <div class="suno-plan-heading"><div><h3>{{ plan.title }}</h3><p>{{ plan.description }}</p></div><span v-if="plan.badge" class="suno-plan-badge">{{ plan.badge }}</span></div>
            <div class="suno-plan-price">{{ plan.price }}</div>
            <router-link :to="plan.to" class="suno-plan-button">{{ plan.action }}</router-link>
            <ul><li v-for="item in plan.features" :key="item"><span class="suno-plan-check">✓</span>{{ item }}</li></ul>
          </article>
        </div>
      </section>

      <section class="suno-trust-section">
        <div class="suno-trust-aura"></div><div class="suno-replica-grain"></div>
        <div class="suno-trust-content">
          <h2 class="suno-replica-section-title">{{ t('homeReplica.builtLike') }} <span>{{ t('homeReplica.operatedLike') }}</span></h2>
          <p class="suno-replica-section-copy">{{ t('homeReplica.intelligenceLayer') }}</p>
          <div class="suno-trust-cards"><div><span>OpenAI-compatible</span><strong>1 API</strong><small>{{ t('homeReplica.protocolSurface') }}</small></div><div><span>{{ t('homeReplica.providers') }}</span><strong>40+</strong><small>{{ t('homeReplica.modelNetwork') }}</small></div></div>
        </div>
      </section>

      <section class="suno-explore-section">
        <div class="suno-section-heading"><h2 class="suno-replica-section-title">{{ t('homeReplica.firstRequest') }} <span>{{ t('homeReplica.fullControl') }}</span></h2><p class="suno-replica-section-copy">{{ t('homeReplica.exploreCopy') }}</p></div>
        <div class="suno-trace-rail">
          <article v-for="model in models" :key="model.name" class="suno-trace-card"><div class="suno-trace-meta"><span>REQUEST TRACE</span><span>200</span></div><div class="suno-trace-pulse"><ModelIcon :model="model.name" size="68px" /></div><h3>{{ model.name }}</h3><p><span>{{ model.route }}</span><span>{{ model.latency }}</span></p></article>
        </div>
      </section>

      <section class="suno-faq-section">
        <div class="suno-section-heading"><h2 class="suno-replica-section-title">{{ t('homeReplica.faq') }}</h2><p class="suno-replica-section-copy">{{ t('homeReplica.faqCopy', { brand: siteName }) }}</p></div>
        <div class="suno-faq-list"><div v-for="(faq, index) in faqs" :key="faq[0]" class="suno-faq-item" :data-open="openFaq === index"><button type="button" :aria-expanded="openFaq === index" @click="openFaq = openFaq === index ? null : index"><span>{{ faq[0] }}</span><span>＋</span></button><div v-show="openFaq === index"><p>{{ faq[1] }}</p></div></div></div>
      </section>

      <footer class="suno-home-footer">
        <router-link :to="isAuthenticated ? dashboardPath : '/register'" class="suno-footer-cta">{{ isAuthenticated ? t('home.goToDashboard') : t('homeReplica.createGateway') }}<Icon name="arrowRight" size="xs" /></router-link>
        <div class="suno-footer-grid"><div><div class="suno-brand-lockup"><img src="/assets/brand/omnio-mark.svg?v=3" alt="" /><strong>{{ siteName }}</strong></div><p>{{ t('homeReplica.intelligenceLayer') }}</p></div><div><h3>{{ t('homeReplica.product') }}</h3><router-link to="/key-usage">{{ t('homeReplica.modelPricing') }}</router-link><router-link to="/available-channels">{{ t('nav.availableChannels') }}</router-link><router-link :to="dashboardPath">{{ t('home.dashboard') }}</router-link></div><div><h3>{{ t('homeReplica.resources') }}</h3><a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ t('home.docs') }}</a><router-link to="/legal/user-agreement">{{ t('homeReplica.userAgreement') }}</router-link><router-link to="/login">{{ t('home.login') }}</router-link></div></div>
        <div class="suno-footer-bottom"><span>© {{ currentYear }} {{ siteName }}</span><span class="suno-footer-attribution">{{ t('homeReplica.builtOnSub2api') }}</span></div>
      </footer>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import { sanitizeUrl } from '@/utils/url'

type SnippetId = 'curl' | 'python' | 'node'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const isDark = ref(document.documentElement.classList.contains('dark'))
const headlineIndex = ref(0)
const activeSnippet = ref<SnippetId>('curl')
const billingMode = ref<'usage' | 'pro'>('usage')
const openFaq = ref<number | null>(null)
const headerScrolled = ref(false)
const mobileOpen = ref(false)
let headlineTimer: number | undefined

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Omnio')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isHomeContentUrl = computed(() => /^https?:\/\//i.test(homeContent.value.trim()))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => authStore.isAdmin ? '/admin/dashboard' : '/dashboard')
const profileInitial = computed(() => (authStore.user?.username || authStore.user?.email || 'A').trim().charAt(0).toUpperCase())
const currentYear = new Date().getFullYear()
const gatewayOrigin = computed(() => typeof window === 'undefined' ? '' : window.location.origin)

const headlines = computed(() => [t('homeReplica.headlineGateway'), t('homeReplica.headlineVisible'), t('homeReplica.headlineControl')])
const models = computed(() => [
  { name: 'GPT-5.6 Sol', provider: 'OpenAI', latency: '36ms', route: '/v1/responses', accent: 'cyan', artwork: '/assets/home/models/gpt-5-6-sol.jpg', capability: t('homeReplica.capabilityReasoning') },
  { name: 'Claude Sonnet 5', provider: 'Anthropic', latency: '44ms', route: '/v1/messages', accent: 'amber', artwork: '/assets/home/models/claude-sonnet-5.jpg', capability: t('homeReplica.capabilityContext') },
  { name: 'Gemini 3.5 Flash', provider: 'Google', latency: '29ms', route: '/v1beta/models', accent: 'violet', artwork: '/assets/home/models/gemini-3-5-flash.jpg', capability: t('homeReplica.capabilityMultimodal') },
  { name: 'DeepSeek V4 Pro', provider: 'DeepSeek', latency: '32ms', route: '/v1/chat/completions', accent: 'teal', artwork: '/assets/home/models/deepseek-v4-pro.jpg', capability: t('homeReplica.capabilityEfficient') },
  { name: 'Qwen3.7 Max', provider: 'Alibaba', latency: '34ms', route: '/v1/chat/completions', accent: 'jade', artwork: '/assets/home/models/qwen3-7-max.jpg', capability: t('homeReplica.capabilityMultilingual') }
])
const railModels = computed(() => [...models.value, ...models.value])
const showcaseModels = computed(() => [...models.value, ...models.value.slice(0, 2)].map((model, index) => ({ ...model, key: `${index}-${model.name}` })))
const snippets = computed(() => [
  { id: 'curl' as const, label: 'cURL', code: `curl ${gatewayOrigin.value}/v1/chat/completions \\\n  -H "Authorization: Bearer $SUB2API_KEY" \\\n  -H "Content-Type: application/json" \\\n  -d '{"model":"gpt-5.6-sol","messages":[{"role":"user","content":"Hello"}]}'` },
  { id: 'python' as const, label: 'Python', code: `from openai import OpenAI\nclient = OpenAI(base_url="${gatewayOrigin.value}/v1", api_key=SUB2API_KEY)\nclient.chat.completions.create(model="gpt-5.6-sol", messages=[{"role": "user", "content": "Hello"}])` },
  { id: 'node' as const, label: 'Node', code: `const client = new OpenAI({ baseURL: '${gatewayOrigin.value}/v1', apiKey: process.env.SUB2API_KEY })\nawait client.chat.completions.create({ model: 'gpt-5.6-sol', messages: [{ role: 'user', content: 'Hello' }] })` }
])
const selectedSnippet = computed(() => snippets.value.find((snippet) => snippet.id === activeSnippet.value) || snippets.value[0])
const features = computed(() => [
  { icon: 'server' as const, title: t('homeReplica.modelNetwork'), description: t('homeReplica.modelNetworkCopy') },
  { icon: 'chart' as const, title: t('homeReplica.requestVisible'), description: t('homeReplica.requestVisibleCopy') },
  { icon: 'shield' as const, title: t('homeReplica.policySpeed'), description: t('homeReplica.policySpeedCopy') },
  { icon: 'swap' as const, title: t('homeReplica.failover'), description: t('homeReplica.failoverCopy') },
  { icon: 'dollar' as const, title: t('homeReplica.costSignal'), description: t('homeReplica.costSignalCopy') },
  { icon: 'terminal' as const, title: t('homeReplica.protocolSurface'), description: t('homeReplica.protocolSurfaceCopy') }
])
const plans = computed(() => [
  { title: t('homeReplica.usageBased'), description: t('homeReplica.usageBasedCopy'), price: t('homeReplica.payAsYouGo'), action: t('homeReplica.viewPricing'), to: '/key-usage', features: [t('homeReplica.perModelPricing'), t('homeReplica.noCommitment'), t('homeReplica.usageAnalytics'), t('homeReplica.automaticRouting')] },
  { title: 'Omnio Pro', description: '独立倍率、专属分组和优先权益', price: t('homeReplica.monthlyPlans'), badge: t('homeReplica.mostPopular'), action: '查看 Omnio Pro', to: '/omnio-pro', features: [t('homeReplica.includedQuota'), 'Pro 专属倍率', 'Pro 专属分组', '优先支持'] },
  { title: t('homeReplica.enterprise'), description: t('homeReplica.enterpriseCopy'), price: t('homeReplica.custom'), badge: t('homeReplica.tailored'), action: t('homeReplica.contactUs'), to: '/register', features: [t('homeReplica.customQuota'), t('homeReplica.rolesPermissions'), t('homeReplica.onboardingSupport'), t('homeReplica.customRouting')] }
])
const faqs = computed(() => [
  [t('homeReplica.whatIs', { brand: siteName.value }), t('homeReplica.whatIsAnswer', { brand: siteName.value })],
  [t('homeReplica.existingSdk'), t('homeReplica.existingSdkAnswer', { brand: siteName.value })],
  [t('homeReplica.routingQuestion'), t('homeReplica.routingAnswer')],
  [t('homeReplica.trackQuestion'), t('homeReplica.trackAnswer')]
])
const featureBars = [28, 46, 35, 66, 48, 80, 58, 94, 72, 86, 64, 98]
const policyRows = [['model.access', 'ALLOW'], ['budget.limit', 'ACTIVE'], ['fallback.route', 'READY']]
const protocolRoutes = ['/v1/chat', '/v1/responses', '/v1/messages', '/v1beta']

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function updateHeaderState() {
  headerScrolled.value = window.scrollY > 20
}

async function copyText(value: string) {
  try { await navigator.clipboard.writeText(value) } catch { /* clipboard may be unavailable */ }
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  isDark.value = true
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
  updateHeaderState()
  window.addEventListener('scroll', updateHeaderState, { passive: true })
  headlineTimer = window.setInterval(() => { headlineIndex.value = (headlineIndex.value + 1) % headlines.value.length }, 3600)
})

onBeforeUnmount(() => {
  if (headlineTimer !== undefined) window.clearInterval(headlineTimer)
  window.removeEventListener('scroll', updateHeaderState)
})
</script>
