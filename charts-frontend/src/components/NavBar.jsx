import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { isAuthenticated, removeToken } from "../utils/auth";
import { logout } from "../api/users";

function Navbar() {
  const [auth, setAuth] = useState(isAuthenticated);
  const navigate = useNavigate();

  useEffect(() => {
    setAuth(isAuthenticated());
  }, []);

  useEffect(() => {
    function handleAuthChange() {
      setAuth(isAuthenticated());
    }

    window.addEventListener("authChanged", handleAuthChange);

    return () => {
      window.removeEventListener("authChanged", handleAuthChange);
    };
  }, []);

  const handleLogout = async () => {
    try {
      await logout();
    } catch (e) {
      console.warn("Logout request failed");
    }

    removeToken();
    navigate("/");
  };

  return (
    <div className="header">
      <div className="header-left">
      </div>

      <div className="header-logo" onClick={() => navigate("/")}>
        Charter
      </div>

      <div className="header-right">
        {!auth ? (
          <>
            <button onClick={() => navigate("/login")}>Войти</button>
            <button onClick={() => navigate("/register")}>Регистрация</button>
          </>
        ) : (
          <>
            <button onClick={() => navigate("/me")}>Личный кабинет</button>
            <button onClick={handleLogout}>Выйти</button>
          </>
        )}
      </div>
    </div>
  );
}

export default Navbar;