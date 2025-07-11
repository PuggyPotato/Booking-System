import { useState } from 'react'
import './App.css'
import {BrowserRouter,Routes,Route} from "react-router-dom"
import Login from './Login'
import Register from './Register'
import Home from './Home'
import Appointment from './Appointment'
import Admin from './Admin'
import AdminLogin from './AdminLogin'

function App() {

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<Login/>}></Route>
          <Route path='/Register' element={<Register/>}></Route>
          <Route path="/Home" element={<Home/>}></Route>
          <Route path="/Appointment" element={<Appointment/>}></Route>
          <Route path="/admin/appointments" element={<Admin/>}></Route>
          <Route path="/admin/login" element={<AdminLogin/>}></Route>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
