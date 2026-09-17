import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ViewSnippet } from './ViewSnippet'
import * as api from '../api/snippets'

vi.mock('../api/snippets', () => ({
  getSnippet: vi.fn(),
  deleteSnippet: vi.fn(),
}))

const mockGetSnippet = vi.fn<typeof api.getSnippet>(api.getSnippet)
const mockDeleteSnippet = vi.fn<typeof api.deleteSnippet>(api.deleteSnippet)
vi.mocked(api.getSnippet).mockImplementation(mockGetSnippet)
vi.mocked(api.deleteSnippet).mockImplementation(mockDeleteSnippet)

function renderWithRouter(initialPath = '/s/test-slug') {
  return render(
    <MemoryRouter initialEntries={[initialPath]}>
      <Routes>
        <Route path="/s/:slug" element={<ViewSnippet />} />
        <Route path="/" element={<span data-testid="home">Home</span>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('ViewSnippet', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockGetSnippet.mockReset()
    mockDeleteSnippet.mockReset()
  })

  it('shows loading state', () => {
    mockGetSnippet.mockImplementation(() => new Promise(() => {}))
    renderWithRouter()
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })

  it('displays snippet content on success', async () => {
    mockGetSnippet.mockResolvedValueOnce({
      id: 1,
      slug: 'test-slug',
      content: 'console.log("hi")',
      created_at: new Date().toISOString(),
    })
    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByText('console.log("hi")')).toBeInTheDocument()
    })
    expect(screen.getByText(/created:/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
  })

  it('shows 404 when snippet not found', async () => {
    const err = Object.assign(new Error('Not found'), { status: 404 })
    mockGetSnippet.mockRejectedValueOnce(err)
    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByText('Snippet not found')).toBeInTheDocument()
    })
  })

  it('shows error on API failure', async () => {
    mockGetSnippet.mockRejectedValueOnce(new Error('Server error'))
    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByText('Server error')).toBeInTheDocument()
    })
  })

  it('calls DELETE endpoint on delete click', async () => {
    mockGetSnippet.mockResolvedValueOnce({
      id: 1,
      slug: 'test-slug',
      content: 'code',
      created_at: new Date().toISOString(),
    })
    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /delete/i }))
    expect(mockDeleteSnippet).toHaveBeenCalledWith('test-slug')
  })

  it('redirects to / after successful delete', async () => {
    mockGetSnippet.mockResolvedValueOnce({
      id: 1,
      slug: 'test-slug',
      content: 'code',
      created_at: new Date().toISOString(),
    })
    mockDeleteSnippet.mockResolvedValueOnce(undefined)
    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /delete/i }))
    await waitFor(() => {
      expect(screen.getByTestId('home')).toBeInTheDocument()
    })
  })
})
