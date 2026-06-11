import { Navigate } from "react-router-dom";

function ProtectedRoute({ children, allowedRoles }) {
  const token = localStorage.getItem("access_token");

  // Not logged in → send to landing (which has your login form)
  if (!token) return <Navigate to="/" replace />;

  // Optional role gate
  if (allowedRoles) {
    const user = JSON.parse(localStorage.getItem("user") || "null");
    if (!user || !allowedRoles.includes(user.role)) {
      return <Navigate to="/" replace />;
    }
  }

  return children;
}

export default ProtectedRoute;
