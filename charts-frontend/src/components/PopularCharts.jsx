import { useEffect, useState } from "react";
import { getPopularCharts,} from "../api/charts";
import { getChartStats } from "../api/analysis"
import { useNavigate } from "react-router-dom";
import { getUserById } from "../api/users";

function PopularCharts() {
  const [popularCharts, setPopularCharts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState({});
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchPopular() {
      setLoading(true);

      try {
        const charts = await getPopularCharts();

        if (!Array.isArray(charts)) {
          setPopularCharts([]);
          return;
        }

        const chartsWithStats = await Promise.all(
          charts.map(async (chart) => {
            try {
              const stats = await getChartStats(chart.id);
              return { ...chart, stats };
            } catch {
              return {
                ...chart,
                stats: { LikesCount: 0, DislikesCount: 0 },
              };
            }
          })
        );

        setPopularCharts(chartsWithStats);
      } catch (e) {
        console.error("Ошибка загрузки популярных:", e);
        setPopularCharts([]);
      } finally {
        setLoading(false);
      }
    }

    fetchPopular();
  }, []);

  useEffect(() => {
  if (!popularCharts.length) return;

  async function loadUsers() {
    try {
      const uniqueIds = [...new Set(popularCharts.map(c => c.user_id))];

      const results = await Promise.all(
        uniqueIds.map(async (id) => {
          try {
            const user = await getUserById(id);
            return [id, user.name];
          } catch {
            return [id, "Неизвестно"];
          }
        })
      );
      setUsers(Object.fromEntries(results));
    } catch (e) {
      console.error("Ошибка пользователей:", e);
    }
  }
  loadUsers();
}, [popularCharts]);

  if (loading) return <p>Загрузка...</p>;
  if (popularCharts.length === 0) return <p>Нет популярных чартов</p>;

  return (
    <>
     {popularCharts.map((chart) => (
      <div
        key={chart.id}
        className="card chart-card"
      >
        <div className="chart-card-content">
          <h3 className="chart-title">{chart.title}</h3>

          <p>
            <span className="muted">Автор:</span>{" "}
            {users[chart.user_id] || "Загрузка..."}
          </p>

          <p>
            <span className="muted">Жанр:</span> {chart.genre}
          </p>

          <p>
            <span className="muted">Позиций:</span> {chart.position_count}
          </p>

          <div className="chart-stats">
            <span>👍 {chart.stats?.LikesCount ?? 0}</span>
            <span>👎 {chart.stats?.DislikesCount ?? 0}</span>
          </div>
        </div>

        <button
          className="auth-button"
          onClick={() => navigate(`/charts/${chart.id}`)}
        >
          Открыть
        </button>
      </div>
    ))}
    </>
  );
}

export default PopularCharts;