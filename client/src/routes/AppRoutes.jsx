import { Routes, Route, Navigate } from "react-router-dom";
import App from "../App";
import Login from "../pages/Login/Login";

function AppRoutes() {
    return(
        <Routes>
            <Route path="/" element= { <App /> } ></Route>
            <Route path="/login" element= { <Login /> } ></Route>
        </Routes>
    )
}

export default AppRoutes