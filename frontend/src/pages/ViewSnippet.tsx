import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getSnippet, deleteSnippet } from '../api/snippets'

type State = 'loading' | 'success' | 'not-found' | 'error'

export function ViewSnippet() {
  const { slug } = useParams<{ slug: string }>()!
  const navigate = useNavigate()
  const [state, setState] = useState<State>('loading')
  const [snippet, setSnippet] = useState<{ content: string; created_at: string } | null>(null)
  const [errorMsg, setErrorMsg] = useState('')

  useEffect(() => {
    getSnippet(slug)
      .then((data) => {
        setSnippet({ content: data.content, created_at: data.created_at })
        setState('success')
      })
      .catch((err) => {
        if ((err as Error & { status?: number }).status === 404) {
          setState('not-found')
        } else {
          setErrorMsg(err instanceof Error ? err.message : 'Failed to fetch snippet')
          setState('error')
        }
      })
  }, [slug])

  async function handleDelete() {
    try {
      await deleteSnippet(slug)
      navigate('/')
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to delete snippet')
      setState('error')
    }
  }

  if (state === 'loading') {
    return (
      <div className="container">
        <h1>Pastebin Lite</h1>
        <p>Loading...</p>
      </div>
    )
  }

  if (state === 'not-found') {
    return (
      <div className="container">
        <h1>Pastebin Lite</h1>
        <p>Snippet not found</p>
        <a href="/">← Back to home</a>
      </div>
    )
  }

  if (state === 'error') {
    return (
      <div className="container">
        <h1>Pastebin Lite</h1>
        <p className="error">{errorMsg}</p>
        <a href="/">← Back to home</a>
      </div>
    )
  }

  return (
    <div className="container">
      <h1>Pastebin Lite</h1>
      <pre>{snippet?.content}</pre>
      <p>Created: {new Date(snippet!.created_at).toLocaleString()}</p>
      <button onClick={handleDelete}>Delete</button>
    </div>
  )
}
