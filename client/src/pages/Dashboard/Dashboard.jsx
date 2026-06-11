import "../Dashboard/Dashboard.css";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { logout } from "../../API/auth";
import { getStats } from "../../API/dashboard";
import { listEmployees, getMyProfile } from "../../API/employees";


function Dashboard() {

    const navigate = useNavigate();
    const user = JSON.parse(localStorage.getItem("user")) || null;
    const isManager = user?.role === "Admin" || user?.role === "HR"

    const [stats, setStats] = useState(null);
    const [employees, setEmployees] = useState([]);
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        async function load() {
            try{
                if(isManager) {
                    const s = await getStats();
                    const list = await listEmployees({ limit: 10 });
                    setStats(s);
                    setEmployees(list.items);
                } else {
                    const me = await getMyProfile();
                    setProfile(me);
                }
            } catch(err) {
                if (err.status === 401) { logout(); navigate("/"); }
                else setError(err.message);
            } finally {
                setLoading(false);
            }
        }
        load();
    },[])

    return (
        <div className="dashboard">
            {/* ===== Sidebar ===== */}
            <aside className="sidebar">
                <div className="sidebar-brand">
                    <h2>EMS</h2>
                    <span>Employee Management</span>
                </div>

                <nav className="sidebar-nav">
                    <a className="nav-item active" href="#">📊 Dashboard</a>
                    <a className="nav-item" href="#">👥 Employees</a>
                    <a className="nav-item" href="#">🏢 Departments</a>
                    <a className="nav-item" href="#">⚙️ Settings</a>
                </nav>

                {/* TODO: wire onClick to your logout() from API/auth.js */}
                <button className="logout-btn" onClick={() => {
                    logout(); navigate("/");
                }}>⤶ Logout</button>
            </aside>

            {/* ===== Main area ===== */}
            <div className="dashboard-main">
            {/* ----- Topbar ----- */}
                <header className="topbar">
                    <div>
                        <h1>Dashboard</h1>
                        <p>Welcome back 👋</p>
                    </div>

                    <div className="topbar-user">
                        {/* TODO: replace with the logged-in user's name/role from localStorage */}
                        <div className="user-meta">
                            <strong>{user?.name}</strong>
                            <span>{user?.role}</span>
                        </div>
                        <div className="user-avatar">{user?.name?.[0]?.toUpperCase()}</div>
                    </div>
                </header>

            {/* ----- Stat cards ----- */}
                <section className="stats-grid">
                    <div className="stat-card">
                        <div className="stat-icon blue">👥</div>
                        <div className="stat-info">
                            <span className="stat-label">Total Employees</span>
                            {/* TODO: {stats.total_employees} */}
                            <strong className="stat-value">0</strong>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon green">✅</div>
                        <div className="stat-info">
                            <span className="stat-label">Active</span>
                            {/* TODO: {stats.active_employees} */}
                            <strong className="stat-value">0</strong>
                        </div>
                    </div>

                    <div className="stat-card">
                        <div className="stat-icon red">⛔</div>
                        <div className="stat-info">
                                <span className="stat-label">Inactive</span>
                                {/* TODO: {stats.inactive_employees} */}
                                <strong className="stat-value">0</strong>
                        </div>
                    </div>

                    <div className="stat-card">
                        <div className="stat-icon purple">🏢</div>
                        <div className="stat-info">
                            <span className="stat-label">Departments</span>
                            {/* TODO: {stats.department_wise_count.length} */}
                            <strong className="stat-value">0</strong>
                        </div>
                    </div>
                </section>

                {/* ----- Bottom panels ----- */}
                <section className="dashboard-bottom">
                    {/* Department-wise breakdown */}
                    <div className="panel dept-panel">
                        <h3 className="panel-title">
                            Department-wise Count
                        </h3>

                        {/* TODO: map stats.department_wise_count -> one .dept-row per item */}
                        <div className="dept-row">
                            <span className="dept-name">Engineering</span>
                            <div className="dept-bar">
                            <div className="dept-bar-fill" style={{ width: "70%" }} />
                            </div>
                            <span className="dept-count">0</span>
                        </div>
                        <div className="dept-row">
                            <span className="dept-name">Sales</span>
                            <div className="dept-bar">
                            <div className="dept-bar-fill" style={{ width: "45%" }} />
                            </div>
                            <span className="dept-count">0</span>
                        </div>

                        <div className="dept-row">
                            <span className="dept-name">HR</span>
                            <div className="dept-bar">
                            <div className="dept-bar-fill" style={{ width: "25%" }} />
                            </div>
                            <span className="dept-count">0</span>
                        </div>
                    </div>

                    {/* Employees table */}
                    <div className="panel table-panel">
                        <div className="panel-header">
                            <h3 className="panel-title">Employees</h3>
                            {/* TODO: wire search input + filters to GET /employees query params */}
                            <input
                            className="table-search"
                            type="text"
                            placeholder="Search employees..."
                            />
                        </div>

                        <div className="table-wrap">
                            <table className="emp-table">
                                <thead>
                                    <tr>
                                        <th>Code</th>
                                        <th>Name</th>
                                        <th>Email</th>
                                        <th>Department</th>
                                        <th>Designation</th>
                                        <th>Salary</th>
                                        <th>Status</th>
                                        <th>Joining Date</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {/* TODO: map employee list -> one <tr> per employee */}
                                    <tr>
                                        <td>EMP-0001</td>
                                        <td>John Smith</td>
                                        <td>john.smith@acme.com</td>
                                        <td>Engineering</td>
                                        <td>Senior Engineer</td>
                                        <td>$95,000</td>
                                        <td><span className="badge badge-active">active</span></td>
                                        <td>2023-04-01</td>
                                    </tr>
                                    <tr>
                                        <td>EMP-0002</td>
                                        <td>Jane Doe</td>
                                        <td>jane.doe@acme.com</td>
                                        <td>Sales</td>
                                        <td>Account Manager</td>
                                        <td>$72,000</td>
                                        <td><span className="badge badge-inactive">inactive</span></td>
                                        <td>2022-11-15</td>
                                    </tr>

                                    {/* Empty-state row (show when there are no employees) */}
                                    {/*
                                    <tr>
                                    <td colSpan={8} className="empty-row">No employees found</td>
                                    </tr>
                                    */}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
}

export default Dashboard;
