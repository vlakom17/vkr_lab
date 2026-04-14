function ChartCard({ chart, author, onOpen }) {
  return (
    <div className="card chart-card">
      <div className="chart-card-content">
        <h3 className="chart-title">{chart.title}</h3>

        <p>
          <span className="muted">Автор:</span>{" "}
          <span>{author}</span>
        </p>

        <p>
          <span className="muted">Жанр:</span> {chart.genre}
        </p>

        <p>
          <span className="muted">Позиций:</span> {chart.position_count}
        </p>
      </div>

      <button className="auth-button" onClick={onOpen}>
        Открыть
      </button>
    </div>
  );
}

export default ChartCard;