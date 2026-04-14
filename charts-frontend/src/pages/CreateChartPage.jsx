import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createChart } from "../api/charts";

function CreateChartPage() {
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const [form, setForm] = useState({
    title: "",
    description: "",
    genre: "",
    position_count: 10,
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await createChart(form);

      navigate("/me");
    } catch (e) {
      console.error(e);
      
      setError("Ошибка создания чарта");
    }
  };

  return (
    <div className="container">
      <div className="auth-card">
        <h1>Создание чарта</h1>

        <form onSubmit={handleSubmit} className="auth-form">
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

          <button className="auth-button" type="submit">
            Создать
          </button>
        </form>

        {error && <p className="auth-error">{error}</p>}
      </div>
    </div>
  );
}

export default CreateChartPage;