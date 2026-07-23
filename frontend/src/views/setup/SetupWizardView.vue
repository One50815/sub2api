<template>
  <div class="setup-page">
    <div class="setup-locale">
      <LocaleSwitcher />
    </div>

    <main class="setup-container">
      <header class="setup-brand">
        <div class="setup-logo-wrap">
          <img
            :src="setupLogo"
            :alt="t('setup.logoAlt')"
            class="setup-logo"
            @error="handleLogoError"
          />
        </div>
        <h1>{{ t('setup.initializeTitle', { siteName: appStore.siteName }) }}</h1>
        <p>{{ t('setup.description') }}</p>
      </header>

      <section class="setup-card">
        <header class="setup-card-header">
          <h2>{{ t('setup.wizardTitle') }}</h2>
          <p>{{ t('setup.wizardDescription') }}</p>
        </header>

        <div class="setup-card-content">
          <ol class="setup-steps" :aria-label="t('setup.progressLabel')">
            <li
              v-for="(step, index) in steps"
              :key="step.id"
              class="setup-step"
              :class="{
                'is-active': currentStep === index,
                'is-complete': currentStep > index
              }"
              :aria-current="currentStep === index ? 'step' : undefined"
            >
              <span class="setup-step-number" aria-hidden="true">
                <Icon v-if="currentStep > index" name="check" size="xs" :stroke-width="2.4" />
                <span v-else>{{ index + 1 }}</span>
              </span>
              <span class="setup-step-copy">
                <strong>{{ step.title }}</strong>
                <small>{{ step.description }}</small>
              </span>
            </li>
          </ol>

          <div class="setup-panel">
            <div v-if="currentStep === 0" class="setup-panel-inner">
              <header class="setup-panel-heading">
                <h3>{{ t('setup.database.title') }}</h3>
                <p>{{ t('setup.database.description') }}</p>
              </header>

              <div class="setup-form-grid">
                <div class="setup-field">
                  <label class="input-label" for="setup-db-host">{{ t('setup.database.host') }}</label>
                  <input id="setup-db-host" v-model="formData.database.host" type="text" class="input" placeholder="localhost" autocomplete="off" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-db-port">{{ t('setup.database.port') }}</label>
                  <input id="setup-db-port" v-model.number="formData.database.port" type="number" class="input" placeholder="5432" inputmode="numeric" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-db-user">{{ t('setup.database.username') }}</label>
                  <input id="setup-db-user" v-model="formData.database.user" type="text" class="input" placeholder="postgres" autocomplete="username" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-db-password">{{ t('setup.database.password') }}</label>
                  <input id="setup-db-password" v-model="formData.database.password" type="password" class="input" :placeholder="t('setup.database.passwordPlaceholder')" autocomplete="new-password" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-db-name">{{ t('setup.database.databaseName') }}</label>
                  <input id="setup-db-name" v-model="formData.database.dbname" type="text" class="input" placeholder="sub2api" autocomplete="off" />
                </div>
                <div class="setup-field">
                  <label class="input-label">{{ t('setup.database.sslMode') }}</label>
                  <Select
                    v-model="formData.database.sslmode"
                    :options="[
                      { value: 'disable', label: t('setup.database.ssl.disable') },
                      { value: 'require', label: t('setup.database.ssl.require') },
                      { value: 'verify-ca', label: t('setup.database.ssl.verifyCa') },
                      { value: 'verify-full', label: t('setup.database.ssl.verifyFull') }
                    ]"
                  />
                </div>
              </div>

              <button type="button" class="btn btn-secondary setup-test-button" :disabled="testingDb" @click="testDatabaseConnection">
                <Icon v-if="testingDb" name="refresh" size="sm" class="animate-spin" :stroke-width="2" />
                <Icon v-else-if="dbConnected" name="checkCircle" size="sm" class="setup-success-icon" :stroke-width="2" />
                {{ testingDb ? t('setup.status.testing') : dbConnected ? t('setup.status.success') : t('setup.status.testConnection') }}
              </button>
            </div>

            <div v-else-if="currentStep === 1" class="setup-panel-inner">
              <header class="setup-panel-heading">
                <h3>{{ t('setup.redis.title') }}</h3>
                <p>{{ t('setup.redis.description') }}</p>
              </header>

              <div class="setup-form-grid">
                <div class="setup-field">
                  <label class="input-label" for="setup-redis-host">{{ t('setup.redis.host') }}</label>
                  <input id="setup-redis-host" v-model="formData.redis.host" type="text" class="input" placeholder="localhost" autocomplete="off" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-redis-port">{{ t('setup.redis.port') }}</label>
                  <input id="setup-redis-port" v-model.number="formData.redis.port" type="number" class="input" placeholder="6379" inputmode="numeric" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-redis-password">{{ t('setup.redis.password') }}</label>
                  <input id="setup-redis-password" v-model="formData.redis.password" type="password" class="input" :placeholder="t('setup.redis.passwordPlaceholder')" autocomplete="new-password" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-redis-database">{{ t('setup.redis.database') }}</label>
                  <input id="setup-redis-database" v-model.number="formData.redis.db" type="number" class="input" placeholder="0" inputmode="numeric" />
                </div>
              </div>

              <div class="setup-toggle-row">
                <div>
                  <strong>{{ t('setup.redis.enableTls') }}</strong>
                  <span>{{ t('setup.redis.enableTlsHint') }}</span>
                </div>
                <Toggle v-model="formData.redis.enable_tls" />
              </div>

              <button type="button" class="btn btn-secondary setup-test-button" :disabled="testingRedis" @click="testRedisConnection">
                <Icon v-if="testingRedis" name="refresh" size="sm" class="animate-spin" :stroke-width="2" />
                <Icon v-else-if="redisConnected" name="checkCircle" size="sm" class="setup-success-icon" :stroke-width="2" />
                {{ testingRedis ? t('setup.status.testing') : redisConnected ? t('setup.status.success') : t('setup.status.testConnection') }}
              </button>
            </div>

            <div v-else-if="currentStep === 2" class="setup-panel-inner">
              <header class="setup-panel-heading">
                <h3>{{ t('setup.admin.title') }}</h3>
                <p>{{ t('setup.admin.description') }}</p>
              </header>

              <div class="setup-form-grid">
                <div class="setup-field">
                  <label class="input-label" for="setup-admin-email">{{ t('setup.admin.email') }}</label>
                  <input id="setup-admin-email" v-model="formData.admin.email" type="email" class="input" placeholder="admin@example.com" autocomplete="username" />
                </div>
                <div class="setup-field">
                  <label class="input-label" for="setup-admin-password">{{ t('setup.admin.password') }}</label>
                  <input id="setup-admin-password" v-model="formData.admin.password" type="password" class="input" :placeholder="t('setup.admin.passwordPlaceholder')" autocomplete="new-password" />
                </div>
                <div class="setup-field setup-field-wide">
                  <label class="input-label" for="setup-admin-confirm">{{ t('setup.admin.confirmPassword') }}</label>
                  <input id="setup-admin-confirm" v-model="confirmPassword" type="password" class="input" :aria-invalid="Boolean(confirmPassword && formData.admin.password !== confirmPassword)" :placeholder="t('setup.admin.confirmPasswordPlaceholder')" autocomplete="new-password" />
                  <p v-if="confirmPassword && formData.admin.password !== confirmPassword" class="input-error-text">
                    {{ t('setup.admin.passwordMismatch') }}
                  </p>
                </div>
              </div>
            </div>

            <div v-else class="setup-panel-inner setup-ready">
              <div class="setup-ready-icon">
                <Icon name="checkCircle" size="xl" :stroke-width="1.8" />
              </div>
              <header class="setup-panel-heading">
                <h3>{{ t('setup.ready.title') }}</h3>
                <p>{{ t('setup.ready.description') }}</p>
              </header>

              <dl class="setup-summary">
                <div>
                  <dt>{{ t('setup.ready.database') }}</dt>
                  <dd>{{ formData.database.user }}@{{ formData.database.host }}:{{ formData.database.port }}/{{ formData.database.dbname }}</dd>
                </div>
                <div>
                  <dt>{{ t('setup.ready.redis') }}</dt>
                  <dd>{{ formData.redis.host }}:{{ formData.redis.port }}</dd>
                </div>
                <div>
                  <dt>{{ t('setup.ready.adminEmail') }}</dt>
                  <dd>{{ formData.admin.email }}</dd>
                </div>
              </dl>
            </div>
          </div>

          <div v-if="errorMessage" class="setup-alert setup-alert-error" role="alert">
            <Icon name="exclamationCircle" size="md" />
            <p>{{ errorMessage }}</p>
          </div>

          <div v-if="installSuccess" class="setup-alert setup-alert-success" role="status">
            <Icon v-if="!serviceReady" name="refresh" size="md" class="animate-spin" />
            <Icon v-else name="checkCircle" size="md" />
            <div>
              <strong>{{ t('setup.status.completed') }}</strong>
              <p>{{ serviceReady ? t('setup.status.redirecting') : t('setup.status.restarting') }}</p>
            </div>
          </div>
        </div>

        <footer v-if="!installSuccess" class="setup-card-footer">
          <div>
            <button v-if="currentStep > 0" type="button" class="btn btn-secondary" @click="currentStep--">
              {{ t('common.back') }}
            </button>
          </div>
          <button v-if="currentStep < 3" type="button" class="btn btn-primary" :disabled="!canProceed" @click="nextStep">
            {{ t('common.next') }}
          </button>
          <button v-else type="button" class="btn btn-primary" :disabled="installing" @click="performInstall">
            <Icon v-if="installing" name="refresh" size="sm" class="animate-spin" :stroke-width="2" />
            <Icon v-else name="checkCircle" size="sm" :stroke-width="2" />
            {{ installing ? t('setup.status.installing') : t('setup.status.completeInstallation') }}
          </button>
        </footer>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { testDatabase, testRedis, install, type InstallRequest } from '@/api/setup'
