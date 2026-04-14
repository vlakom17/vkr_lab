function ChartEditForm({ form, setForm, onSave, onCancel }) {
  return (
    <div className="chart-edit">
      <h3>Редактирование чарта</h3>

      <div className="auth-form">
        <input
          className="auth-input"
          placeholder="Название"
          value={form.title}
          onChange={(e) =>
            setForm({ ...form, title: e.target.value })
          }
        />

        <textarea
          className="auth-input"
          placeholder="Описание"
          value={form.description}
          onChange={(e) =>
            setForm({ ...form, description: e.target.value })
          }
        />

        <input
          className="auth-input"
          placeholder="Жанр"
          value={form.genre}
          onChange={(e) =>
            setForm({ ...form, genre: e.target.value })
          }
        />

        <div className="select-block">
          <span className="muted">Количество позиций</span>

          <select
            className="auth-input"
            value={form.position_count}
            onChange={(e) =>
              setForm({
                ...form,
                position_count: Number(e.target.value),
              })
            }
          >
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={25}>25</option>
            <option value={30}>30</option>
            <option value={40}>40</option>
            <option value={50}>50</option>
          </select>
        </div>

        <div className="chart-edit-actions">
          <button className="auth-button" onClick={onSave}>
            Сохранить
          </button>
          <button onClick={onCancel}>
            Отмена
          </button>
        </div>
      </div>
    </div>
  );
}

export default ChartEditForm;