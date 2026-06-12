import "../Dashboard/Dashboard.css";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { logout } from "../../API/auth";
import { getStats } from "../../API/dashboard";
import { listEmployees, getMyProfile } from "../../API/employees";
import AdminDashboard from "../../components/AdminDashboard";
import EmployeeDashboard from "../../components/EmployeeDashboard";


function Dashboard() {

    const navigate = useNavigate();
    const user = JSON.parse(localStorage.getItem("user")) || null;
    const isManager = user?.role === "Admin" || user?.role === "HR"

    function handleLogout() {
        logout();
        navigate("/");
    }

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
                    handleLogout();
                }}>⤶ Logout</button>
            </aside>

            {/* ===== Main area ===== */}
            <div className="dashboard-main">
                <header className="topbar">
                <div>
                    <h1>Dashboard</h1>
                    <p>Welcome back 👋</p>
                </div>
                <div className="topbar-user">
                    <div className="user-meta">
                    <strong>{user?.name}</strong>
                    <span>{user?.role}</span>
                    </div>
                    <div className="user-avatar">{user?.name?.[0]?.toUpperCase()}</div>
                </div>
                </header>
                {isManager ? <AdminDashboard /> : <EmployeeDashboard />}
            </div>
        </div>
    );
}

export default Dashboard;
