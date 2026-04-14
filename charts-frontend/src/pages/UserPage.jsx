import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getUserById } from "../api/users";

function UserPage() {
  const { id } = useParams();
  const [user, setUser] = useState(null);

  useEffect(() => {
    getUserById(id)
      .then(setUser)
      .catch((e) => {
        console.error("Ошибка загрузки пользователя:", e);
        setUser(null);
      });
  }, [id]);

  if (!user) return <p style={{ padding: "20px" }}>Загрузка...</p>;

  return (
    <div className="container">
      <div className="card user-card">
        <h2 className="user-name">{user.name}</h2>

        <p>
          <span className="muted">О себе:</span>{" "}
          {user.about || "—"}
        </p>

        <p>
          <span className="muted">Дата регистрации:</span>{" "}
          {user.createdAt
            ? new Date(user.createdAt).toLocaleDateString()
            : user.created_at
            ? new Date(user.created_at).toLocaleDateString()
            : "—"}
        </p>
      </div>
    </div>
  );
}

export default UserPage;