import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { login, register, saveSession } from "../API/auth";

function AuthForm() {
	const navigate = useNavigate();
	const [isLogin, setIsLogin] = useState(true);
	const [form, setForm] = useState({
		name: "", email: "", password: "", confirmPassword: "",
	});
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	function update(e) {
		setForm({ ...form, [e.target.name]: e.target.value });
	}

	async function handleSubmit(e) {
		e.preventDefault();
		setError("");
		setLoading(true);
		try {
			let data;
			if (isLogin) {
				data = await login({ email: form.email, password: form.password });
			} else {
				if (form.password !== form.confirmPassword) {
					throw new Error("Passwords do not match");
				}
				data = await register({
					name: form.name, email: form.email, password: form.password,
				});
			}
			saveSession(data);          // store token + user
			navigate("/dashboard");     // go to dashboard on success
		} catch (err) {
			setError(err.message);
		} finally {
			setLoading(false);
		}
	}

  return (
	<div className="auth-card">
		<h1>Employee Management System</h1>  
		<form onSubmit={handleSubmit}>
			{isLogin ? <h2>Login</h2> : <h2>Create Account</h2>}

			{!isLogin && (
				<input name="name" type="text" placeholder="Full Name"
					value={form.name} onChange={update} />
			)}

			<input name="email" type="email" placeholder="Email Address"
				value={form.email} onChange={update} />

			<input name="password" type="password" placeholder="Password"
				value={form.password} onChange={update} />

			{!isLogin && (
				<input name="confirmPassword" type="password" placeholder="Confirm Password"
				value={form.confirmPassword} onChange={update} />
			)}

			{error && <p className="auth-error">{error}</p>}

			<button className="auth-btn" type="submit" disabled={loading}>
				{loading ? "Please wait..." : isLogin ? "Login" : "Register"}
			</button>
			<p className="switch-text">
				{isLogin ? "New user? " : "Already have an account? "}
				<span onClick={() => setIsLogin(!isLogin)}>
					{isLogin ? "Register" : "Login"}
				</span>
			</p>
		</form>
	</div>
  );
}

export default AuthForm;
