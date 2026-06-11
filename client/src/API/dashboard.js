import { apiRequest } from "./client";

export function getStats() {
  return apiRequest("/dashboard", { auth: true });
}

// const user = JSON.parse(localStorage.getItem("user"));

// if (user.role === "Admin" || user.role === "HR") {
//   const stats = await getStats();              // /dashboard
//   const emps  = await listEmployees({ limit: 10 });  // emps.items
// } else {
//   const me = await getMyProfile();             // /me
// }