import { buildGatewayUrl } from '@/api/client'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const appStore = useAppStore()
const defaultLogo = '/assets/brand/omnio-mark.svg?v=3'
const setupLogo = computed(() => appStore.siteLogo || defaultLogo)

const steps = computed(() => [
  {
    id: 'database',
    title: t('setup.database.title'),
    description: t('setup.database.stepDescription')
  },
  {
    id: 'redis',
    title: t('setup.redis.title'),
    description: t('setup.redis.stepDescription')
  },
  {
    id: 'admin',
    title: t('setup.admin.title'),
    description: t('setup.admin.stepDescription')
  },
  {
    id: 'complete',
    title: t('setup.ready.title'),
    description: t('setup.ready.stepDescription')
  }
])

const currentStep = ref(0)
const errorMessage = ref('')
const installSuccess = ref(false)

// Connection test states
const testingDb = ref(false)
const testingRedis = ref(false)
const dbConnected = ref(false)
const redisConnected = ref(false)
const installing = ref(false)
const confirmPassword = ref('')
const serviceReady = ref(false)

function handleLogoError(event: Event) {
  const image = event.currentTarget as HTMLImageElement
  if (!image.src.includes('/assets/brand/omnio-mark.svg')) {
    image.src = defaultLogo
  }
}

// Default server port
const getCurrentPort = (): number => {
  const port = window.location.port
  if (port) {
    return parseInt(port, 10)
  }

  return window.location.protocol === 'https:' ? 443 : 80
}

