import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { getEpisodeById } from "../api/archive";
import { useNavigate } from "react-router-dom";
import { getChartByIdWithoutView } from "../api/charts";

function EpisodePage() {
  const { id } = useParams();
  const [episode, setEpisode] = useState(null);
  const [chartName, setChartName] = useState("Загрузка...");
  const chartId = episode?.chart_id || episode?.ChartID;
  const createdAt = episode?.created_at || episode?.CreatedAt;
  const tracks = episode?.tracks || episode?.Tracks;

  useEffect(() => {
    getEpisodeById(id)
      .then((data) => {
        setEpisode(data);
      })
      .catch((e) => console.error("Ошибка загрузки эпизода:", e));
  }, [id]);

  useEffect(() => {
    if (!episode) return;
    getChartByIdWithoutView(chartId)
      .then((data) => {
        setChartName(data.title);
      })
      .catch(() => setChartName("Неизвестный чарт"));
  }, [chartId]);

  const navigate = useNavigate();
  if (!episode) {
    return <p>Загрузка...</p>;
  }
  
  return (
    <div className="container">
      <button className="back-button" onClick={() => navigate(-1)}>
        ← Назад
      </button>

      <div className="card episode-header">
        <h1
          className="link"
          onClick={() => navigate(`/charts/${chartId}`)}
        >
          {chartName}
        </h1>

        <p className="muted">
          Эпизод от {new Date(createdAt).toLocaleDateString()}
        </p>
      </div>

      <h3>Треки</h3>

      {!tracks || tracks.length === 0 ? (
        <p>Нет треков</p>
      ) : (
        <table className="episode-table">
          <thead>
            <tr>
              <th>#</th>
              <th>Исполнитель</th>
              <th>Название</th>
              <th>Было</th>
              <th>Пик</th>
              <th className="appearances-column">Появлений в чарте</th>
              <th></th>
            </tr>
          </thead>

          <tbody>
            {[...tracks]
              .sort((a, b) => a.current_position - b.current_position)
              .map((track) => {
                let color = "var(--text)";
                let symbol = "";

                if (track.previous_position === 0) {
                  color = "#5039e1";
                  symbol = "NEW";
                } else if (track.current_position < track.previous_position) {
                  color = "green";
                  symbol = "↑";
                } else if (track.current_position > track.previous_position) {
                  color = "red";
                  symbol = "↓";
                } else {
                  color = "orange";
                  symbol = "=";
                }

                const peak =
                  track.times_at_peak_position > 1
                    ? `${track.highest_position} (${track.times_at_peak_position})`
                    : `${track.highest_position}`;

                return (
                  <tr
                    key={track.current_position}
                    className="episode-row"
                  >
                    <td>
                      <strong>{track.current_position}</strong>
                    </td>

                      <td className="track-main">{track.artist}</td>

                      <td className="track-main">{track.title}</td>

                    <td className="center">
                      {track.previous_position === 0
                        ? "—"
                        : track.previous_position}
                    </td>
                    
                    <td className="center">{peak}</td>

                    <td className="center">{track.episodes_count}</td>

                    <td className="track-change" style={{ color }}>
                      {symbol}
                    </td>

                    <td className="track-links">
                      <a href={track.listen_links.apple_music} target="_blank" rel="noopener noreferrer">
                        Apple
                      </a>{" "}
                      |{" "}
                      <a href={track.listen_links.yandex_music} target="_blank" rel="noopener noreferrer">
                        Yandex
                      </a>
                    </td>
                  </tr>
                );
              })}
          </tbody>
       </table>
    )}
    </div>
  );
}

export default EpisodePage;