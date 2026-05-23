import { useEffect, useState } from "react";
import { getRecommendations } from "../api/analysis";
import RecommendedEpisodeCard from "./RecommendedEpisodeCard";

function RecommendedEpisodes() {
  const [recommended, setRecommended] = useState([]);
  const [loading, setLoading] = useState(false);

  async function fetchRecommended() {
    setLoading(true);

    try {
      const data = await getRecommendations();
      setRecommended(Array.isArray(data) ? data : []);
    } catch (e) {
      console.error("Ошибка рекомендаций:", e);
      setRecommended([]);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchRecommended();
  }, []);

  useEffect(() => {
    function handleAuthChange() {
      fetchRecommended();
    }

    window.addEventListener("authChanged", handleAuthChange);

    return () => {
      window.removeEventListener("authChanged", handleAuthChange);
    };
  }, []);
  
  if (loading) return <p>Загрузка...</p>;
  if (recommended.length === 0) return <p>Нет рекомендаций</p>;

  return (
    <>
      {recommended.map((ep) => (
        <RecommendedEpisodeCard key={ep.ID || ep.id} episode={ep} />
      ))}
    </>
  );
}

export default RecommendedEpisodes;