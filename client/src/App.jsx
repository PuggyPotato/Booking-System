import { useState } from 'react'
import './App.css'
import {BrowserRouter,Routes,Route} from "react-router-dom"
import Login from './Login'
import Register from './Register'
import Home from './Home'

function App() {

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<Login/>}></Route>
          <Route path='/Register' element={<Register/>}></Route>
          <Route path="/Home" element={<Home/>}></Route>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
