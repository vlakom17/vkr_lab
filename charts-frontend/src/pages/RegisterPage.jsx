import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { registerUser } from "../api/users";

function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [about, setAbout] = useState("");
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const handleRegister = async (e) => {
    e.preventDefault();
    setError("");

    try {
      await registerUser(name, email, password, about);
      navigate("/");
    } catch (err) {
      setError("Пользователь уже существует или данные неверны");
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h2>Регистрация</h2>

        <form onSubmit={handleRegister} className="auth-form">
          <input
            className="auth-input"
            placeholder="Имя"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />

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

          <textarea
            className="auth-input"
            placeholder="О себе"
            value={about}
            onChange={(e) => setAbout(e.target.value)}
            rows={3}
          />

          <button className="auth-button" type="submit">
            Зарегистрироваться
          </button>
        </form>

        {error && <p className="auth-error">{error}</p>}
      </div>
    </div>
  );
}

export default RegisterPage;