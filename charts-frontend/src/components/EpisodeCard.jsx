import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getChartByIdWithoutView } from "../api/charts";

function EpisodeCard({ episode }) {
  const navigate = useNavigate();
  const [chart, setChart] = useState(null);
  const Id = episode.id || episode.ID;
  const chartId = episode.chart_id || episode.ChartID;
  const createdAt = episode.created_at || episode.CreatedAt;
  const date = createdAt
    ? new Date(createdAt).toLocaleDateString()
    : "—";

  useEffect(() => {
    getChartByIdWithoutView(chartId)
      .then((data) => setChart(data))
      .catch(() => setChart(null));
  }, [chartId]);

  return (
    <div
      className="card episode-card"
      onClick={() => navigate(`/episodes/${Id}`)}
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