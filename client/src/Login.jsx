import { useState } from "react"
import { useNavigate } from "react-router-dom";




function Login(){
    const [username,setUsername] = useState("");
    const [password,setPassword] = useState("");
    const navigate = useNavigate();

    function navToRegister(){
        navigate("/Register")
    }

    function loginUser(){

    }


    return(
        <>
            <div id="loginContainer">
                <h1>Login</h1>
                <div>
                    <label id="usernameLoginLabel" className="label">Username:</label>
                    <input value={username} onChange={(e) => setUsername(e.target.value)}/>
                </div>
                <div>
                    <label id="passwordLoginLabel" className="label">Password:</label>
                    <input value={password} onChange={(e) =>setPassword(e.target.value)}/>
                </div>
                <div>
                    <button onClick={navToRegister}>Back To Register</button>
                    <button onClick={loginUser}>Login</button>
                </div>

            </div>
        </>
    )
}


export default Login