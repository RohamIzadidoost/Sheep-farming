import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import SheepPage from './pages/SheepPage';
import LambingsPage from './pages/LambingsPage';
import VaccinationsPage from './pages/VaccinationsPage';
import VaccinesPage from './pages/VaccinesPage';
import TreatmentsPage from './pages/TreatmentsPage';
import { getToken } from './utils/api';

export default function App() {
  const token = getToken();
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={token ? <Navigate to="/dashboard" /> : <LoginPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/sheep" element={<SheepPage />} />
        <Route path="/lambings" element={<LambingsPage />} />
        <Route path="/vaccinations" element={<VaccinationsPage />} />
        <Route path="/vaccines" element={<VaccinesPage />} />
        <Route path="/treatments" element={<TreatmentsPage />} />
      </Routes>
    </BrowserRouter>
  );
}
