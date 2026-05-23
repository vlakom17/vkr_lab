import { useEffect, useState } from "react";
import { getMyChart, updateChart } from "../api/charts";
import { getMe, updateProfile } from "../api/users";
import { useNavigate } from "react-router-dom";
import MyChartSection from "../components/MyChartSection";
import ProfileSection from "../components/ProfileSection";
import ChartEditForm from "../components/ChartEditForm";

function ProfilePage() {
  const [user, setUser] = useState(null);
  const [chart, setChart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [editChartMode, setEditChartMode] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    about: "",
  });
  const [chartForm, setChartForm] = useState({
    title: "",
    description: "",
    genre: "",
    position_count: 10,
  });

  const handleChartSave = async () => {
    if (!chart?.id) return;
    setError("");

    try {
      const updated = await updateChart(chart.id, chartForm);

      setChart(updated);
      setEditChartMode(false);
    } catch (e) {
      setError("Ошибка обновления чарта");
    }
  };
  
  const handleSave = async () => {
    try {
      const dataToSend = {
        name: form.name,
        email: form.email,
        about: form.about,
      };

      if (form.password) {
        dataToSend.password = form.password;
      }

      await updateProfile(dataToSend);

      setUser((prev) => ({
        ...prev,
        ...dataToSend,
      }));

      setForm((prev) => ({
        ...prev,
        password: "",
      }));

      setEditMode(false);
    } catch (e) {
      setError("Ошибка обновления профиля");
    }
  };

  useEffect(() => {
    async function fetchData() {
      try {
        const userData = await getMe();
        setUser(userData);

        setForm({
          name: userData.name,
          email: userData.email,
          about: userData.about || "",
          password: "",
        });

        let chartData = null;

        try {
          chartData = await getMyChart();
        } catch {
          chartData = null;
        }

        if (chartData) {
          setChart(chartData);

          setChartForm({
            title: chartData.title || "",
            description: chartData.description || "",
            genre: chartData.genre || "",
            position_count: chartData.position_count || 10,
          });
        } else {
          setChart(null);
        }

      } catch (e) {
        console.error("Ошибка загрузки профиля:", e);
      } finally {
        setLoading(false);
      }
    }

    fetchData();
  }, []);
  
  if (loading) return <p style={{ padding: "20px" }}>Загрузка...</p>;
  if (!user) {
    return <p>Профиль не загружен</p>;
  }

  return (
    <div className="container">
      <h1>Личный кабинет</h1>

      <div className="profile-grid">
        
        <div className="card profile-card">
          <h3>Профиль</h3>

          <ProfileSection
            user={user}
            form={form}
            setForm={setForm}
            editMode={editMode}
            setEditMode={setEditMode}
            onSave={handleSave}
          />
        </div>

        <div className="card profile-card">
          <h3>Мой чарт</h3>

          <MyChartSection
            chart={chart}
            navigate={navigate}
            onEdit={() => setEditChartMode(true)}
          />
        </div>

      </div>

      {editChartMode && chart && (
        <ChartEditForm
          form={chartForm}
          setForm={setChartForm}
          onSave={handleChartSave}
          onCancel={() => setEditChartMode(false)}
        />
      )}

      <div className="profile-actions">
        <button onClick={() => navigate("/me/likes")}>
          👍 Понравившиеся чарты
        </button>

        <button onClick={() => navigate("/me/dislikes")}>
          👎 Непонравившиеся чарты
        </button>
      </div>
    </div>
  );
}

export default ProfilePage;