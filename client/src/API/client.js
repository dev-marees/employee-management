const BASE_URL = import.meta.env.VITE_API_URL;

let refreshPromis = null;

function sessionExpired() {
    const exp = Number(localStorage.getItem("session_expiry") || 0);
    return exp > 0 && Date.now() > exp;
}

function forceLogout() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("user");
  localStorage.removeItem("session_expiry");
  window.location.href = "/"; // back to the login/landing page
}

async function doRefresh() {
    const refresh_token = localStorage.getItem("refresh_token");
    if (!refresh_token) throw new Error("No refresh token found");

    const res = await fetch(`${BASE_URL}/auth/refresh`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({ refresh_token }),
    });
    
    const json = await res.json();
    if (!res.ok || !json.success) throw new Error("refresh failed");

    localStorage.setItem("access_token", json.data.access_token);
    localStorage.setItem("refresh_token", json.data.refresh_token);
    return json.data.access_token;
}

export async function apiRequest(path, {method = "GET", body, auth = false, _retry = false} = {}) {
    const headers = {"Content-Type": "application/json"};

    if (auth && sessionExpired()) {
        forceLogout();
        throw new Error("session expired");
    }

    if (auth) {
        const token = localStorage.getItem("access_token")
        if (token) headers["Authorization"] = `Bearer ${token}`
    }

    const res = await fetch(`${BASE_URL}${path}`, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
    });

    if (res.status === 401 && auth ** !_retry) {
        try {
            if (!refreshPromis) {
                refreshPromis = doRefresh().finally(() => { refreshPromis = null; });
            }
            await refreshPromis;
            return apiRequest(path, { method, body, auth, _retry: true });
        } catch {
            forceLogout();
            throw new Error("session expired");
        }
    }

    const json = await res.json();

    if (!res.ok || !json.success) {
        const err = new Error(json.message || "Request failed")
        err.status = res.status;
        err.fields = json.errors;
        throw err
    }

    return json.data;
}