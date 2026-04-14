function ProfileSection({ user, form, setForm, editMode, setEditMode, onSave }) {
  if (!editMode) {
    return (
      <div className="profile-info">
        <p><span className="muted">Имя:</span> {user.name}</p>
        <p><span className="muted">Email:</span> {user.email}</p>
        <p><span className="muted">
            Дата регистрации:</span>{" "}
            {new Date(user.createdAt).toLocaleDateString()}
          </p>
        <p><span className="muted">О себе:</span> {user.about || "—"}</p>
        <button
          onClick={() => setEditMode(true)}
          style={{ marginTop: "10px" }}
        >
          Редактировать
        </button>
      </div>
    );
  }

  return (
    <div className="auth-form">
      <input
        className="auth-input"
        placeholder="Имя"
        value={form.name}
        onChange={(e) =>
          setForm({ ...form, name: e.target.value })
        }
      />

      <input
        className="auth-input"
        placeholder="Email"
        value={form.email}
        onChange={(e) =>
          setForm({ ...form, email: e.target.value })
        }
      />

      <textarea
        className="auth-input"
        placeholder="О себе"
        value={form.about}
        onChange={(e) =>
          setForm({ ...form, about: e.target.value })
        }
      />

      <input
        className="auth-input"
        type="password"
        placeholder="Новый пароль"
        value={form.password}
        onChange={(e) =>
          setForm({ ...form, password: e.target.value })
        }
      />

      <div className="profile-actions-row">
        <button className="auth-button" onClick={onSave}>
          Сохранить
        </button>
        <button onClick={() => setEditMode(false)}>
          Отмена
        </button>
      </div>
    </div>
  );
}

export default ProfileSection;