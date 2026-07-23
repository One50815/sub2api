<template>
  <AuthLayout>
    <div v-if="isProcessing" class="oauth-state oauth-state-centered">
      <div class="oauth-state-icon">
        <Icon name="shield" size="xl" :stroke-width="1.6" />
      </div>
      <div class="oauth-heading">
        <h2>{{ t('auth.oauth.processingTitle') }}</h2>
        <p>{{ t('auth.oauth.processingDescription') }}</p>
      </div>
      <div class="oauth-progress">
        <Icon name="refresh" size="sm" class="animate-spin" :stroke-width="2" />
        <span>{{ t('auth.oauth.processing') }}</span>
      </div>
      <p class="oauth-note">{{ t('auth.oauth.processingNote') }}</p>
    </div>

    <div v-else-if="needsRegistrationCompletion" class="oauth-state">
      <div class="oauth-state-header">
        <div class="oauth-state-icon oauth-state-icon-small">
          <Icon name="userPlus" size="lg" :stroke-width="1.7" />
        </div>
        <div class="oauth-heading oauth-heading-left">
          <h2>{{ t('auth.oidc.callbackTitle', { providerName }) }}</h2>
          <p>{{ registrationHint }}</p>
        </div>
      </div>

      <div class="oauth-form">
        <div class="oauth-field">
          <label class="input-label" for="oauth-registration-email">{{ t('auth.emailLabel') }}</label>
          <input id="oauth-registration-email" class="input w-full" type="email" :value="registrationEmail" readonly disabled />
        </div>
        <div class="oauth-field">
          <label class="input-label" for="oauth-registration-password">{{ t('auth.passwordLabel') }}</label>
          <input id="oauth-registration-password" v-model="password" type="password" class="input w-full" :placeholder="t('auth.createPasswordPlaceholder')" :disabled="isSubmitting" autocomplete="new-password" @keyup.enter="handleSubmitRegistration" />
        </div>
        <div class="oauth-field">
          <label class="input-label" for="oauth-registration-confirm">{{ t('auth.confirmPassword') }}</label>
          <input id="oauth-registration-confirm" v-model="confirmPassword" type="password" class="input w-full" :placeholder="t('auth.confirmPasswordPlaceholder')" :disabled="isSubmitting" autocomplete="new-password" @keyup.enter="handleSubmitRegistration" />
        </div>
        <div v-if="invitationRequired" class="oauth-field">
          <label class="input-label" for="oauth-invitation-code">{{ t('auth.invitationCodeLabel') }}</label>
          <input id="oauth-invitation-code" v-model="invitationCode" type="text" class="input w-full" :placeholder="t('auth.invitationCodePlaceholder')" :disabled="isSubmitting" @keyup.enter="handleSubmitRegistration" />
        </div>
        <p v-if="registrationError" class="oauth-error" role="alert">
          <Icon name="exclamationCircle" size="sm" />
          <span>{{ registrationError }}</span>
        </p>
        <button class="btn btn-primary oauth-primary-action" type="button" :disabled="isSubmitting || !canSubmitRegistration" @click="handleSubmitRegistration">
          <Icon v-if="isSubmitting" name="refresh" size="sm" class="animate-spin" />
          {{ isSubmitting ? t('common.processing') : t('auth.oidc.completeRegistration') }}
        </button>
      </div>
    </div>

    <div v-else-if="invalidCallback" class="oauth-state oauth-state-centered">
      <div class="oauth-state-icon oauth-state-icon-error">
        <Icon name="exclamationCircle" size="xl" :stroke-width="1.6" />
      </div>
      <div class="oauth-heading">
        <h2>{{ t('auth.oauth.invalidCallbackTitle') }}</h2>
        <p>{{ t('auth.oauth.invalidCallbackHint') }}</p>
      </div>
      <button class="btn btn-primary oauth-inline-action" type="button" @click="router.replace('/login')">
        {{ t('auth.backToLogin') }}
      </button>
    </div>

    <div v-else class="oauth-state">
      <div class="oauth-state-header">
        <div class="oauth-state-icon oauth-state-icon-small">
          <Icon name="key" size="lg" :stroke-width="1.7" />
        </div>
        <div class="oauth-heading oauth-heading-left">
          <h2>{{ t('auth.oauth.callbackTitle') }}</h2>
          <p>{{ t('auth.oauth.callbackHint') }}</p>
        </div>
      </div>

      <div class="oauth-code-list">
        <div class="oauth-field">
          <label class="input-label" for="oauth-code">{{ t('auth.oauth.code') }}</label>
          <div class="oauth-copy-row">
            <input id="oauth-code" class="input oauth-code-input" :value="code" readonly />
            <button class="btn btn-secondary" type="button" :disabled="!code" @click="copy(code)">{{ t('common.copy') }}</button>
          </div>
        </div>
        <div class="oauth-field">
          <label class="input-label" for="oauth-state">{{ t('auth.oauth.state') }}</label>
          <div class="oauth-copy-row">
            <input id="oauth-state" class="input oauth-code-input" :value="state" readonly />
            <button class="btn btn-secondary" type="button" :disabled="!state" @click="copy(state)">{{ t('common.copy') }}</button>
          </div>
        </div>
        <div class="oauth-field">
          <label class="input-label" for="oauth-url">{{ t('auth.oauth.fullUrl') }}</label>
          <div class="oauth-copy-row">
            <input id="oauth-url" class="input oauth-code-input oauth-url-input" :value="fullUrl" readonly />
            <button class="btn btn-secondary" type="button" :disabled="!fullUrl" @click="copy(fullUrl)">{{ t('common.copy') }}</button>
          </div>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useClipboard } from '@/composables/useClipboard'
