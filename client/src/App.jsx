import { useState } from 'react'
import './App.css'
import {BrowserRouter,Routes,Route} from "react-router-dom"
import Login from './Login'
import Register from './Register'

function App() {

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<Login/>}></Route>
          <Route path='/Register' element={<Register/>}></Route>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
