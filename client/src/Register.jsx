import { useNavigate } from "react-router-dom";



function Register(){
    const [username,setUsername] = useState("");
    const [password,setPassword] = useState("");

    const navigate = useNavigate();

    function navToLogin(){
        navigate("/")
    }

    function registerUser(){

    }

    return(
        <>
            <div id="registerContainer">
                <h1>Register</h1>
                <div>
                    <label id="usernameRegisterLabel" className="label">Username:
                        <input value={username} onChange={(e) => setUsername(e.target.value)}/>
                    </label>
                </div>
                <div>
                    <label id="passwordRegisterLabel" className="label">Password:<
                        input value={password} onChange={(e) =>setPassword(e.target.value)}/>
                    </label>
                </div>
                <div>
                    <button onClick={navToLogin}>Back To Login</button>
                    <button onClick={registerUser}></button>
                </div>

            </div>
        </>
    )
}


export default Register