import { useAppStore, useAuthStore } from '@/stores'
import AuthLayout from '@/components/layout/AuthLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { apiClient } from '@/api/client'
import { buildApiUrl } from '@/api/url'
import {
  exchangePendingOAuthCompletion,
  persistOAuthTokenContext,
  type OAuthTokenResponse
} from '@/api/auth'
import {
  clearAllAffiliateReferralCodes,
  loadOAuthAffiliateCode,
  oauthAffiliatePayload
} from '@/utils/oauthAffiliate'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const { copyToClipboard } = useClipboard()
const appStore = useAppStore()
const authStore = useAuthStore()
const isProcessing = ref(false)
const isSubmitting = ref(false)
const needsRegistrationCompletion = ref(false)
const invitationRequired = ref(false)
const registrationEmail = ref('')
const password = ref('')
const confirmPassword = ref('')
const invitationCode = ref('')
const registrationError = ref('')
const pendingProvider = ref<'github' | 'google'>('github')
const redirectTo = ref('/dashboard')
const invalidCallback = ref(false)
const EMAIL_OAUTH_PENDING_PROVIDER_KEY = 'email_oauth_pending_provider'

type EmailOAuthPendingCompletion = Partial<OAuthTokenResponse> & {
  error?: string
  provider?: string
  redirect?: string
  email?: string
  resolved_email?: string
  invitation_required?: boolean
}

const code = computed(() => (route.query.code as string) || '')
const state = computed(() => (route.query.state as string) || '')
const error = computed(
  () => (route.query.error as string) || (route.query.error_description as string) || ''
)

const fullUrl = computed(() => {
  if (typeof window === 'undefined') return ''
  return window.location.href
})
const providerName = computed(() =>
  pendingProvider.value === 'google' ? 'Google' : 'GitHub'
)
const registrationHint = computed(() =>
  invitationRequired.value
    ? t('auth.oidc.invitationRequired', { providerName: providerName.value })
    : t('auth.oidc.completeRegistration')
)
const canSubmitRegistration = computed(() => {
  if (!registrationEmail.value.trim()) return false
  if (password.value.length < 6) return false
  if (password.value !== confirmPassword.value) return false
  if (invitationRequired.value && !invitationCode.value.trim()) return false
  return true
})

function parseFragmentParams(): URLSearchParams {
  const raw = typeof window !== 'undefined' ? window.location.hash : ''
  const hash = raw.startsWith('#') ? raw.slice(1) : raw
  return new URLSearchParams(hash)
}

function readTokenResponse(params: URLSearchParams): OAuthTokenResponse | null {
  const accessToken = params.get('access_token')?.trim() || ''
  if (!accessToken) return null

  const response: OAuthTokenResponse = { access_token: accessToken }
  const refreshToken = params.get('refresh_token')?.trim() || ''
  if (refreshToken) response.refresh_token = refreshToken
  const expiresIn = Number.parseInt(params.get('expires_in')?.trim() || '', 10)
  if (Number.isFinite(expiresIn) && expiresIn > 0) response.expires_in = expiresIn
  const tokenType = params.get('token_type')?.trim() || ''
  if (tokenType) response.token_type = tokenType
  return response
}

