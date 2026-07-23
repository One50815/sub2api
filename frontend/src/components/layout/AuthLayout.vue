<template>
  <div class="sub2-auth-shell dark" @pointermove="handlePointerMove">
    <section class="sub2-auth-stage">
      <div class="auth-pointer-glow"></div>
      <div class="auth-grid"></div>
      <div class="auth-particles"></div>
      <div class="auth-scan"></div>
      <div class="auth-aurora auth-aurora-one"></div>
      <div class="auth-aurora auth-aurora-two"></div>

      <header class="sub2-auth-header">
        <router-link to="/" class="sub2-auth-brand">
          <span class="sub2-auth-logo">
            <img :src="siteLogo || '/logo.png'" alt="" />
          </span>
          <span>{{ siteName }}</span>
        </router-link>
        <div class="sub2-auth-toolbar">
          <LocaleSwitcher />
          <button type="button" class="auth-tool-button" :aria-label="isDark ? 'Light theme' : 'Dark theme'" @click="toggleTheme">
            <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
          </button>
        </div>
      </header>

      <div class="sub2-auth-content">
        <section class="auth-gateway-visual" aria-hidden="true">
          <div class="auth-gateway-copy">
            <span class="auth-gateway-eyebrow">{{ t('auth.gatewayEyebrow') }}</span>
            <h1>
              {{ t('auth.gatewayTitle') }}
              <span>{{ t('auth.gatewayTitleAccent') }}</span>
            </h1>
            <p>{{ t('auth.gatewayDescription') }}</p>
            <div class="auth-value-list">
              <div v-for="item in valueItems" :key="item.title" class="auth-value-item">
                <span class="auth-value-icon"><Icon :name="item.icon" size="sm" /></span>
                <span>
                  <strong>{{ item.title }}</strong>
                  <small>{{ item.detail }}</small>
                </span>
              </div>
            </div>
          </div>

          <div class="auth-network">
            <img
              class="auth-core-asset"
              src="/assets/auth/omnio-core-orbit.webp?v=2"
              alt=""
              width="1000"
              height="1000"
              fetchpriority="high"
            />
            <svg class="auth-network-routes" viewBox="0 0 600 600" fill="none">
              <path
                v-for="(path, index) in routePaths"
                :key="path"
                :d="path"
                class="auth-route-line"
                :style="{ animationDelay: `${index * -0.38}s` }"
              />
            </svg>
            <div class="auth-network-core">
              <img src="/assets/brand/omnio-mark.svg?v=3" alt="" />
              <strong>{{ siteName }}</strong>
              <span>{{ t('auth.gatewayEyebrow') }}</span>
            </div>
            <div
              v-for="node in modelNodes"
              :key="node.label"
              class="auth-model-node"
              :style="{ left: node.left, top: node.top, animationDelay: node.delay }"
            >
              <span class="auth-model-mark">
                <ModelIcon :model="node.iconModel" size="22px" />
              </span>
              <span><strong>{{ node.label }}</strong><small>{{ node.provider }}</small></span>
            </div>
          </div>
        </section>

        <div class="sub2-auth-card-column">
          <div class="sub2-auth-card">
            <slot />
            <div class="sub2-auth-footer"><slot name="footer" /></div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import Icon from '@/components/icons/Icon.vue'

const appStore = useAppStore()
const { t } = useI18n()
const isDark = ref(document.documentElement.classList.contains('dark'))

const siteName = computed(() => appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))

const valueItems = computed(() => [
  { icon: 'terminal' as const, title: t('auth.gatewayUnified'), detail: t('auth.gatewayUnifiedDetail') },
  { icon: 'chart' as const, title: t('auth.gatewayRouting'), detail: t('auth.gatewayRoutingDetail') },
  { icon: 'shield' as const, title: t('auth.gatewayReliable'), detail: t('auth.gatewayReliableDetail') }
])

const modelNodes = [
  { iconModel: 'gpt-5.6-sol', label: 'GPT-5.6 Sol', provider: 'OpenAI', left: '50%', top: '6%', delay: '0s' },
  { iconModel: 'claude-fable-5', label: 'Claude Fable 5', provider: 'Anthropic', left: '84%', top: '25%', delay: '-.7s' },
  { iconModel: 'gemini-3.5-flash', label: 'Gemini 3.5 Flash', provider: 'Google', left: '91%', top: '55%', delay: '-1.4s' },
  { iconModel: 'grok-4.5', label: 'Grok 4.5', provider: 'xAI', left: '78%', top: '84%', delay: '-2.1s' },
  { iconModel: 'deepseek-v4-pro', label: 'DeepSeek V4 Pro', provider: 'DeepSeek', left: '48%', top: '94%', delay: '-2.8s' },
  { iconModel: 'qwen3.7-max', label: 'Qwen3.7 Max', provider: 'Alibaba', left: '9%', top: '57%', delay: '-3.5s' },
  { iconModel: 'openrouter', label: '100+', provider: 'Models', left: '17%', top: '24%', delay: '-4.2s' }
]

const routePaths = [
  'M300 300 C300 210 300 110 300 42',
  'M300 300 C390 260 470 185 505 150',
  'M300 300 C410 300 510 320 548 330',
  'M300 300 C390 365 455 455 468 504',
  'M300 300 C300 395 288 505 288 564',
  'M300 300 C205 320 105 340 54 342',
  'M300 300 C215 245 135 180 102 144'
]

function handlePointerMove(event: PointerEvent) {
  if (event.pointerType === 'touch') return
  const target = event.currentTarget as HTMLElement
  target.style.setProperty('--auth-pointer-x', `${event.clientX}px`)
  target.style.setProperty('--auth-pointer-y', `${event.clientY}px`)
}

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>
