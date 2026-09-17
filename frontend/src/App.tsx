import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { CreateSnippet } from './pages/CreateSnippet'
import { ViewSnippet } from './pages/ViewSnippet'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<CreateSnippet />} />
        <Route path="/s/:slug" element={<ViewSnippet />} />
      </Routes>
    </BrowserRouter>
  )
}
