export interface Snippet {
  id: number
  slug: string
  content: string
  created_at: string
}

export interface CreateSnippetInput {
  content: string
}

export async function createSnippet(input: CreateSnippetInput): Promise<Snippet> {
  const res = await fetch('/api/snippets', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? 'Failed to create snippet')
  }
  return res.json()
}

export async function getSnippet(slug: string): Promise<Snippet> {
  const res = await fetch(`/api/snippets/${slug}`)
  if (res.status === 404) {
    const err = await res.json().catch(() => ({ error: 'Not found' }))
    throw Object.assign(new Error(err.error ?? 'Snippet not found'), { status: 404 })
  }
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? 'Failed to fetch snippet')
  }
  return res.json()
}

export async function deleteSnippet(slug: string): Promise<void> {
  const res = await fetch(`/api/snippets/${slug}`, { method: 'DELETE' })
  if (!res.ok && res.status !== 204) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? 'Failed to delete snippet')
  }
}