const formData = reactive<InstallRequest>({
  database: {
    host: 'localhost',
    port: 5432,
    user: 'postgres',
    password: '',
    dbname: 'sub2api',
    sslmode: 'disable'
  },
  redis: {
    host: 'localhost',
    port: 6379,
    password: '',
    db: 0,
    enable_tls: false
  },
  admin: {
    email: '',
    password: ''
  },
  server: {
    host: '0.0.0.0',
    port: getCurrentPort(), // Use current port from browser
    mode: 'release'
  }
})

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 0:
      return dbConnected.value
    case 1:
      return redisConnected.value
    case 2:
      return (
        formData.admin.email &&
        formData.admin.password.length >= 8 &&
        formData.admin.password === confirmPassword.value
      )
    default:
      return true
  }
})

async function testDatabaseConnection() {
  testingDb.value = true
  errorMessage.value = ''
  dbConnected.value = false

  try {
    await testDatabase(formData.database)
    dbConnected.value = true
  } catch (error: unknown) {
    const err = error as { response?: { data?: { detail?: string; message?: string } }; message?: string }
    errorMessage.value =
      err.response?.data?.detail || err.response?.data?.message || err.message || 'Connection failed'
  } finally {
    testingDb.value = false
  }
}

async function testRedisConnection() {
  testingRedis.value = true
  errorMessage.value = ''
  redisConnected.value = false

  try {
    await testRedis(formData.redis)
    redisConnected.value = true
  } catch (error: unknown) {
    const err = error as { response?: { data?: { detail?: string; message?: string } }; message?: string }
    errorMessage.value =
      err.response?.data?.detail || err.response?.data?.message || err.message || 'Connection failed'
  } finally {
    testingRedis.value = false
  }
}

