import { apiClient } from './client'
import type {
  CreateTicketPayload,
  ReplyTicketPayload,
  Ticket,
  TicketAssignee,
  TicketConfig,
  TicketDetail,
  TicketListParams,
  TicketMessage,
  TicketPage,
  TicketSummary,
  UpdateAdminTicketPayload,
} from '@/types/ticket'

const cleanParams = (params: TicketListParams = {}) =>
  Object.fromEntries(Object.entries(params).filter(([, value]) => value !== '' && value !== undefined && value !== null))

export const ticketsAPI = {
  async getConfig(): Promise<TicketConfig> {
    const { data } = await apiClient.get<TicketConfig>('/tickets/config')
    return data
  },
  async list(params: TicketListParams = {}): Promise<TicketPage> {
    const { data } = await apiClient.get<TicketPage>('/tickets', { params: cleanParams(params) })
    return data
  },
  async create(payload: CreateTicketPayload): Promise<TicketDetail> {
    const { data } = await apiClient.post<TicketDetail>('/tickets', payload)
    return data
  },
  async summary(): Promise<TicketSummary> {
    const { data } = await apiClient.get<TicketSummary>('/tickets/summary')
    return data
  },
  async get(id: number): Promise<TicketDetail> {
    const { data } = await apiClient.get<TicketDetail>(`/tickets/${id}`)
    return data
  },
  async markRead(id: number): Promise<void> {
    await apiClient.post(`/tickets/${id}/read`)
  },
  async reply(id: number, payload: ReplyTicketPayload): Promise<TicketMessage> {
    const { data } = await apiClient.post<TicketMessage>(`/tickets/${id}/messages`, payload)
    return data
  },
  async updateStatus(id: number, status: 'resolved' | 'closed'): Promise<Ticket> {
    const { data } = await apiClient.patch<Ticket>(`/tickets/${id}/status`, { status })
    return data
  },
  admin: {
    async getConfig(): Promise<TicketConfig> {
      const { data } = await apiClient.get<TicketConfig>('/admin/tickets/config')
      return data
    },
    async updateConfig(payload: TicketConfig): Promise<TicketConfig> {
      const { data } = await apiClient.put<TicketConfig>('/admin/tickets/config', payload)
      return data
    },
    async list(params: TicketListParams = {}): Promise<TicketPage> {
      const { data } = await apiClient.get<TicketPage>('/admin/tickets', { params: cleanParams(params) })
      return data
    },
    async summary(): Promise<TicketSummary> {
      const { data } = await apiClient.get<TicketSummary>('/admin/tickets/summary')
      return data
    },
    async assignees(): Promise<TicketAssignee[]> {
      const { data } = await apiClient.get<TicketAssignee[]>('/admin/tickets/assignees')
      return data
    },
    async get(id: number): Promise<TicketDetail> {
      const { data } = await apiClient.get<TicketDetail>(`/admin/tickets/${id}`)
      return data
    },
    async markRead(id: number): Promise<void> {
      await apiClient.post(`/admin/tickets/${id}/read`)
    },
    async reply(id: number, payload: ReplyTicketPayload): Promise<TicketMessage> {
      const { data } = await apiClient.post<TicketMessage>(`/admin/tickets/${id}/messages`, payload)
      return data
    },
    async update(id: number, payload: UpdateAdminTicketPayload): Promise<Ticket> {
      const { data } = await apiClient.patch<Ticket>(`/admin/tickets/${id}`, payload)
      return data
    },
  },
}

