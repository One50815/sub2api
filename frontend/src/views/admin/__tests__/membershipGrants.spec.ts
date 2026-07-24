import { describe, expect, it } from 'vitest'
import type { MembershipGrant } from '@/types/membership'
import {
  canRevokeMembershipGrant,
  membershipGrantDisplayStatus,
  membershipGrantSourceLabel,
  membershipGrantStatusLabel
} from '../membershipGrants'

const now = new Date('2026-07-24T04:00:00Z').getTime()

function grant(overrides: Partial<MembershipGrant> = {}): MembershipGrant {
  return {
    id: 7,
    user_id: 42,
    level_id: 3,
    level_name: 'Omnio Pro Max',
    level_rank: 20,
    badge_color: '#2563eb',
    source_type: 'manual',
    source_id: 'manual:42:1',
    starts_at: '2026-07-23T04:00:00Z',
    expires_at: '2026-08-23T04:00:00Z',
    status: 'active',
    notes: '',
    ...overrides
  }
}

describe('membership grant presentation', () => {
  it('allows revoking active and scheduled grants', () => {
    expect(membershipGrantDisplayStatus(grant(), now)).toBe('active')
    expect(canRevokeMembershipGrant(grant(), now)).toBe(true)

    const scheduled = grant({ starts_at: '2026-07-25T04:00:00Z' })
    expect(membershipGrantStatusLabel(scheduled, now)).toBe('待生效')
    expect(canRevokeMembershipGrant(scheduled, now)).toBe(true)
  })

  it('does not offer revoke for expired or revoked records', () => {
    const expired = grant({ expires_at: '2026-07-24T03:59:59Z' })
    expect(membershipGrantDisplayStatus(expired, now)).toBe('expired')
    expect(canRevokeMembershipGrant(expired, now)).toBe(false)

    const revoked = grant({ status: 'revoked' })
    expect(membershipGrantStatusLabel(revoked, now)).toBe('已撤销')
    expect(canRevokeMembershipGrant(revoked, now)).toBe(false)
  })

  it('formats known and fallback grant sources', () => {
    expect(membershipGrantSourceLabel('manual')).toBe('管理员赠送')
    expect(membershipGrantSourceLabel('payment_order')).toBe('Omnio Pro 订单')
    expect(membershipGrantSourceLabel('custom')).toBe('custom')
  })
})