function nextStep() {
  if (canProceed.value) {
    errorMessage.value = ''
    currentStep.value++
  }
}

async function performInstall() {
  installing.value = true
  errorMessage.value = ''

  try {
    await install(formData)
    installSuccess.value = true
    // Start polling for service restart
    waitForServiceRestart()
  } catch (error: unknown) {
    const err = error as { response?: { data?: { detail?: string; message?: string } }; message?: string }
    errorMessage.value =
      err.response?.data?.detail || err.response?.data?.message || err.message || 'Installation failed'
  } finally {
    installing.value = false
  }
}

// Wait for service to restart and become available
async function waitForServiceRestart() {
  const maxAttempts = 60 // Increase to 60 attempts, ~60 seconds max
  const interval = 1000 // 1 second between attempts

  // Wait a moment for the service to start restarting
  await new Promise((resolve) => setTimeout(resolve, 3000))

  for (let attempt = 0; attempt < maxAttempts; attempt++) {
    try {
      // Use setup status endpoint as it tells us the real mode
      // Service might return 404 or connection refused while restarting
      const response = await fetch(buildGatewayUrl('/setup/status'), {
        method: 'GET',
        cache: 'no-store'
      })

      if (response.ok) {
        const data = await response.json()
        // If needs_setup is false, service has restarted in normal mode
        if (data.data && !data.data.needs_setup) {
          serviceReady.value = true
          // Redirect to login page after a short delay
          setTimeout(() => {
            window.location.href = '/login'
          }, 1500)
          return
        }
      }
    } catch {
      // Service not ready or network error during restart, continue polling
    }

    await new Promise((resolve) => setTimeout(resolve, interval))
  }

  // If we reach here, service didn't restart in time
  // Show a message to refresh manually
  errorMessage.value = t('setup.status.timeout')
}
</script>

<style scoped>
.setup-page {
  position: relative;
  min-height: 100svh;
  padding: 2.5rem 1rem;
  color: var(--omnio-foreground);
  background: color-mix(in srgb, var(--omnio-surface-subtle) 78%, var(--omnio-bg));
}

.setup-locale {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  z-index: 10;
}

.setup-container {
  width: 100%;
  max-width: 64rem;
  margin: 0 auto;
}

.setup-brand {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2rem;
  text-align: center;
}

.setup-logo-wrap {
  display: grid;
  width: 3rem;
  height: 3rem;
  place-items: center;
  overflow: hidden;
  border: 1px solid var(--omnio-border);
  border-radius: 999px;
  background: var(--omnio-surface);
  box-shadow: 0 1px 2px rgb(15 23 42 / 8%);
}

