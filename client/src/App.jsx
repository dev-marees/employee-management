import About from "./pages/About/About";
import Login from "./pages/Login/Login";
import logo from "./assets/logo.png";
import './App.css';
import { BrowserRouter as Router, Route, Link, Routes } from 'react-router-dom';

function App() {
  return <>
    <div className="landing-container">
      <section className="hero-section">
        <div className="hero-left">
          <h1>Employee Management System</h1>

          <p>
            Manage employees, departments, payroll information, and workforce
            data from a single platform.
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
          <div className="login-card">
            <h2>Employee Login</h2>

            <form>
              <div className="form-group">
                <label>Email</label>
                <input
                  type="email"
                  placeholder="Enter your email"
                />
              </div>

              <div className="form-group">
                <label>Password</label>
                <input
                  type="password"
                  placeholder="Enter your password"
                />
              </div>

              <button
                type="submit"
                className="login-btn"
              >
                Login
              </button>
            </form>
          </div>
        </div>
      </section>
      <section className="features-section">
        <h2>Why Choose Our System?</h2>

        <div className="feature-grid">
          <div className="feature-card">
            <h3>Employee Records</h3>
            <p>
              Maintain employee information in a centralized system.
            </p>
          </div>

          <div className="feature-card">
            <h3>Dashboard Reports</h3>
            <p>
              View workforce statistics and department-wise reports.
            </p>
          </div>

          <div className="feature-card">
            <h3>Secure Access</h3>
            <p>
              Role-based authentication for Admin, HR, and Employees.
            </p>
          </div>

          <div className="feature-card">
            <h3>Fast Search</h3>
            <p>
              Search employees instantly using filters and sorting.
            </p>
          </div>
        </div>
        <div className="company-logo">
          <div className="company-logo">
            <img src={logo} alt="Company Logo" />
          </div>
        </div>
      </section>
    </div>
  </>;
}

export default App;