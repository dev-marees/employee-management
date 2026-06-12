import "../pages/Dashboard/Dashboard.css";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { logout } from "../API/auth";
import { getStats } from "../API/dashboard";
import { getMyProfile } from "../API/employees";


function EmployeeDashboard() {

    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        async function load() {
            try{
                const me = await getMyProfile();
                setProfile(me);
            } catch(err) {
                if (err.status === 401) { logout(); navigate("/"); }
                else setError(err.message);
            } finally {
                setLoading(false);
            }
        }
        load();
    },[])

    if (loading) return <p className="dash-state">Loading...</p>
    if (error) return <p className="dash-state dash-error">{error}</p>

    return (
        <section className="profile-section">
            <div className="panel">
                <h3 className="panel-title">My Profile</h3>
                <div className="profile-grid">
                    <div><span className="profile-label">Employee Code</span><strong>{profile.employee_code}</strong></div>
                    <div><span className="profile-label">Name</span><strong>{profile.first_name} {profile.last_name}</strong></div>
                    <div><span className="profile-label">Email</span><strong>{profile.email}</strong></div>
                    <div><span className="profile-label">Phone</span><strong>{profile.phnoe || "-"}</strong></div>
                    <div><span className="profile-label">Department</span><strong>{profile.department}</strong></div>
                    <div><span classN ame="profile-label">Designation</span><strong>{profile.designation || "-"}</strong></div>
                    <div><span className="profile-label">Joining Date</span><strong>{profile.joining_date}</strong></div>
                    <div>
                        <span className="profile-label">Status</span>
                        <span className={`badge badge-${profile.status}`}>{profile.status}</span>
                    </div>
                </div>
            </div>
        </section>
    );
}

export default EmployeeDashboard;
