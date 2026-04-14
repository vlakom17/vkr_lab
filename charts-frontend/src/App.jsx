import { Routes, Route } from "react-router-dom";
import Navbar from "./components/NavBar";
import HomePage from "./pages/HomePage";
import EpisodePage from "./pages/EpisodePage";
import LoginPage from "./pages/LoginPage";
import ProtectedRoute from "./components/ProtectedRoute";
import ProfilePage from "./pages/ProfilePage";
import ChartPage from "./pages/ChartPage";
import RegisterPage from "./pages/RegisterPage";
import MyChartsPage from "./pages/RatedChartsPage";
import UserPage from "./pages/UserPage";
import CreateChartPage from "./pages/CreateChartPage";
import CreateEpisodePage from "./pages/CreateEpisodePage";

function App() {
  return (
    <>
      <Navbar />

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/episodes/:id" element={<EpisodePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/charts/:id" element={<ChartPage />} />
        <Route path="/users/:id" element={<UserPage />} />

        <Route element={<ProtectedRoute />}>
          <Route path="/me" element={<ProfilePage />} />
          <Route path="/me/likes" element={<MyChartsPage type="likes" />} />
          <Route path="/me/dislikes" element={<MyChartsPage type="dislikes" />} />
          <Route path="/create-chart" element={<CreateChartPage />} />
          <Route
            path="/charts/:id/create-episode"
            element={<CreateEpisodePage />}
          />
        </Route>
      </Routes>
    </>
  );
}

export default App;