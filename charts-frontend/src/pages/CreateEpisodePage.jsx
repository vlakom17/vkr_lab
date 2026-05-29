import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { createEpisode, searchTracks } from "../api/archive";
import { getChartByIdWithoutView } from "../api/charts";

function capitalizeWords(str = "") {
  return str
    .toLowerCase()
    .replace(/(^|[\s\-–—([{])([a-zа-яё])/gi, (match, prefix, letter) =>
      prefix + letter.toUpperCase()
    );
}

function CreateEpisodePage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [chart, setChart] = useState(null);
  const [tracks, setTracks] = useState([]);
  const [suggestions, setSuggestions] = useState([]);
  const [activeIndex, setActiveIndex] = useState(null);
  const [searchValue, setSearchValue] = useState("");
  const suggestionsRef = useRef(null);

  useEffect(() => {
    async function loadChart() {
      setError("");
      try {
        const data = await getChartByIdWithoutView(id);
        setChart(data);

        const initial = Array.from(
          { length: data.position_count },
          () => ({ artist: "", title: "" })
        );

        setTracks(initial);
      } catch (e) {
        console.error(e);
        setError("Ошибка загрузки чарта");
      }
    }
    loadChart();
  }, [id]);

  useEffect(() => {
    if (!searchValue || searchValue.length < 2) {
      setSuggestions([]);
      return;
    }
    let currentValue = searchValue;

    const timeout = setTimeout(async () => {
      try {
        const res = await searchTracks(currentValue);
        if (currentValue === searchValue) {
          setSuggestions(Array.isArray(res) ? res : []);
        }
      } catch (e) {
        console.error(e);
        setSuggestions([]);
      }
    }, 400);
    return () => clearTimeout(timeout);
  }, [searchValue]);

  useEffect(() => {
    function handleClickOutside(e) {
      if (suggestionsRef.current && !suggestionsRef.current.contains(e.target)) {
        setSuggestions([]);
        setActiveIndex(null);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const updateTrack = (index, field, value) => {
    const updated = [...tracks];
    updated[index][field] = value;
    setTracks(updated);
  };

  const handleSubmit = async (e) => {
    e?.preventDefault();
    setError("");
    try {
      const prepared = tracks
        .filter((t) => t.artist && t.title)
        .map((t, i) => ({
          artist: t.artist,
          title: t.title,
          current_position: i + 1,
        }));

      if (prepared.length !== tracks.length) {
        setError("Заполните все позиции");
        return;
      }

      await createEpisode(id, { tracks: prepared });

      navigate(`/charts/${id}`, { state: { refresh: true } });
    } catch (e) {
      console.error(e);
      setError("Ошибка создания эпизода");
    }
  };

  if (!chart) return <p style={{ padding: "20px" }}>Загрузка...</p>;

  return (
    <div className="container">
      <div className="card">
        <h1>Создание эпизода</h1>
        <p className="muted">{chart.title}</p>

        <div style={{ marginTop: "20px" }}>
          <div className="episode-form">
            {tracks.map((track, index) => (
              <div key={index} className="track-row">
                
                <div className="track-index">#{index + 1}</div>

                <input
                  className="auth-input"
                  placeholder="Название трека"
                  value={track.title}
                  onChange={(e) => {
                    const value = e.target.value;

                    updateTrack(index, "title", value);
                    setActiveIndex(index);
                    setSearchValue(value);
                  }}
                />

                <input
                  className="auth-input"
                  placeholder="Исполнитель"
                  value={track.artist}
                  onChange={(e) =>
                    updateTrack(index, "artist", e.target.value)
                  }
                  onFocus={() => setActiveIndex(index)}
                />

                {activeIndex === index && suggestions.length > 0 && (
                  <div className="suggestions" ref={suggestionsRef}>
                    {suggestions.map((s, i) => (
                      <div
                        key={i}
                        className="suggestion-item"
                        onClick={() => {
                          updateTrack(index, "title",  capitalizeWords(s.title));
                          updateTrack(index, "artist",  capitalizeWords(s.artist));

                          setSuggestions([]);
                          setSearchValue("");
                          setActiveIndex(null);
                        }}
                      >
                        {capitalizeWords(s.artist)} — {capitalizeWords(s.title)}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
          <button
            className="auth-button"
            style={{ marginTop: "20px" }}
            onClick={handleSubmit}
          >
            Создать эпизод
          </button> 
          <p
            className="muted"
            style={{
              marginTop: "16px",
              fontSize: "12px",
              lineHeight: "1.4",
              color: "#6d28d9",
            }}
          >
            При добавлении в эпизод ремикса указывайте автора ремикса в поле
            «Исполнитель», а не в названии трека.
          </p>
      </div>
  </div>
  );
}

export default CreateEpisodePage;