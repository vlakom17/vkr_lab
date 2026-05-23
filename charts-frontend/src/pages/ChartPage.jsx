import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getEpisodesByChart, getEpisodeById } from "../api/archive";
import { getChartStats, sendChartReaction, getMyReaction } from "../api/analysis";
import { getUserById } from "../api/users";
import { getChartById } from "../api/charts"
import EpisodeListItem from "../components/EpisodeListItem";

function ChartPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [chart, setChart] = useState(null);
  const [episodes, setEpisodes] = useState([]);
  const [topTracks, setTopTracks] = useState({});
  const [stats, setStats] = useState(null);
  const [loadingReaction, setLoadingReaction] = useState(false);
  const [myReaction, setMyReaction] = useState(null);
  const [userName, setUserName] = useState("Загрузка...");
  const [showUser, setShowUser] = useState(false);
  const [user, setUser] = useState(null);
  const [page, setPage] = useState(1);
  const PAGE_SIZE = 5;
  
  useEffect(() => {
    async function fetchData() {
      try {
        const [chartData, statsData, episodesData] = await Promise.all([
          getChartById(id),
          getChartStats(id),
          getEpisodesByChart(id),
        ]);

        setChart(chartData);
        setStats(statsData);
        setEpisodes(Array.isArray(episodesData) ? episodesData : []);
      } catch (e) {
        console.error("Ошибка загрузки чарта:", e);
      }
    }
    fetchData();
  }, [id]);

  const refreshStatsWithDelay = () => {
    setTimeout(async () => {
      try {
        const updated = await getChartStats(id);
        setStats(updated);
      } catch (e) {
        console.error("Ошибка обновления статистики:", e);
      }
    }, 500);
  };

  useEffect(() => {
    if (!chart?.user_id) return;

    async function fetchUser() {
      try {

        const userData = await getUserById(chart.user_id);
        setUser(userData);
        setUserName(userData.name);

      } catch (e) {
        console.error("Ошибка загрузки пользователя:", e);
        setUserName("Неизвестный пользователь");
      }
    }
    fetchUser();
  }, [chart?.user_id]);

  const handleReaction = async (type) => {
    if (loadingReaction) return;
    setLoadingReaction(true);

    try {
      const newType = myReaction === type ? "remove" : type;

      await sendChartReaction(id, newType);

      setMyReaction((prev) => (prev === type ? null : type));

      refreshStatsWithDelay();
    } catch (e) {
      console.error("Ошибка реакции:", e);
    } finally {
      setLoadingReaction(false);
    }
  };

  useEffect(() => {
    async function loadMyReaction() {
      try {
        const data = await getMyReaction(id);

        setMyReaction(data?.type || data?.Type || null);
      } catch (e) {
        console.error("Ошибка загрузки реакции:", e);
        setMyReaction(null);
      }
    }
    loadMyReaction();
  }, [id]);

  useEffect(() => {
    if (!episodes || episodes.length === 0) return;

    async function loadTopTracks() {
      try {
        const results = await Promise.all(
          episodes.map(async (ep) => {
            const episodeId = ep.id || ep.ID;
            if (!episodeId) {
             return [null, null];
            }
            try {
              const full = await getEpisodeById(episodeId);
              const tracks = full?.tracks || full?.Tracks;
              return [episodeId, tracks?.[0] || null];
            } catch {
              return [episodeId, null];
            }
          })
        );

        const map = Object.fromEntries(results);
        setTopTracks(map);
      } catch (e) {
        console.error("Ошибка загрузки топ треков:", e);
      }
    }
    loadTopTracks();
  }, [episodes]);

  if (!chart) return <p>Загрузка...</p>;
  const visibleEpisodes = (episodes || []).slice(0, page * PAGE_SIZE);

  const hasMore = visibleEpisodes.length < episodes.length;
  return (
    <div className="container">
      <div className="card chart-header">
        <h1>{chart.title}</h1>

        <div>
          <span className="muted">Автор:</span>{" "}
         
            <span
              className="user-toggle"
              onClick={() => setShowUser((prev) => !prev)}
            >
              {userName}
            </span>
                      
            {showUser && user && (
              <div className="user-inline">
                <p><b>О себе:</b> {user.about || "—"}</p>
                <p className="muted">
                  Зарегистрирован: {new Date(user.created_at).toLocaleDateString()}
                </p>
              </div>
            )}
        </div>

        <div className="chart-meta">
          <p><span className="muted">Жанр:</span> {chart.genre || "—"}</p>
          <p>
            <span className="muted">Позиций:</span> {chart.position_count}
          </p>
          <p><span className="muted">Описание:</span> {chart.description || "—"}</p>
          <p>
            <span className="muted">Создан:</span>{" "}
            {chart.created_at
              ? new Date(chart.created_at).toLocaleDateString()
              : "—"}
          </p>
        </div>
      </div>

      {stats && (
        <div className="card chart-stats">
          <div>👍 {stats.LikesCount}</div>
          <div>👎 {stats.DislikesCount}</div>
          <div>👁 {stats.ViewsCount}</div>
        </div>
      )}

      <div className="chart-reactions">
        <button
          className={myReaction === "like" ? "active" : ""}
          onClick={() => handleReaction("like")}
          disabled={loadingReaction}
        >
          👍 Нравится
        </button>

        <button
          className={myReaction === "dislike" ? "active" : ""}
          onClick={() => handleReaction("dislike")}
          disabled={loadingReaction}
        >
          👎 Не нравится
        </button>
      </div>

      <h2>Эпизоды</h2>

      <div className="list">
        {visibleEpisodes.map((ep, index) => {
          const episodeId = ep?.id || ep?.ID;
          const safeId = episodeId || index;
          const topTrack = episodeId ? topTracks?.[episodeId] : null;
          return (
            <EpisodeListItem
              key={safeId}
              episode={ep}
              topTrack={topTrack}
              onClick={() => episodeId && navigate(`/episodes/${episodeId}`)}
            />
          );
        })}
      </div>
      {hasMore && (
        <button
          className="load-more-button"
          onClick={() => setPage((p) => p + 1)}
        >
          Показать ещё
        </button>
      )}
    </div>
  );
}

export default ChartPage;