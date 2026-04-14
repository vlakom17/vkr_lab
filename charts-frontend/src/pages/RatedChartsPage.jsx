import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  getMyLikedCharts,
  getMyDislikedCharts,
} from "../api/charts";
import ChartCard from "../components/ChartCard";
import { getUserById } from "../api/users";

function RatedChartsPage({ type }) {
  const [charts, setCharts] = useState([]);
  const navigate = useNavigate();
  const [users, setUsers] = useState({});

  useEffect(() => {
    async function fetchData() {
      try {
        const data =
          type === "likes"
            ? await getMyLikedCharts()
            : await getMyDislikedCharts();

        setCharts(Array.isArray(data) ? data : []);
      } catch {
        console.error("Ошибка загрузки чартов");
      }
    }
    fetchData();
  }, [type]);

  useEffect(() => {
    if (!charts.length) return;

    async function loadUsers() {
      try {
        const uniqueIds = [...new Set(charts.map(c => c.user_id))];

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
        console.error("Ошибка загрузки пользователей:", e);
      }
    }
    loadUsers();
  }, [charts]);

  return (
    <div className="container">
      <h2>
        {type === "likes"
          ? "Понравившиеся чарты"
          : "Непонравившиеся чарты"}
      </h2>

      {charts.length === 0 ? (
        <p>Пока ничего нет</p>
      ) : (
        charts.map((chart) => (
          <ChartCard
            key={chart.id}
            chart={chart}
            author={users[chart.user_id] || "Загрузка..."}
            onOpen={() => navigate(`/charts/${chart.id}`)}
          />
        ))
      )}
    </div>
  );
}

export default RatedChartsPage;