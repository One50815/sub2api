export default {
  tickets: {
    title: 'My Tickets',
    description: 'Submit issues and continue the conversation with support.',
    newTicket: 'New ticket',
    intakePaused: 'New ticket intake is paused. You can still view and reply to existing tickets.',
    create: {
      title: 'Create ticket', subject: 'Subject', subjectPlaceholder: 'Briefly describe the issue', category: 'Category',
      content: 'Problem description', contentPlaceholder: 'Include steps, expected behavior, and actual behavior', requestId: 'Related request ID',
      requestIdHint: 'Optional. This helps administrators locate the matching request log.', submit: 'Submit ticket', success: 'Ticket created',
    },
    status: { pending_admin: 'Waiting for support', pending_user: 'Waiting for user', resolved: 'Resolved', closed: 'Closed' },
    category: { account: 'Account', billing: 'Billing & payment', api_model: 'API & models', incident: 'Service incident', feature: 'Feature request', other: 'Other' },
    priority: { low: 'Low', normal: 'Normal', high: 'High', urgent: 'Urgent' },
    summary: { total: 'All tickets', pendingAdmin: 'Waiting for support', pendingUser: 'Waiting for user', unread: 'Unread messages' },
    filters: { status: 'Status', category: 'Category', priority: 'Priority', userId: 'User ID', assignee: 'Assignee', requestId: 'Request ID', all: 'All', unassigned: 'Unassigned', apply: 'Apply', reset: 'Reset' },
    table: { subject: 'Subject', user: 'User', category: 'Category', status: 'Status', priority: 'Priority', assignee: 'Assignee', updatedAt: 'Updated', unread: '{count} unread' },
    empty: 'No tickets yet', emptyDescription: 'Create your first ticket when you need help.', loadFailed: 'Failed to load tickets', notFound: 'Ticket not found or unavailable', back: 'Back to tickets',
    detail: { title: 'Ticket details', requestId: 'Related request ID', createdAt: 'Created', updatedAt: 'Updated', conversation: 'Conversation', userMessage: 'You', adminMessage: 'Support', userLabel: 'User', adminLabel: 'Administrator', messageRequestId: 'Request ID: {id}' },
    reply: { placeholder: 'Write a reply…', requestId: 'Related request ID (optional)', send: 'Send reply', success: 'Reply sent', closed: 'This ticket is closed and cannot be replied to by the user.' },
    actions: { resolve: 'Mark resolved', close: 'Close ticket', reopenHint: 'Replying to a resolved ticket returns it to the support queue.', resolved: 'Ticket marked resolved', closed: 'Ticket closed' },
    admin: { title: 'Ticket Management', description: 'Handle user issues, assign owners, and track progress.', myTickets: 'View my tickets', detailTitle: 'Handle ticket', claim: 'Claim', save: 'Save ticket settings', updated: 'Ticket updated', noAssignee: 'Unassigned', permissions: 'Ticket access is split into read, reply, and manage capabilities.' },
    settings: { title: 'Sidebar modules', description: 'Control the ticket entry and new-ticket intake in the personal area.', personalArea: 'Personal area', showCenter: 'Show user ticket center', showCenterHint: 'Controls whether users can see and access My Tickets. Admin ticket management always remains available.', acceptNew: 'Accept new tickets', acceptNewHint: 'When disabled, users can still view and reply to existing tickets but cannot create new ones.', save: 'Save ticket settings', saved: 'Ticket settings saved', confirmTitle: 'Hide the user ticket center?', confirmDescription: 'Users will no longer see, access, or reply to tickets. Existing data remains and administrators retain access.', pendingWarning: '{count} non-closed tickets still remain.', confirm: 'Hide ticket center' },
  },
}

