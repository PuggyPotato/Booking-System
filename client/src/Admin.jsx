import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

function Admin() {
  const [appointments, setAppointments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [unauthorized, setUnauthorized] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("adminToken");
    if (!token) {
      setUnauthorized(true);
      return;
    }

    fetch("http://localhost:8080/adminAppointments", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text);
        }
        return res.json();
      })
      .then((data) => {
        setAppointments(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error:", err.message);
        setUnauthorized(true);
        setLoading(false);
      });
  }, []);

  if (unauthorized) {
    return (
      <div style={{ textAlign: "center", marginTop: "2rem" }}>
        <h1>🚫 Unauthorized</h1>
        <p>You must be an admin to view this page.</p>
        <button onClick={() => navigate("/")}>Go to Home</button>
      </div>
    );
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h1>📋 All Appointments</h1>
      {loading ? (
        <p>Loading...</p>
      ) : appointments.length === 0 ? (
        <p>No appointments found.</p>
      ) : (
        <table
          style={{
            borderCollapse: "collapse",
            width: "100%",
            marginTop: "1rem",
          }}
        >
          <thead>
            <tr>
              <th style={thStyle}>👤 Username</th>
              <th style={thStyle}>📅 Date</th>
              <th style={thStyle}>⏰ Time</th>
              <th style={thStyle}>📝 Reason</th>
            </tr>
          </thead>
          <tbody>
            {appointments.map((a, i) => (
              <tr key={i}>
                <td style={tdStyle}>{a.username}</td>
                <td style={tdStyle}>{a.date}</td>
                <td style={tdStyle}>{a.time}</td>
                <td style={tdStyle}>{a.reason}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

const thStyle = {
  border: "1px solid #ccc",
  padding: "10px",
  backgroundColor: "#f5f5f5",
  fontWeight: "bold",
};

const tdStyle = {
  border: "1px solid #ccc",
  padding: "10px",
};

export default Admin;
