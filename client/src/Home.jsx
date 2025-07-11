import { useEffect, useState } from "react";
import {jwtDecode} from "jwt-decode";
import { useNavigate } from "react-router-dom";





function Home(){
    const [username,setUsername] = useState("Guest");
    const [loggedIn,setLoggedIn] = useState(false);


    const navigate = useNavigate();
    useEffect(() =>{
        const token = localStorage.getItem("token");
        if(token){
            try{
                setLoggedIn(true);
                const decoded = jwtDecode(token);
                console.log(decoded);
                setUsername(decoded.username);
            }
            catch(err){
                console.log("Error:",err);
            }
        }
        else{
            setUsername("Guest");
        }
    },[])

    function logOut(){
        localStorage.removeItem("token");
        setLoggedIn(false);
        navigate("/");
    }




    return(
        <>
            <div id="container">
                <div id="navBar">
                        <button id="home" onClick={() =>navigate("/Home")}>Home</button>
                        <button id="logOut" onClick={() =>logOut()}>Log Out</button>
                </div>
                <div id="homeContainer">
                    <h1>Welcome {username}</h1>
                    <button id="appointment" onClick={() =>navigate("/Appointment")}>Make An Appointment</button>
                </div>
            </div>
        </>
    )
}

export default Home