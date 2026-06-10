import About from "./pages/About/About";
import Login from "./pages/Login/Login";
import './App.css';
import { BrowserRouter as Router, Route, Link, Routes } from 'react-router-dom';

function App() {
  return <>
    <div className="landing-container">
      <header className="hero-section">
        <div className="hero-content">
          <h1>Employee Management System</h1>
          <p>
            A modern platform to manage employees, departments, attendance,
            performance, and organizational data efficiently.
          </p>

          <button
            className="login-btn"
            onClick={() => navigate("/login")}
          >
            Go to Login Page
          </button>
        </div>
      </header>

      <section className="features-section">
        <h2>Key Features</h2>

        <div className="feature-grid">
          <div className="feature-card">
            <h3>Employee Management</h3>
            <p>
              Create, update, search, and manage employee records with ease.
            </p>
          </div>

          <div className="feature-card">
            <h3>Dashboard Analytics</h3>
            <p>
              View employee statistics, department reports, and workforce
              insights.
            </p>
          </div>

          <div className="feature-card">
            <h3>Role Based Access</h3>
            <p>
              Secure access for Admin, HR, and Employees using JWT
              authentication.
            </p>
          </div>

          <div className="feature-card">
            <h3>Search & Filters</h3>
            <p>
              Quickly find employees using advanced search, sorting, and
              filtering.
            </p>
          </div>

          <div className="feature-card">
            <h3>Employee Profiles</h3>
            <p>
              Maintain complete employee information including department,
              designation, and joining date.
            </p>
          </div>

          <div className="feature-card">
            <h3>Scalable Architecture</h3>
            <p>
              Built with React, Golang, PostgreSQL, JWT, Docker, and Swagger.
            </p>
          </div>
        </div>
      </section>

      <section className="tech-section">
        <h2>Technology Stack</h2>

        <div className="tech-container">
          <div className="tech-box">
            <h3>Frontend</h3>
            <ul>
              <li>React</li>
              <li>React Router</li>
              <li>Axios</li>
              <li>Vite</li>
              <li>CSS / Tailwind CSS</li>
            </ul>
          </div>

          <div className="tech-box">
            <h3>Backend</h3>
            <ul>
              <li>Golang</li>
              <li>Gin Framework</li>
              <li>PostgreSQL</li>
              <li>GORM</li>
              <li>JWT Authentication</li>
            </ul>
          </div>
        </div>
      </section>

      <footer className="footer">
        <p>
          Employee Management System © 2026 | React + Golang + PostgreSQL
        </p>
      </footer>
    </div>
  </>;
}

export default App;