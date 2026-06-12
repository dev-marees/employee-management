import { apiRequest } from "./client";

const SESSION_MAX_MS = 24 * 60 * 60 * 1000;

export function register({name, email, password, role}) {
    return apiRequest("/auth/register", {
        method: "POST",
        body: { name, email, password, ...(role?{role}:{}) }
    });
}

export function login({email, password}) {
    return apiRequest("/auth/login", {
        method: "POST",
        body: {email, password},
    });
}

export function saveSession(data) {
    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("refresh_token", data.refresh_token);
    localStorage.setItem("user", JSON.stringify(data.user));
    localStorage.setItem("session_expiry", String(Date.now() + SESSION_MAX_MS));
}

export function logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    localStorage.removeItem("session_expiry");
}