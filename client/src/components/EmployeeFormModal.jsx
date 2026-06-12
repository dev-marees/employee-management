import { useState } from "react";
import { createEmployee } from "../API/employees";

function EmployeeFormModal({ onClose, onCreated }) {
    const [form, setForm] = useState({
        employee_code:"", first_name: "", last_name: "", email: "",
        department: "", designation: "", phone: "", salary: 0, joining_date: "",
        status: "active"
    });
    const [errors, setErrors] = useState({});
    const [error, setError] = useState("");
    const [saving, setSaving] = useState(false);

    const update = (e) => setForm({ ...form, [e.target.name]: e.target.value });

    async function handleSubmit(e) {
        e.preventDefault();
        setError(""); setErrors({}); setSaving(true);
        try {
            await createEmployee({...form, salary:Number(form.salary)});
            onCreated();
            onClose()
        } catch (err) {
            if (err.fields) setErrors(err.fields);
            else setError(err.message);
        } finally {
            setSaving(false);
        }
    }

    return (
        <>
            <div className="modal-overlay" onClick={onClose}>
                <div className="modal" onClick={(e) => e.stopPropagation()}>
                    <h3>Add Employee</h3>
                    {error && <p className="auth-error">{error}</p>}

                    <form onSubmit={handleSubmit} className="emp-form">
                    <input name="employee_code" placeholder="Employee Code" value={form.employee_code} onChange={update} />
                    {errors.employee_code && <small className="field-err">{errors.employee_code}</small>}

                    <input name="first_name" placeholder="First Name" value={form.first_name} onChange={update} />
                    <input name="last_name"  placeholder="Last Name"  value={form.last_name}  onChange={update} />
                    <input name="email" type="email" placeholder="Email" value={form.email} onChange={update} />
                    <input name="phone" placeholder="Phone (optional)" value={form.phone} onChange={update} />
                    <input name="department" placeholder="Department" value={form.department} onChange={update} />
                    <input name="designation" placeholder="Designation (optional)" value={form.designation} onChange={update} />
                    <input name="salary" type="number" placeholder="Salary" value={form.salary} onChange={update} />
                    <input name="joining_date" type="date" value={form.joining_date} onChange={update} />
                    <select name="status" value={form.status} onChange={update}>
                        <option value="active">Active</option>
                        <option value="inactive">Inactive</option>
                    </select>

                    <div className="modal-actions">
                        <button type="button" className="btn-secondary" onClick={onClose}>Cancel</button>
                        <button type="submit" className="auth-btn" disabled={saving}>
                        {saving ? "Saving…" : "Create"}
                        </button>
                    </div>
                    </form>
                </div>
            </div>
        </>
    )

}

export default EmployeeFormModal;