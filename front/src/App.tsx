import { useState } from 'react'
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material'
import Signup from './components/Signup'
import './App.css'

const theme = createTheme({
  palette: {
    mode: 'light',
  },
})

function App() {
  const [currentView, setCurrentView] = useState<'home' | 'signup'>('signup')

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {currentView === 'signup' && <Signup />}
      {currentView === 'home' && (
        <div>
          <h1>Welcome to Real-Time Chat</h1>
          <button onClick={() => setCurrentView('signup')}>
            Go to Signup
          </button>
        </div>
      )}
    </ThemeProvider>
  )
}

export default App
