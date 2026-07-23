<template>
  <main class="not-found-page">
    <section class="not-found-content">
      <h1 aria-label="404">404</h1>
      <strong>{{ t('errors.pageNotFoundTitle') }}</strong>
      <p>
        {{ t('errors.pageNotFoundLineOne') }}<br />
        {{ t('errors.pageNotFoundLineTwo') }}
      </p>
      <div class="not-found-actions">
        <button type="button" class="btn btn-secondary" @click="goBack">
          {{ t('errors.goBack') }}
        </button>
        <RouterLink to="/home" class="btn btn-primary">
          {{ t('errors.backToHome') }}
        </RouterLink>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()

function goBack(): void {
  router.back()
}
</script>

<style scoped>
.not-found-page {
  display: grid;
  min-height: 100svh;
  place-items: center;
  padding: 1.5rem;
  color: var(--omnio-foreground);
  background: var(--omnio-bg);
}

.not-found-content {
  display: flex;
  width: 100%;
  max-width: 32rem;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  text-align: center;
}

.not-found-content h1 {
  margin: 0;
  color: var(--omnio-foreground);
  font-size: clamp(5.5rem, 20vw, 7rem);
  font-weight: 750;
  line-height: 1;
  letter-spacing: -0.07em;
}

.not-found-content strong {
  color: var(--omnio-foreground);
  font-size: 1rem;
  font-weight: 620;
}

.not-found-content p {
  margin: 0;
  color: var(--omnio-muted);
  font-size: 0.9rem;
  line-height: 1.6;
}

.not-found-actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
  margin-top: 1rem;
}

.not-found-actions .btn {
  min-width: 7.5rem;
  min-height: 2.4rem;
  border-radius: 0.5rem;
}

@media (max-width: 420px) {
  .not-found-actions {
    width: 100%;
    flex-direction: column-reverse;
    gap: 0.65rem;
  }

  .not-found-actions .btn {
    width: 100%;
  }
}
</style>
