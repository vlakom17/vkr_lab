import { useEffect, useState } from "react";
import { getLatestEpisodes } from "../api/archive";
import { getMe } from "../api/users";
import { isAuthenticated, removeToken } from "../utils/auth";
import EpisodeCard from "../components/EpisodeCard";
import PopularCharts from "../components/PopularCharts";
import RecommendedEpisodes from "../components/RecommendedEpisodes";

function HomePage() {
  const [episodes, setEpisodes] = useState([]);
  const [page, setPage] = useState(1);
  const pageSize = 10;
  const [tab, setTab] = useState("latest");

  useEffect(() => {
    async function loadEpisodes() {
      try {
        const data = await getLatestEpisodes(page, pageSize);
        setEpisodes(Array.isArray(data) ? data : []);
      } catch (e) {
        console.error("Ошибка загрузки эпизодов:", e);
        setEpisodes([]);
      }
    }
    loadEpisodes();
  }, [page]);

  useEffect(() => {
    async function checkAuth() {

      if (!isAuthenticated()) {
        return;
      }

      try {
        await getMe();
      } catch {
        removeToken();
      }
    }
    
    checkAuth();
  }, [tab]);

  const isFirstPage = page === 1;
  const isLastPage = !episodes || episodes.length < pageSize;

  return (
    
    <div className="container">
        <h2>Charter - интерактивный сервис пользовательских музыкальных чартов</h2>
        <div className="tabs">
    <button
      onClick={() => setTab("latest")}
      className={tab === "latest" ? "active" : ""}
    >
      Новые эпизоды
    </button>

    <button
      onClick={() => setTab("popular")}
      className={tab === "popular" ? "active" : ""}
    >
      Популярные чарты
    </button>

    <button
      onClick={() => setTab("recommended")}
      className={tab === "recommended" ? "active" : ""}
    >
      Может понравиться
    </button>
    </div>
    {tab === "latest" &&(
      <div className="pagination">
        {!isFirstPage && (
          <button onClick={() => setPage((p) => p - 1)}>
            ← Назад
          </button>
        )}

        <span className="muted">Страница {page}</span>

        {!isLastPage && (
          <button onClick={() => setPage((p) => p + 1)}>
            Вперёд →
          </button>
        )}
      </div>
    )}
      <div className="list">
        {tab === "latest" && (
          episodes.length === 0 ? (
            <p>Нет эпизодов</p>
          ) : (
            episodes.map((ep) => {
              const episodeId = ep.id || ep.ID;

              return (
                <EpisodeCard
                  key={episodeId}
                  episode={ep}
                />
              );
            })
          )
        )}
    {tab === "popular" && <PopularCharts />}
    {tab === "recommended" && <RecommendedEpisodes />}
    </div>
  </div>
);}

export default HomePage;