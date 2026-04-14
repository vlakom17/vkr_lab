function EpisodeListItem({ episode, topTrack, onClick }) {
  let date = "Нет даты";

  try {
    if (episode?.CreatedAt) {
      date = new Date(episode.CreatedAt).toLocaleDateString();
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