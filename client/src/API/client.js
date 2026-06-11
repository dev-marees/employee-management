const BASE_URL = import.meta.env.VITE_API_URL;

export async function apiRequest(path, {method = "GET", body, auth = false} = {}) {
    const headers = {"Content-Type": "application/json"};

    if (auth) {
        const token = localStorage.getItem("access_token")
        if (token) headers["Authorization"] = `Bearer ${token}`
    }

    const res = await fetch(`${BASE_URL}${path}`, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
    });

    const json = await res.json();

    if (!res.ok || !json.success) {
        const err = new Error(json.message || "Request failed")
        err.status = res.status;
        err.fields = json.errors;
        throw err
    }

    return json.data;
}