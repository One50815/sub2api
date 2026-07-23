<template>
  <div class="space-y-4">
    <button type="button" :disabled="disabled" class="btn btn-secondary w-full" @click="startLogin">
      <svg
        class="icon mr-2"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        width="20"
        height="20"
        aria-hidden="true"
        style="flex-shrink: 0"
      >
        <circle cx="12" cy="12" r="12" fill="#1677FF" />
        <path
          fill="white"
          d="M17.9 7.15c-.16-.08-.76-.36-1.52-.63-1.75-.65-4.35-1.48-7.54-2.22-.28-.06-.42.31-.19.48 1.13.82 2.08 1.55 2.86 2.18L8.54 6.4a.31.31 0 0 0-.28.52c.86.77 1.63 1.43 2.31 1.98l-2.1-.35a.31.31 0 0 0-.27.53c1.11.96 2.08 1.7 2.91 2.26l-1.41.02a.32.32 0 0 0-.19.57c.84.61 1.58 1.08 2.23 1.43-.45.78-1.02 1.72-1.73 2.82-.12.19.07.42.28.34 2.49-.94 4.41-2.43 5.67-4.45.15-.24.06-.55-.19-.68l-1.13-.57c1.41-.82 2.55-1.76 3.43-2.82.24-.29.15-.69-.17-.85Z"
        />
      </svg>
      {{ t('auth.dingtalk.signIn') }}
    </button>

    <div v-if="showDivider" class="flex items-center gap-3">
      <div class="h-px flex-1 bg-gray-200 dark:bg-dark-700"></div>
      <span class="text-xs text-gray-500 dark:text-dark-400">
        {{ t('auth.oauthOrContinue') }}
      </span>
      <div class="h-px flex-1 bg-gray-200 dark:bg-dark-700"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { resolveAffiliateReferralCode, storeOAuthAffiliateCode } from '@/utils/oauthAffiliate'

const props = withDefaults(defineProps<{
  disabled?: boolean
  affCode?: string
  showDivider?: boolean
}>(), {
  showDivider: true
})

const route = useRoute()
const { t } = useI18n()

function startLogin(): void {
  const redirectTo = (route.query.redirect as string) || '/dashboard'
  storeOAuthAffiliateCode(resolveAffiliateReferralCode(props.affCode, route.query.aff, route.query.aff_code))
  const apiBase = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api/v1'
  const normalized = apiBase.replace(/\/$/, '')
  const startURL = `${normalized}/auth/oauth/dingtalk/start?redirect=${encodeURIComponent(redirectTo)}`
  window.location.href = startURL
}
</script>
