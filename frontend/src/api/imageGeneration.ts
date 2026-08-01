import { buildGatewayUrl } from './client'

export interface ImageGenerationModel {
  id: string
  object?: string
  owned_by?: string
}

export interface ImageGenerationModelsResponse {
  object: string
  data: ImageGenerationModel[]
}

export interface ImageGenerationRequest {
  model: string
  prompt: string
  n?: number
  size?: string
  response_format?: 'url' | 'b64_json'
  output_format?: 'png' | 'jpeg' | 'webp'
}

export interface GeneratedImageOutput {
  url?: string
  b64_json?: string
  revised_prompt?: string
}

export interface ImageGenerationResponse {
  created?: number
  data: GeneratedImageOutput[]
}

export type ImageTaskStatus = 'processing' | 'completed' | 'failed' | string

export interface ImageGenerationTask {
  id: string
  task_id: string
  object: string
  status: ImageTaskStatus
  http_status?: number
  image_url?: string
  result?: ImageGenerationResponse
  error?: {
    type?: string
    code?: string
    message?: string
  }
  created_at: number
  completed_at?: number
  expires_at: number
}

export type ImageGenerationSubmission =
  | { mode: 'async'; task: ImageGenerationTask }
  | { mode: 'sync'; result: ImageGenerationResponse }

export interface ImageGenerationSubmissionOptions {
  signal?: AbortSignal
  preferAsync?: boolean
}

export class ImageGenerationAPIError extends Error {
  status: number
  code: string
  requestId: string

  constructor(message: string, status = 0, code = '', requestId = '') {
    super(message)
    this.name = 'ImageGenerationAPIError'
    this.status = status
    this.code = code
    this.requestId = requestId
  }
}

function authHeaders(apiKey: string, extra?: HeadersInit): HeadersInit {
  return {
    Authorization: `Bearer ${apiKey}`,
    ...extra,
  }
}

async function parseImageGenerationError(response: Response): Promise<ImageGenerationAPIError> {
  let message = response.statusText || `HTTP ${response.status}`
  let code = String(response.status)
  try {
    const body = await response.json()
    const error = body?.error ?? body
    message = error?.message || body?.message || message
    code = error?.code || error?.type || body?.code || code
  } catch {
    // Keep the HTTP fallback when the upstream returned a non-JSON error page.
  }
  return new ImageGenerationAPIError(
    message,
    response.status,
    String(code || ''),
    response.headers.get('X-Request-Id') || '',
  )
}

function isAsyncFeatureUnavailable(error: ImageGenerationAPIError): boolean {
  return error.status === 404 && /async image tasks are not enabled/i.test(error.message)
}

export function isImageGenerationModel(model: string): boolean {
  const normalized = String(model || '').trim().toLowerCase()
  return normalized.startsWith('gpt-image-') ||
    normalized === 'grok-imagine' ||
    normalized.startsWith('grok-imagine-image')
}

export async function listImageGenerationModels(
  apiKey: string,
  signal?: AbortSignal,
): Promise<ImageGenerationModelsResponse> {
  const response = await fetch(buildGatewayUrl('/v1/models'), {
    headers: authHeaders(apiKey),
    signal,
  })
  if (!response.ok) throw await parseImageGenerationError(response)
  const body = await response.json()
  return {
    object: String(body?.object || 'list'),
    data: Array.isArray(body?.data) ? body.data : [],
  }
}

async function submitSynchronousImageGeneration(
  apiKey: string,
  payload: ImageGenerationRequest,
  signal?: AbortSignal,
): Promise<ImageGenerationResponse> {
  const response = await fetch(buildGatewayUrl('/v1/images/generations'), {
    method: 'POST',
    headers: authHeaders(apiKey, { 'Content-Type': 'application/json' }),
    body: JSON.stringify(payload),
    signal,
  })
  if (!response.ok) throw await parseImageGenerationError(response)
  return response.json()
}

export async function submitImageGeneration(
  apiKey: string,
  payload: ImageGenerationRequest,
  options: ImageGenerationSubmissionOptions = {},
): Promise<ImageGenerationSubmission> {
  if (options.preferAsync === false) {
    return {
      mode: 'sync',
      result: await submitSynchronousImageGeneration(apiKey, payload, options.signal),
    }
  }

  const response = await fetch(buildGatewayUrl('/v1/images/generations/async'), {
    method: 'POST',
    headers: authHeaders(apiKey, { 'Content-Type': 'application/json' }),
    body: JSON.stringify(payload),
    signal: options.signal,
  })

  if (response.ok) {
    return { mode: 'async', task: await response.json() }
  }

  const error = await parseImageGenerationError(response)
  if (!isAsyncFeatureUnavailable(error)) throw error

  return {
    mode: 'sync',
    result: await submitSynchronousImageGeneration(apiKey, payload, options.signal),
  }
}

export async function getImageGenerationTask(
  apiKey: string,
  taskId: string,
  signal?: AbortSignal,
): Promise<ImageGenerationTask> {
  const response = await fetch(
    buildGatewayUrl(`/v1/images/tasks/${encodeURIComponent(taskId)}`),
    { headers: authHeaders(apiKey), signal },
  )
  if (!response.ok) throw await parseImageGenerationError(response)
  return response.json()
}

export function imageOutputSource(output: GeneratedImageOutput, outputFormat = 'png'): string {
  if (output.url) return output.url
  if (!output.b64_json) return ''
  const mimeType = outputFormat === 'jpeg' ? 'image/jpeg' : outputFormat === 'webp' ? 'image/webp' : 'image/png'
  return `data:${mimeType};base64,${output.b64_json}`
}

export function imageTaskResult(task: ImageGenerationTask): ImageGenerationResponse | null {
  if (task.result && Array.isArray(task.result.data)) return task.result
  if (task.image_url) return { data: [{ url: task.image_url }] }
  return null
}
