import { apiClient } from './client'
import type { MembershipAuditLog, MembershipCatalog, MembershipLevel, MembershipOffer, MembershipSummary, MembershipGroupBenefit, MembershipGrant, MembershipPlanBenefit } from '@/types/membership'

export const membershipAPI = {
  getSummary() {
    return apiClient.get<MembershipSummary>('/membership/summary')
  },
  getOffers() {
    return apiClient.get<MembershipOffer[]>('/membership/offers')
  },
  setOveragePolicy(id: number, policy: 'block' | 'wallet') {
    return apiClient.put(`/membership/subscription-entitlements/${id}/overage-policy`, { policy })
  },
  adminCatalog() {
    return apiClient.get<MembershipCatalog>('/admin/membership/catalog')
  },
  saveLevel(level: Partial<MembershipLevel> & { id?: number }) {
    return level.id ? apiClient.put<MembershipLevel>(`/admin/membership/levels/${level.id}`, level) : apiClient.post<MembershipLevel>('/admin/membership/levels', level)
  },
  deleteLevel(id: number) {
    return apiClient.delete(`/admin/membership/levels/${id}`)
  },
  saveOffer(offer: Partial<MembershipOffer> & { id?: number }) {
    return offer.id ? apiClient.put<MembershipOffer>(`/admin/membership/offers/${offer.id}`, offer) : apiClient.post<MembershipOffer>('/admin/membership/offers', offer)
  },
  deleteOffer(id: number) {
    return apiClient.delete(`/admin/membership/offers/${id}`)
  },
  saveBenefit(benefit: Partial<MembershipGroupBenefit>) {
    return apiClient.put('/admin/membership/benefits', benefit)
  },
  deleteBenefit(levelId: number, groupId: number) {
    return apiClient.delete(`/admin/membership/benefits/${levelId}/${groupId}`)
  },
  savePlanBenefit(benefit: MembershipPlanBenefit) {
    return apiClient.put('/admin/membership/plan-benefits', benefit)
  },
  deletePlanBenefit(planId: number) {
    return apiClient.delete(`/admin/membership/plan-benefits/${planId}`)
  },
  grant(data: { user_id: number; level_id: number; days: number; notes?: string }) {
    return apiClient.post<MembershipGrant>('/admin/membership/grants', data)
  },
  listGrants(params: { user_id?: number; status?: 'active' | 'revoked' | 'all'; limit?: number } = {}) {
    return apiClient.get<MembershipGrant[]>('/admin/membership/grants', { params })
  },
  revokeGrant(id: number) {
    return apiClient.post(`/admin/membership/grants/${id}/revoke`)
  },
  auditLogs(limit = 100) {
    return apiClient.get<MembershipAuditLog[]>('/admin/membership/audit-logs', { params: { limit } })
  }
}

export default membershipAPI