function sanitizeRedirectPath(path: string | null | undefined): string {
  if (!path) return '/dashboard'
  if (!path.startsWith('/')) return '/dashboard'
  if (path.startsWith('//')) return '/dashboard'
  if (path.includes('://')) return '/dashboard'
  if (path.includes('\n') || path.includes('\r')) return '/dashboard'
  return path
}

function readPendingEmailOAuthProvider(): 'github' | 'google' | null {
  if (typeof window === 'undefined') return null
  const provider = window.sessionStorage.getItem(EMAIL_OAUTH_PENDING_PROVIDER_KEY)
  if (provider === 'github' || provider === 'google') return provider
  return null
}

function redirectProviderCallbackToBackend(provider: 'github' | 'google'): void {
  if (typeof window === 'undefined') return
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(route.query)) {
    if (Array.isArray(value)) {
      value.forEach((item) => {
        if (item != null) params.append(key, String(item))
      })
    } else if (value != null) {
      params.set(key, String(value))
    }
  }
  const suffix = params.toString() ? `?${params.toString()}` : ''
  window.location.href = buildApiUrl(`/auth/oauth/${provider}/callback${suffix}`)
}

async function finalizeTokenResponse(tokenResponse: OAuthTokenResponse, redirect: string) {
  persistOAuthTokenContext(tokenResponse)
  await authStore.setToken(tokenResponse.access_token)
  if (typeof window !== 'undefined') {
    window.sessionStorage.removeItem(EMAIL_OAUTH_PENDING_PROVIDER_KEY)
  }
  clearAllAffiliateReferralCodes()
  appStore.showSuccess(t('auth.loginSuccess'))
  await router.replace(sanitizeRedirectPath(redirect))
}

function hasOAuthTokenResponse(value: Partial<OAuthTokenResponse>): value is OAuthTokenResponse {
  return typeof value.access_token === 'string' && value.access_token.trim() !== ''
}

async function resumePendingEmailOAuth() {
  isProcessing.value = true
  try {
    const completion = await exchangePendingOAuthCompletion() as EmailOAuthPendingCompletion
    const completionRedirect = completion.redirect || '/dashboard'
    if (hasOAuthTokenResponse(completion)) {
      await finalizeTokenResponse(completion, completionRedirect)
      return
    }

    const provider = String(completion.provider || '').toLowerCase()
    if (provider === 'github' || provider === 'google') {
      pendingProvider.value = provider
    }
    redirectTo.value = sanitizeRedirectPath(completionRedirect)

    if (completion.error === 'invitation_required' || completion.error === 'registration_completion_required') {
      invitationRequired.value = completion.error === 'invitation_required' || completion.invitation_required === true
      registrationEmail.value = String(completion.resolved_email || completion.email || '').trim()
      needsRegistrationCompletion.value = true
      isProcessing.value = false
      return
    }

    appStore.showError(completion.error || t('auth.loginFailed'))
  } catch (e: unknown) {
    const err = e as { message?: string; response?: { data?: { message?: string } } }
    const message = err.response?.data?.message || err.message || t('auth.loginFailed')
    appStore.showError(message)
    invalidCallback.value = true
  } finally {
    if (!needsRegistrationCompletion.value) {
      isProcessing.value = false
    }
  }
}

async function handleSubmitRegistration() {
  registrationError.value = ''
  if (!registrationEmail.value.trim()) {
    registrationError.value = t('auth.emailRequired')
    return
  }
  if (password.value.length < 6) {
    registrationError.value = t('auth.passwordMinLength')
    return
  }
  if (password.value !== confirmPassword.value) {
    registrationError.value = t('auth.passwordsDoNotMatch')
    return
  }
  const code = invitationCode.value.trim()
  if (invitationRequired.value && !code) return

  isSubmitting.value = true
  try {
    const payload: { password: string; invitation_code?: string; aff_code?: string } = {
      password: password.value,
      ...oauthAffiliatePayload(loadOAuthAffiliateCode())
    }
    if (invitationRequired.value) {
      payload.invitation_code = code
    }
    const { data } = await apiClient.post<OAuthTokenResponse>(
      `/auth/oauth/${pendingProvider.value}/complete-registration`,
      payload
    )
    await finalizeTokenResponse(data, redirectTo.value)
  } catch (e: unknown) {
    const err = e as { message?: string; response?: { data?: { message?: string } } }
    registrationError.value =
      err.response?.data?.message || err.message || t('auth.oidc.completeRegistrationFailed')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  const params = parseFragmentParams()
  const tokenResponse = readTokenResponse(params)
  const fragmentError = params.get('error') || ''
  const fragmentErrorDescription =
    params.get('error_description') || params.get('error_message') || ''

  if (fragmentError) {
    appStore.showError(fragmentErrorDescription || fragmentError)
    return
  }
  if (!tokenResponse) {
    if (route.path === '/auth/oauth/callback') {
      const pendingEmailOAuthProvider = readPendingEmailOAuthProvider()
      if (pendingEmailOAuthProvider && code.value && state.value) {
        redirectProviderCallbackToBackend(pendingEmailOAuthProvider)
        return
      }
      await resumePendingEmailOAuth()
    }
    return
  }

  isProcessing.value = true
  try {
    await finalizeTokenResponse(tokenResponse, params.get('redirect') || '/dashboard')
  } catch (error: unknown) {
    const message = (error as { message?: string })?.message || t('auth.loginFailed')
    appStore.showError(message)
    isProcessing.value = false
  }
})

