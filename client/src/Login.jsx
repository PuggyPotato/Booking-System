import { useState } from "react"
import { useNavigate } from "react-router-dom";




function Login(){
    const [username,setUsername] = useState("");
    const [password,setPassword] = useState("");
    const navigate = useNavigate();

    function navToRegister(){
        navigate("/Register")
    }

    const loginUser = async () =>{
        try{
            const response = await fetch("http://localhost:8080/login",{
                method:"POST",
                headers:{
                    "Content-Type":"application/json"
                },
                body:JSON.stringify({username,password})
            });
            if (!response.ok){
                const errorText = await response.text();
                console.log("Backend Error:",errorText);
                return
            }
            const result = await response.json();
            console.log(result);
            navigate("/Home");
        }
        catch(err){
            console.log("Error Sending Data:",err)
        }
    }


    return(
        <>
            <div id="loginContainer">
                <h1 id="loginHeader">Login</h1>
                <div id="credentials">
                    <div>
                        <label id="usernameLoginLabel" className="label">Username:</label>
                        <input value={username} onChange={(e) => setUsername(e.target.value)} id="usernameInput"/>
                    </div>
                    <div>
                        <label id="passwordLoginLabel" className="label">Password:</label>
                        <input value={password} onChange={(e) =>setPassword(e.target.value)} id="passwordInput"/>
                    </div>
                </div>
                <div id="authNav">
                    <button onClick={navToRegister} id="navRegister">Back To Register</button>
                    <button onClick={loginUser} id="loginButton">Login</button>
                </div>

            </div>
        </>
    )
}


export default Login