import "../pages/Dashboard/Dashboard.css";
import { useState, useEffect } from "react";
import { getStats } from "../API/dashboard";
import { listEmployees } from "../API/employees";
import EmployeeFormModal from "./EmployeeFormModal";


function AdminDashboard() {

    const [stats, setStats] = useState(null);
    const [employees, setEmployees] = useState([]);
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [showForm, setShowForm] = useState(false);

    async function loadData(){
        setLoading(true);
        try {
            const s = await getStats();
            const list = await listEmployees({ limit: 10 });
            setStats(s);
            setEmployees(list.data);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => { loadData(); },[])

    if (loading) return <p className="dash-state">Loading...</p>
    if (error) return <p className="dash-state dash-error">{error}</p>

    const maxDept = stats?.department_wise_count?.[0]?.count || 1;

    return (
        <>
            <section className="stats-grid">
                <div className="stat-card">
                    <div className="stat-icon blue">👥</div>
                    <div className="stat-info">
                        <span className="stat-label">Total Employees</span>
                        <strong className="stat-value">{stats.total_employees}</strong>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon green">✅</div>
                    <div className="stat-info">
                        <span className="stat-label">Active</span>
                        {/* TODO: {stats.active_employees} */}
                        <strong className="stat-value">{stats.active_employees}</strong>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-icon red">⛔</div>
                    <div className="stat-info">
                            <span className="stat-label">Inactive</span>
                            {/* TODO: {stats.inactive_employees} */}
                            <strong className="stat-value">{stats.inactive_employees}</strong>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-icon purple">🏢</div>
                    <div className="stat-info">
                        <span className="stat-label">Departments</span>
                        <strong className="stat-value">{stats.department_wise_count.length}</strong>
                    </div>
                </div>
            </section>
            <section className="dashboard-bottom">
                <div className="panel dept-panel">
                    <h3 className="panel-title">
                        Department-wise Count
                    </h3>
                    {stats.department_wise_count.map((d) => (
                        <div className="dept-row" key={d.department}>
                            <span className="dept-name">{d.department}</span>
                            <div className="dept-bar">
                            <div className="dept-bar-fill" style={{ width: `${(d.count / maxDept) * 100}%` }} />
                            </div>
                            <span className="dept-count">{d.count}</span>
                        </div>
                    ))}
                </div>
                <div className="panel table-panel">
                    <div className="panel-header">
                        <h3 className="panel-title">Employees</h3>
                        <input className="table-search" type="text" placeholder="Search employees..."
                        />
                        <button className="add-btn" onClick={() => setShowForm(true)}>+ Add Employee</button>
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
                                {employees.length === 0 ? (<tr><td colSpan={8} className="empty-row">No employees found </td></tr>): (employees.map((e) => (
                                    <tr key={e.id}>
                                        <td>{e.employee_code}</td>
                                        <td>{e.first_name} {e.last_name}</td>
                                        <td>{e.email}</td>
                                        <td>{e.department}</td>
                                        <td>{e.designation}</td>
                                        <td>{e.salary.toLocaleString()}</td>
                                        <td><span className={`badge badge-${e.status}`}>{e.status}</span></td>
                                        <td>{e.joining_date}</td>
                                    </tr>
                                )))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </section>
            {showForm && (
                <EmployeeFormModal
                    onClose={() => setShowForm(false)}
                    onCreated={loadData}
                />
            )}
        </>
    );
}

export default AdminDashboard;
