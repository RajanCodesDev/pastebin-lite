import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { CreateSnippet } from './CreateSnippet'
import * as api from '../api/snippets'

const mockCreateSnippet = vi.fn<typeof api.createSnippet>(api.createSnippet)

describe('CreateSnippet', () => {
  beforeEach(() => {
    vi.resetAllMocks()
    vi.stubEnv('VITE_API_URL', 'http://localhost:8080')
  })

  it('renders textarea and create button', () => {
    render(<CreateSnippet />)
    expect(screen.getByRole('textbox')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create snippet/i })).toBeInTheDocument()
  })

  it('shows validation error on empty submission', async () => {
    render(<CreateSnippet />)
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.click(btn)
    await waitFor(() => {
      expect(screen.getByText(/content cannot be empty/i)).toBeInTheDocument()
    })
  })

  it('shows loading state while submitting', async () => {
    mockCreateSnippet.mockImplementationOnce(() => new Promise<typeof api.Snippet>((resolve) =>
      setTimeout(() => resolve({ id: 1, slug: 'abc123', content: 'hello', created_at: new Date().toISOString() }), 50)))

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'hello world' } })
    fireEvent.click(btn)

    expect(btn).toBeDisabled()
    expect(screen.getByText(/creating\.\.\./i)).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText(/snippet created/i)).toBeInTheDocument()
    })
    expect(mockCreateSnippet).toHaveBeenCalledWith({ content: 'hello world' })
  })

  it('displays generated URL on success', async () => {
    const snippet = { id: 1, slug: 'a8f31c', content: 'hello world', created_at: new Date().toISOString() }
    mockCreateSnippet.mockResolvedValueOnce(snippet)

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'hello world' } })
    fireEvent.click(btn)

    await waitFor(() => {
      const link = screen.getByRole('link', { name: /http:\/\/localhost\/s\/a8f31c/i })
      expect(link).toBeInTheDocument()
      expect(link).toHaveAttribute('href', 'http://localhost/s/a8f31c')
    })
  })

  it('displays error on API failure', async () => {
    mockCreateSnippet.mockRejectedValueOnce(new Error('Network error'))

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'hello' } })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeInTheDocument()
    })
  })

  it('clears textarea after successful creation', async () => {
    mockCreateSnippet.mockResolvedValueOnce({ id: 1, slug: 'xyz', content: 'test', created_at: new Date().toISOString() })

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    fireEvent.change(textarea, { target: { value: 'my code' } })
    fireEvent.click(screen.getByRole('button'))

    await waitFor(() => {
      expect(textarea).toHaveValue('')
    })
  })
})
