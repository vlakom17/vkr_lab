import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { getChartByIdWithoutView } from "../api/charts";

function capitalizeWords(str) {
  return str?.replace(/\b\w/g, (c) => c.toUpperCase());
}

function RecommendedEpisodeCard({ episode }) {
    const navigate = useNavigate();
    const [chartName, setChartName] = useState("Загрузка...");
    const topTrack = episode.tracks?.[0];

    const date = episode.created_at
        ? new Date(episode.created_at).toLocaleDateString()
        : "—";

    useEffect(() => {
        const chartId = episode.ChartID || episode.chart_id;

        if (!chartId) return;

        getChartByIdWithoutView(chartId)
            .then((data) => setChartName(data.title))
            .catch(() => setChartName("Неизвестный чарт"));
        }, [episode]);

    return (
        <div
        className="card episode-card rec-card"
        onClick={() => navigate(`/episodes/${episode.ID || episode.id}`)}
        >
        <div className="rec-left">
            <div className="episode-date">📅 {date}</div>
            <div className="episode-title">{chartName}</div>
        </div>

        {topTrack && (
            <div className="rec-right">
            🏆 {capitalizeWords(topTrack.artist)} — {capitalizeWords(topTrack.title)}
            </div>
        )}
    </div>
  );
}
export default RecommendedEpisodeCard;