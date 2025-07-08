import { useNavigate } from "react-router-dom";
import { useState } from "react";


function Register(){
    const [username,setUsername] = useState("");
    const [password,setPassword] = useState("");

    const navigate = useNavigate();

    function navToLogin(){
        navigate("/")
    }

    const registerUser = async () =>{
        try{
            const response = await fetch("http://localhost:8080/register",{
                method:"POST",
                headers:{
                    "Content-Type":"application/json",
                },
                body: JSON.stringify({username,password})
            });
            if(!response.ok){
                const errorText = response.text();
                console.log("Backend Error:",errorText);
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
            <div id="registerContainer">
                <h1 id="registerHeader">Register</h1>
                <div id="credentials">
                    <div>
                        <label id="usernameRegisterLabel" className="label">Username:
                            <input value={username} onChange={(e) => setUsername(e.target.value)} id="usernameInput"/>
                        </label>
                    </div>
                    <div>
                        <label id="passwordRegisterLabel" className="label">Password:
                            <input value={password} onChange={(e) =>setPassword(e.target.value)} id="passwordInput"/>
                        </label>
                    </div>
                </div>
                <div id="authNav">
                    <button onClick={navToLogin} id="navLogin">Back To Login</button>
                    <button onClick={registerUser} id="registerButton">Register</button>
                </div>

            </div>
        </>
    )
}


export default Register