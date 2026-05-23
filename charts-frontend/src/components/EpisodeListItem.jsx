function EpisodeListItem({ episode, topTrack, onClick }) {
  const createdAt = episode.created_at || episode.CreatedAt;
  let date = "Нет даты";

  try {
    if (createdAt) {
      date = new Date(createdAt).toLocaleDateString();
    }
  } catch {}

  return (
    <div className="episode-item" onClick={onClick}>
      <div className="episode-item-date">{date}</div>

      {topTrack ? (
        <div className="episode-item-track">
          🏆 {topTrack.artist} — {topTrack.title}
        </div>
      ) : (
        <div className="muted">Загрузка трека...</div>
      )}
    </div>
  );
}

export default EpisodeListItem;