watch(
  error,
  (message) => {
    if (message) {
      appStore.showError(message)
    }
  },
  { immediate: true }
)

const copy = (value: string) => {
  if (!value) return
  copyToClipboard(value)
}
</script>

<style scoped>
.oauth-state {
  display: grid;
  width: 100%;
  gap: 1.5rem;
}

.oauth-state-centered {
  justify-items: center;
  padding: 0.25rem 0;
  text-align: center;
}

.oauth-state-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.oauth-state-icon {
  display: grid;
  width: 4rem;
  height: 4rem;
  flex: 0 0 4rem;
  place-items: center;
  border-radius: 1rem;
  color: var(--omnio-foreground);
  background: color-mix(in srgb, var(--omnio-foreground) 6%, transparent);
}

.oauth-state-icon-small {
  width: 3rem;
  height: 3rem;
  flex-basis: 3rem;
  border-radius: 0.75rem;
}

.oauth-state-icon-error {
  color: #ef4444;
  background: rgb(239 68 68 / 10%);
}

.oauth-heading {
  display: grid;
  gap: 0.5rem;
}

.oauth-heading-left {
  min-width: 0;
  text-align: left;
}

.oauth-heading h2 {
  margin: 0;
  color: var(--omnio-foreground);
  font-size: 1.5rem;
  font-weight: 650;
  line-height: 1.25;
  letter-spacing: -0.025em;
}

.oauth-heading-left h2 {
  font-size: 1.25rem;
}

.oauth-heading p,
.oauth-note {
  margin: 0;
  color: var(--omnio-muted);
  font-size: 0.875rem;
  line-height: 1.55;
}

.oauth-progress {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: var(--omnio-foreground);
  font-size: 0.875rem;
  font-weight: 580;
}

.oauth-note {
  max-width: 24rem;
  font-size: 0.78rem;
}

.oauth-form,
.oauth-code-list {
  display: grid;
  gap: 1rem;
}

.oauth-field .input-label {
  display: block;
  margin-bottom: 0.4rem;
  color: var(--omnio-foreground);
  font-size: 0.8rem;
  font-weight: 570;
}

.oauth-field .input {
  min-height: 2.5rem;
  border-color: var(--omnio-border-strong) !important;
  border-radius: 0.5rem !important;
  color: var(--omnio-foreground) !important;
  background: var(--omnio-surface) !important;
}

.oauth-field .input:focus {
  border-color: var(--omnio-primary) !important;
  box-shadow: 0 0 0 3px var(--omnio-ring) !important;
}

.oauth-field .input:disabled {
  color: var(--omnio-muted) !important;
  background: var(--omnio-surface-subtle) !important;
}

.oauth-error {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  margin: 0;
  padding: 0.75rem;
  border: 1px solid rgb(239 68 68 / 24%);
  border-radius: 0.625rem;
  color: #ef4444;
  background: rgb(239 68 68 / 8%);
  font-size: 0.8rem;
  line-height: 1.5;
}

.oauth-error svg {
  flex: 0 0 auto;
  margin-top: 0.08rem;
}

.oauth-primary-action {
  width: 100%;
  min-height: 2.5rem;
  gap: 0.5rem;
}

.oauth-inline-action {
  min-width: 8rem;
}

.oauth-copy-row {
  display: flex;
  gap: 0.5rem;
}

.oauth-code-input {
  min-width: 0;
  flex: 1;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace !important;
  font-size: 0.78rem !important;
}

.oauth-url-input {
  font-size: 0.7rem !important;
}

@media (max-width: 520px) {
  .oauth-state-header {
    display: grid;
  }

  .oauth-copy-row {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
  }
}
</style>
