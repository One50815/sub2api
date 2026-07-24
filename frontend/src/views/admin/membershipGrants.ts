import type { MembershipGrant } from '@/types/membership'

export type MembershipGrantFilterStatus = 'active' | 'revoked' | 'all'
export type MembershipGrantDisplayStatus = 'active' | 'scheduled' | 'expired' | 'revoked'

export function membershipGrantDisplayStatus(
  grant: MembershipGrant,
  now = Date.now()
): MembershipGrantDisplayStatus {
  if (grant.status === 'revoked') return 'revoked'
  if (new Date(grant.expires_at).getTime() <= now) return 'expired'
  if (new Date(grant.starts_at).getTime() > now) return 'scheduled'
  return 'active'
}

export function canRevokeMembershipGrant(grant: MembershipGrant, now = Date.now()): boolean {
  const status = membershipGrantDisplayStatus(grant, now)
  return status === 'active' || status === 'scheduled'
}

export function membershipGrantStatusLabel(grant: MembershipGrant, now = Date.now()): string {
  const labels: Record<MembershipGrantDisplayStatus, string> = {
    active: '生效中',
    scheduled: '待生效',
    expired: '已过期',
    revoked: '已撤销'
  }
  return labels[membershipGrantDisplayStatus(grant, now)]
}

export function membershipGrantSourceLabel(sourceType: string): string {
  const labels: Record<string, string> = {
    manual: '管理员赠送',
    payment_order: 'Omnio Pro 订单',
    subscription_order: '订阅赠送'
  }
  return labels[sourceType] || sourceType
}
