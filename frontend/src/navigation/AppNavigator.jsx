import { Routes, Route, HashRouter } from 'react-router-dom'

import { useState } from "react";

import Home from '../pages/home/home'
import Commands from '../pages/DiskCreen/DiskCreen'

export default function AppNavigator() {
  const [ip, setIP] = useState("localhost") 

  const handleChage = (e) => {
    console.log(e.target.value)
    setIP(e.target.value)
  }

  return (
    <HashRouter>
      IP: <input type="text" onChange={handleChage}/> -- {ip}
      <Routes>
         
      <Route path="/" element={<Home/>} />
      <Route path="/DiskCreen" element={<Commands ip={ip}/>} />
      </Routes>
    </HashRouter>
  )
}