import { useState } from "react"
import { useNavigate } from "react-router-dom";




function Appointment(){
    const [date,setDate] = useState();
    const [time,setTime] = useState();
    const [reason,setReason] = useState();
    const navigate = useNavigate();

    const timeSlots = ["08:00","09:00","10:00","11:00","12:00","13:00","14:00","15:00","16:00"]

    const bookAppointment = async () =>{
        const token = localStorage.getItem("token");
        if(!token){
            alert("Please Login First");
            navigate("/");
            return
        }
        const response = await fetch("http://localhost:8080/appointment",{
            method:"POST",
            headers:{
                "Content-Type":"application/json",
                "Authorization":`Bearer ${token}`
            },
            body:JSON.stringify({date,time,reason})

        })
        const data = await response.json();
        alert(data.message)
    }



    return(
        <>
            <div id="appointmentContainer">
                <h1 id="title">Book An Appointment</h1>
                <div id="form-group">
                    <input type="date" value={date} min={new Date().toISOString().split("T")[0]} max={new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split("T")[0]} onChange={e =>setDate(e.target.value)}></input>
                </div>
                <div id="form-group">
                <select value={time} onChange={e => setTime(e.target.value)}>
                    <option value="">Select Time</option>
                    {timeSlots.map(slot => (
                        <option key={slot} value={slot}>{slot}</option>
                    ))}
                </select>
                </div>
                <div id="form-group">
                <input type="text" placeholder="Reason" value={reason} onChange={e => setReason(e.target.value)} />
                <button onClick={bookAppointment}>Book</button>
                </div>
            </div>
        </>
    )
}

export default Appointment