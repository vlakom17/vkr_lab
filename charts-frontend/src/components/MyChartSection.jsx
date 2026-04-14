function MyChartSection({ chart, loading, navigate, onEdit }) {
  if (loading) return <p className="muted">Загрузка чарта...</p>;
    if (!chart) {
      return (
        <div className="chart-empty">
          <p className="muted">У вас пока нет чарта</p>

          <button
            className="auth-button"
            onClick={() => navigate("/create-chart")}
          >
            Создать чарт
          </button>
        </div>
      );
    }

    return (
      <div className="chart-info">
        <p className="chart-title">{chart.title}</p>

        <p>
          <span className="muted">Жанр:</span> {chart.genre}
        </p>

        <p>
          <span className="muted">Позиций:</span> {chart.position_count}
        </p>

        <div className="chart-actions">
          <button onClick={() => navigate(`/charts/${chart.id}`)}>
            Открыть
          </button>

          <button onClick={onEdit}>
            Редактировать
          </button>

          <button
            className="auth-button"
            onClick={() => navigate(`/charts/${chart.id}/create-episode`)}
            sstyle={{ color: "#fff" }}
          >
            <span className="plus">＋</span>
            Эпизод
          </button>
        </div>
      </div>
    )
  }
export default MyChartSection;