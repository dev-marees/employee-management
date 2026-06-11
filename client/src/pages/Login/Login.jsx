import { useState } from "react";
import { useNavigate } from "react-router-dom";
import './Login.css'

function Login() {
    const [isLogin, setIsLogin] = useState(false)
    const navigate = useNavigate();
    const [form, setForm] = useState( { email: "", password: "" } )
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    async function handleSubmit(e) {
        e.preventDefault();
    }

  return (
            <>
                <div className="landing-container">
                    <section className="hero-section">
                        <div className="hero-left">
                            <h1>Employee Management System</h1>

                            <p>
                            Manage employees, departments, payroll information,
                            reports, and workforce analytics from a single platform.
                            </p>

                            <div className="hero-features">
                            <div>✓ Employee Management</div>
                            <div>✓ Dashboard Analytics</div>
                            <div>✓ Search & Filters</div>
                            <div>✓ Role Based Access</div>
                            <div>✓ JWT Authentication</div>
                            <div>✓ Docker Deployment</div>
                            </div>
                        </div>

                        <div className="hero-right">
                            <div className="auth-card">

                            {isLogin ? (
                                <>
                                <h2>Login</h2>

                                <input
                                    type="email"
                                    placeholder="Email Address"
                                />

                                <input
                                    type="password"
                                    placeholder="Password"
                                />

                                <button className="auth-btn">
                                    Login
                                </button>

                                <p className="switch-text">
                                    New user?{" "}
                                    <span onClick={() => setIsLogin(false)}>
                                    Register
                                    </span>
                                </p>
                                </>
                            ) : (
                                <>
                                <h2>Create Account</h2>

                                <input
                                    type="text"
                                    placeholder="Full Name"
                                />

                                <input
                                    type="email"
                                    placeholder="Email Address"
                                />

                                <input
                                    type="password"
                                    placeholder="Password"
                                />

                                <input
                                    type="password"
                                    placeholder="Confirm Password"
                                />

                                <button className="auth-btn">
                                    Register
                                </button>

                                <p className="switch-text">
                                    Already have an account?{" "}
                                    <span onClick={() => setIsLogin(true)}>
                                    Login
                                    </span>
                                </p>
                                </>
                            )}
                            </div>
                        </div>
                    </section>
                </div>
            </>
        );
}

export default Login;