.setup-logo {
  width: 2rem;
  height: 2rem;
  object-fit: contain;
}

.setup-brand h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 650;
  line-height: 1.25;
  letter-spacing: -0.025em;
}

.setup-brand p,
.setup-card-header p,
.setup-panel-heading p {
  margin: 0;
  color: var(--omnio-muted);
}

.setup-brand p {
  max-width: 36rem;
  font-size: 0.925rem;
}

.setup-card {
  overflow: hidden;
  border: 1px solid var(--omnio-border);
  border-radius: 0.875rem;
  background: var(--omnio-surface);
  box-shadow: var(--omnio-card-shadow);
}

.setup-card-header {
  padding: 1.5rem 1.5rem 0;
}

.setup-card-header h2,
.setup-panel-heading h3 {
  margin: 0;
  color: var(--omnio-foreground);
  font-weight: 650;
  letter-spacing: -0.018em;
}

.setup-card-header h2 {
  font-size: 1.25rem;
  line-height: 1.5;
}

.setup-card-header p {
  margin-top: 0.25rem;
  font-size: 0.875rem;
}

.setup-card-content {
  display: grid;
  gap: 1.5rem;
  padding: 1.5rem;
}

.setup-steps {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.75rem;
  margin: 0;
  padding: 0;
  list-style: none;
}

.setup-step {
  display: flex;
  min-width: 0;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.75rem;
  border: 1px solid var(--omnio-border);
  border-radius: 0.75rem;
  background: var(--omnio-surface);
  transition: border-color 150ms ease, box-shadow 150ms ease, background-color 150ms ease;
}

.setup-step.is-active {
  border-color: var(--omnio-primary);
  box-shadow: 0 0 0 2px var(--omnio-ring);
}

.setup-step.is-complete {
  border-color: color-mix(in srgb, var(--omnio-primary) 42%, var(--omnio-border));
  background: color-mix(in srgb, var(--omnio-primary) 5%, var(--omnio-surface));
}

.setup-step-number {
  display: grid;
  width: 1.5rem;
  height: 1.5rem;
  flex: 0 0 1.5rem;
  place-items: center;
  border: 1px solid var(--omnio-border-strong);
  border-radius: 0.375rem;
  color: var(--omnio-muted);
  font-size: 0.75rem;
  font-weight: 650;
}

.setup-step.is-active .setup-step-number,
.setup-step.is-complete .setup-step-number {
  border-color: var(--omnio-primary);
  color: #fff;
  background: var(--omnio-primary);
}

.setup-step-copy {
  min-width: 0;
}

.setup-step-copy strong,
.setup-step-copy small {
  display: block;
}

.setup-step-copy strong {
  color: var(--omnio-foreground);
  font-size: 0.82rem;
  font-weight: 630;
  line-height: 1.25;
}

.setup-step-copy small {
  margin-top: 0.25rem;
  color: var(--omnio-muted);
  font-size: 0.72rem;
  line-height: 1.35;
}

.setup-panel {
  padding-top: 1.5rem;
  border-top: 1px solid var(--omnio-border);
}

.setup-panel-inner {
  display: grid;
  gap: 1.5rem;
}

.setup-panel-heading h3 {
  font-size: 1rem;
  line-height: 1.5;
}

.setup-panel-heading p {
  margin-top: 0.2rem;
  font-size: 0.82rem;
  line-height: 1.5;
}

.setup-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.setup-field {
  min-width: 0;
}

.setup-field-wide {
  grid-column: 1 / -1;
}

.setup-field .input-label {
  display: block;
  margin-bottom: 0.4rem;
  color: var(--omnio-foreground);
  font-size: 0.8rem;
  font-weight: 560;
}

.setup-field .input,
.setup-field :deep(.input) {
  min-height: 2.5rem;
  border-color: var(--omnio-border-strong) !important;
  border-radius: 0.5rem !important;
  color: var(--omnio-foreground) !important;
  background: var(--omnio-surface) !important;
  box-shadow: 0 1px 1px rgb(15 23 42 / 2%);
}

