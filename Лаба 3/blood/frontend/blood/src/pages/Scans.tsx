import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import type { Scan } from '../types';
import './Scans.css';

const Scans = () => {
  const [scans, setScans] = useState<Scan[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useAuth();

  useEffect(() => {
    fetchScans();
  }, []);

  const fetchScans = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch('http://localhost:3002/scans', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Не удалось загрузить анализы');
      }

      const data = await response.json();
      setScans(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки данных');
    } finally {
      setIsLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  };

  if (isLoading) {
    return (
      <div className="scans-container">
        <div className="loading">Загрузка анализов...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="scans-container">
        <div className="error-message">{error}</div>
        <button onClick={fetchScans} className="retry-btn">
          Попробовать снова
        </button>
      </div>
    );
  }

  if (scans?.length === 0) {
    return (
      <div className="scans-container">
        <div className="empty-state">
          <h2>Анализы не найдены</h2>
          <p>У вас пока нет загруженных анализов крови</p>
        </div>
      </div>
    );
  }

  return (
    <div className="scans-container">
      <h1>Мои анализы крови</h1>
      
      <div className="scans-grid">
        {scans.map((scan) => (
          <div key={scan.id} className="scan-card">
            <div className="scan-header">
              <h3>{scan.full_name}</h3>
              <span className="scan-date">{formatDate(scan.created_at)}</span>
            </div>

            <div className="scan-info">
              <div className="info-row">
                <span className="label">Дата рождения:</span>
                <span className="value">{formatDate(scan.birth_date)}</span>
              </div>
              <div className="info-row">
                <span className="label">Пол:</span>
                <span className="value">{scan.sex === 'М' ? 'Мужской' : 'Женский'}</span>
              </div>
            </div>

            <div className="scan-results">
              <h4>Показатели крови</h4>
              
              <div className="results-grid">
                <div className="result-item">
                  <span className="param-name">Гемоглобин</span>
                  <span className="param-value">{scan.hemoglobin} г/л</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Эритроциты</span>
                  <span className="param-value">{scan.erythrocytes} ×10¹²/л</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Гематокрит</span>
                  <span className="param-value">{scan.hematocrit} %</span>
                </div>
                <div className="result-item">
                  <span className="param-name">MCV</span>
                  <span className="param-value">{scan.mcv} фл</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Лейкоциты</span>
                  <span className="param-value">{scan.leukocytes} ×10⁹/л</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Нейтрофилы</span>
                  <span className="param-value">{scan.neutrophils} %</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Лимфоциты</span>
                  <span className="param-value">{scan.lymphocytes} %</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Моноциты</span>
                  <span className="param-value">{scan.monocytes} %</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Эозинофилы</span>
                  <span className="param-value">{scan.eosinophils} %</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Базофилы</span>
                  <span className="param-value">{scan.basophils} %</span>
                </div>
                <div className="result-item">
                  <span className="param-name">Тромбоциты</span>
                  <span className="param-value">{scan.platelets} ×10⁹/л</span>
                </div>
                <div className="result-item">
                  <span className="param-name">MPV</span>
                  <span className="param-value">{scan.mpv} фл</span>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Scans;
