import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Arena from './Arena'

function App() {
  const [count, setCount] = useState(0)

  return (
    <Arena/>
  )
}

export default App
