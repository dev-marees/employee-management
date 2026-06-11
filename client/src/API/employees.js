import { apiRequest } from "./client";

//Get Employees
export function listEmployees(params = {}) {
    const qs = new URLSearchParams(params).toString();
    return apiRequest(`/employees/${qs ? `?${qs}`: ""}`, {auth: true})
}

//Get EmployeeById
export function getEmployeeById(id) {
    return apiRequest(`/employees/${id}`, {auth: true})
}

// PUT /employees/:id  (Admin/HR only)
export function updateEmployee(id, payload) {
    return apiRequest(`/employees/${id}`, { method: "PUT", body: payload, auth: true });
}

// DELETE /employees/:id  (Admin only)
export function deleteEmployee(id) {
    return apiRequest(`/employees/${id}`, { method: "DELETE", auth: true });
}

export function createEmployee(payload) {
  return apiRequest("/employees", { method: "POST", body: payload, auth: true });
}

export function getMyProfile() {
  return apiRequest("/me", { auth: true });
}