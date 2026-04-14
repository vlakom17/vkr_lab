import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginUser } from "../api/users";
import { setToken } from "../utils/auth";

function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  async function handleLogin(e) {
    e.preventDefault();
    setError("");

    try {
      const data = await loginUser(email, password);

      setToken(data.token);

      navigate("/");
    } catch (err) {
      setError("Неверный email или пароль");
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h2>Вход</h2>

        <form onSubmit={handleLogin} className="auth-form">
          <input
            className="auth-input"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />

          <input
            className="auth-input"
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <button className="auth-button" type="submit">
            Войти
          </button>
        </form>

        {error && <p className="auth-error">{error}</p>}
      </div>
    </div>
  );
}

export default LoginPage;