.setup-field .input:focus,
.setup-field :deep(.input:focus) {
  border-color: var(--omnio-primary) !important;
  box-shadow: 0 0 0 3px var(--omnio-ring) !important;
}

.setup-toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.875rem 1rem;
  border: 1px solid var(--omnio-border);
  border-radius: 0.75rem;
  background: var(--omnio-surface);
}

.setup-toggle-row strong,
.setup-toggle-row span {
  display: block;
}

.setup-toggle-row strong {
  font-size: 0.82rem;
  font-weight: 620;
}

.setup-toggle-row span {
  margin-top: 0.15rem;
  color: var(--omnio-muted);
  font-size: 0.75rem;
}

.setup-test-button {
  width: 100%;
  min-height: 2.5rem;
  gap: 0.5rem;
}

.setup-success-icon {
  color: #10b981;
}

.setup-ready {
  justify-items: center;
  text-align: center;
}

.setup-ready-icon {
  display: grid;
  width: 4rem;
  height: 4rem;
  place-items: center;
  border-radius: 1rem;
  color: #059669;
  background: rgb(16 185 129 / 10%);
}

.setup-ready .setup-panel-heading p {
  max-width: 34rem;
}

.setup-summary {
  width: 100%;
  margin: 0;
  padding: 0.25rem 1.5rem;
  border: 1px solid var(--omnio-border);
  border-radius: 0.75rem;
  text-align: left;
  background: var(--omnio-surface);
  box-shadow: 0 1px 2px rgb(15 23 42 / 3%);
}

.setup-summary > div {
  padding: 1rem 0;
}

.setup-summary > div + div {
  border-top: 1px solid var(--omnio-border);
}

.setup-summary dt {
  color: var(--omnio-muted);
  font-size: 0.7rem;
  font-weight: 650;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.setup-summary dd {
  margin: 0.35rem 0 0;
  overflow-wrap: anywhere;
  color: var(--omnio-foreground);
  font-size: 0.85rem;
  font-weight: 600;
}

.setup-alert {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border: 1px solid;
  border-radius: 0.75rem;
  font-size: 0.82rem;
  line-height: 1.5;
}

.setup-alert svg {
  flex: 0 0 auto;
  margin-top: 0.05rem;
}

.setup-alert p {
  margin: 0;
}

.setup-alert strong {
  display: block;
  margin-bottom: 0.1rem;
  font-weight: 650;
}

.setup-alert-error {
  border-color: rgb(239 68 68 / 28%);
  color: #b91c1c;
  background: rgb(254 242 242 / 82%);
}

.setup-alert-success {
  border-color: rgb(16 185 129 / 28%);
  color: #047857;
  background: rgb(236 253 245 / 82%);
}

:global(.dark) .setup-alert-error {
  color: #fca5a5;
  background: rgb(127 29 29 / 18%);
}

:global(.dark) .setup-alert-success {
  color: #6ee7b7;
  background: rgb(6 78 59 / 22%);
}

.setup-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  min-height: 4.5rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--omnio-border);
  background: color-mix(in srgb, var(--omnio-surface-subtle) 42%, var(--omnio-surface));
}

.setup-card-footer .btn {
  min-height: 2.25rem;
  gap: 0.45rem;
  border-radius: 0.5rem;
}

@media (max-width: 767px) {
  .setup-page {
    padding: 4.5rem 1rem 1.5rem;
  }

  .setup-locale {
    top: 1rem;
    right: 1rem;
  }

  .setup-brand {
    margin-bottom: 1.5rem;
  }

  .setup-steps {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 639px) {
  .setup-card-header,
  .setup-card-content,
  .setup-card-footer {
    padding-right: 1rem;
    padding-left: 1rem;
  }

  .setup-form-grid {
    grid-template-columns: 1fr;
  }

  .setup-field-wide {
    grid-column: auto;
  }

  .setup-card-footer {
    align-items: stretch;
  }

  .setup-card-footer .btn {
    min-width: 5.5rem;
  }
}
</style>
