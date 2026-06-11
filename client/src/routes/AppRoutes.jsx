import { Routes, Route, Navigate } from "react-router-dom";
import App from "../App";
import Login from "../pages/Login/Login";
import Dashboard from "../pages/Dashboard/Dashboard";
import ProtectedRoute from "../components/ProtectedRoute";

function AppRoutes() {
    return(
        <Routes>
            <Route path="/" element= { <App /> } />
            <Route path="/dashboard" element= {
                <ProtectedRoute>
                    <Dashboard  />
                </ProtectedRoute>
            } />
            <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
    )
}

export default AppRoutes