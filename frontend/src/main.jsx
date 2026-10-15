import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@fontsource/pridi/600.css'
import '@fontsource-variable/anuphan'
import './index.css'
import App from './App.jsx'
{/*เปิดโหมด StrictMode เพื่อให้ React ตรวจสอบการใช้งานที่ไม่ถูกต้องใน component tree*/}
createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
