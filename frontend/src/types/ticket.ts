export const TICKET_STATUSES = ['pending_admin', 'pending_user', 'resolved', 'closed'] as const
export const TICKET_CATEGORIES = ['account', 'billing', 'api_model', 'incident', 'feature', 'other'] as const
export const TICKET_PRIORITIES = ['low', 'normal', 'high', 'urgent'] as const

export type TicketStatus = (typeof TICKET_STATUSES)[number]
export type TicketCategory = (typeof TICKET_CATEGORIES)[number]
export type TicketPriority = (typeof TICKET_PRIORITIES)[number]
export type TicketSenderType = 'user' | 'admin'

export interface TicketConfig {
  user_center_enabled: boolean
  accept_new_tickets: boolean
}

export interface Ticket {
  id: number
  user_id: number
  subject: string
  category: TicketCategory
  status: TicketStatus
  priority: TicketPriority
  related_request_id: string
  assignee_id?: number
  created_at: string
  updated_at: string
  last_message_at: string
  closed_at?: string
  username?: string
  assignee_username?: string
  unread_count: number
}

export interface TicketMessage {
  id: number
  ticket_id: number
  sender_id: number
  sender_type: TicketSenderType
  content: string
  request_id: string
  created_at: string
}

export interface TicketDetail {
  ticket: Ticket
  messages: TicketMessage[]
}

export interface TicketPage {
  items: Ticket[]
  total: number
  page: number
  page_size: number
}

export interface TicketSummary {
  total: number
  pending_admin: number
  pending_user: number
  resolved: number
  closed: number
  unread_count: number
  pending_admin_count: number
  open_count: number
}

export interface TicketAssignee {
  id: number
  username: string
}

export interface TicketListParams {
  page?: number
  page_size?: number
  status?: TicketStatus | ''
  category?: TicketCategory | ''
  priority?: TicketPriority | ''
  user_id?: number
  assignee_id?: number
  related_request_id?: string
}

export interface CreateTicketPayload {
  subject: string
  category: TicketCategory
  content: string
  related_request_id?: string
}

export interface ReplyTicketPayload {
  content: string
  request_id?: string
}

export interface UpdateAdminTicketPayload {
  status?: TicketStatus
  priority?: TicketPriority
  assignee_id?: number
}

