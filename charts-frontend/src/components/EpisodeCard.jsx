import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getChartByIdWithoutView } from "../api/charts";

function EpisodeCard({ episode }) {
  const navigate = useNavigate();

  const [chart, setChart] = useState(null);

  useEffect(() => {
    getChartByIdWithoutView(episode.ChartID)
      .then((data) => setChart(data))
      .catch(() => setChart(null));
  }, [episode.ChartID]);

  const date = episode.CreatedAt
    ? new Date(episode.CreatedAt).toLocaleDateString()
    : "—";

  return (
    <div
      className="card episode-card"
      onClick={() => navigate(`/episodes/${episode.ID}`)}
    >
      <div className="episode-date">📅 {date}</div>

      <div className="episode-title">
        {chart?.title || "Загрузка..."}

        {chart?.position_count && (
          <span className="episode-meta">
            {" "}• {chart.position_count} позиций
          </span>
        )}
      </div>
    </div>
  );
}

export default EpisodeCard;