export interface MembershipGroupBenefit {
  id: number
  level_id: number
  group_id: number
  group_name?: string
  allow_access: boolean
  pro_only: boolean
  rate_multiplier?: number | null
  rpm_limit?: number | null
}

export interface OmnioProGroupSetting {
  group_id: number
  rate_multiplier?: number | null
  pro_only: boolean
  daily_free_usd?: number | null
  monthly_free_usd?: number | null
}

export interface OmnioProQuotaProgress {
  group_id: number
  group_name: string
  daily_limit_usd: number
  daily_used_usd: number
  daily_remaining_usd: number
  monthly_limit_usd: number
  monthly_used_usd: number
  monthly_remaining_usd: number
}

export interface MembershipLevel {
  id: number
  name: string
  slug: string
  description: string
  rank: number
  badge_color: string
  concurrency_bonus: number
  priority_support: boolean
  active: boolean
  sort_order: number
  group_benefits?: MembershipGroupBenefit[]
  created_at?: string
  updated_at?: string
}

export interface MembershipOffer {
  id: number
  level_id: number
  level_name: string
  level_rank: number
  badge_color: string
  name: string
  description: string
  price: number
  original_price?: number | null
  currency: string
  duration_days: number
  for_sale: boolean
  sort_order: number
}

export interface MembershipGrant {
  id: number
  user_id: number
  level_id: number
  level_name: string
  level_rank: number
  badge_color: string
  source_type: string
  source_id: string
  starts_at: string
  expires_at: string
  status: string
  notes: string
}

export interface SubscriptionEntitlement {
  id: number
  order_id?: number | null
  subscription_id?: number | null
  plan_id?: number | null
  group_id: number
  plan_name: string
  daily_limit_usd?: number | null
  weekly_limit_usd?: number | null
  monthly_limit_usd?: number | null
  rate_multiplier: number
  overage_policy: 'block' | 'wallet'
  starts_at?: string | null
  expires_at?: string | null
  status: string
}

export interface MembershipSummary {
  effective_level?: MembershipLevel | null
  grants: MembershipGrant[]
  benefits: MembershipGroupBenefit[]
  offers: MembershipOffer[]
  quota_progress: OmnioProQuotaProgress[]
}

export interface MembershipPlanBenefit {
  plan_id: number
  plan_name: string
  level_id: number
  level_name: string
  duration_days?: number | null
}

export interface MembershipCatalog {
  levels: MembershipLevel[]
  offers: MembershipOffer[]
  plan_benefits: MembershipPlanBenefit[]
}

export interface MembershipAuditLog {
  id: number
  user_id?: number | null
  level_id?: number | null
  action: string
  source_type: string
  source_id: string
  operator: string
  detail: Record<string, unknown>
  created_at: string
}
