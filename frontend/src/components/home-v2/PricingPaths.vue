<template>
  <section id="pricing" class="ov2-section ov2-pricing-section" data-home-reveal>
    <div class="ov2-shell">
      <div class="ov2-section-intro ov2-section-intro-split">
        <div>
          <span class="ov2-kicker"><i></i>05 / {{ t('homeV2.pricingKicker') }}</span>
          <h2>{{ t('homeV2.pricingTitle') }}</h2>
        </div>
        <p>{{ t('homeV2.pricingDescription') }}</p>
      </div>

      <div class="ov2-pricing-paths">
        <article class="ov2-pricing-card is-personal">
          <span class="ov2-pricing-index">01 / PERSONAL</span>
          <div class="ov2-pricing-type">
            <div><strong>{{ t('homeV2.personal') }}</strong><small>{{ t('homeV2.personalAudience') }}</small></div>
          </div>
          <div class="ov2-pricing-scenario"><i></i>{{ t('homeV2.personalScenario') }}</div>
          <div class="ov2-pricing-offers">
            <span>Free</span><span>Pro</span><span>Max</span>
          </div>
          <p>{{ t('homeV2.personalCopy') }}</p>
          <router-link to="/omnio-pro">{{ t('homeV2.viewPlans') }} <span>→</span></router-link>
        </article>

        <article class="ov2-pricing-card is-developer">
          <span class="ov2-pricing-index">02 / API</span>
          <div class="ov2-pricing-type">
            <div><strong>{{ t('homeV2.developers') }}</strong><small>{{ t('homeV2.developerAudience') }}</small></div>
          </div>
          <div class="ov2-pricing-scenario"><i></i>{{ t('homeV2.developerScenario') }}</div>
          <div class="ov2-token-rate">
            <span><small>{{ t('homeV2.inputTokens') }}</small><strong>$0.10<em>/ 1M</em></strong></span>
            <span><small>{{ t('homeV2.outputTokens') }}</small><strong>$0.40<em>/ 1M</em></strong></span>
          </div>
          <p>{{ t('homeV2.developerPricingCopy') }}</p>
          <a href="/docs/pricing/">{{ t('homeV2.livePricing') }} <span>→</span></a>
        </article>

        <article class="ov2-pricing-card is-enterprise">
          <span class="ov2-pricing-index">03 / ENTERPRISE</span>
          <div class="ov2-pricing-type">
            <div><strong>{{ t('homeV2.enterprise') }}</strong><small>{{ t('homeV2.enterpriseAudience') }}</small></div>
          </div>
          <div class="ov2-pricing-scenario"><i></i>{{ t('homeV2.enterpriseScenario') }}</div>
          <div class="ov2-enterprise-features">
            <span>{{ t('homeV2.teamManagement') }}</span>
            <span>{{ t('homeV2.quotaControl') }}</span>
            <span>{{ t('homeV2.dedicatedSupport') }}</span>
          </div>
          <p>{{ t('homeV2.enterpriseCopy') }}</p>
          <router-link to="/register">{{ t('homeV2.contactTeam') }} <span>→</span></router-link>
        </article>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
