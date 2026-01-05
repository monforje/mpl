import { useState } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import { useAuth } from '../context/AuthContext';
import './SendCSV.css';

const SendCSV = () => {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { token } = useAuth();

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setMessage(null);
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    if (!file) {
      setMessage({ text: 'Пожалуйста, выберите файл', type: 'error' });
      return;
    }

    setIsLoading(true);
    setMessage(null);

    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await fetch('http://localhost:3000/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        body: formData,
      });

      if (response.ok) {
        setMessage({ text: 'Файл успешно загружен!', type: 'success' });
        setFile(null);
        const fileInput = document.getElementById('csvFile') as HTMLInputElement;
        if (fileInput) fileInput.value = '';
      } else {
        const errorData = await response.json().catch(() => null);
        const errorMessage = errorData?.message || errorData?.error || `Ошибка ${response.status}: ${response.statusText}`;
        setMessage({ text: errorMessage, type: 'error' });
        console.error('Server error:', errorData);
      }
    } catch (error) {
      setMessage({ text: 'Ошибка соединения с сервером', type: 'error' });
      console.error('Upload error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="send-csv-container">
      <div className="csv-card">
        <h1>Загрузка CSV файла</h1>
        
        <form onSubmit={handleSubmit} className="csv-form">
          <div className="file-input-wrapper">
            <input
              type="file"
              id="csvFile"
              accept=".csv"
              onChange={handleFileChange}
              required
            />
            <label htmlFor="csvFile" className="file-label">
              {file ? file.name : 'Выберите CSV файл'}
            </label>
          </div>

          <button 
            type="submit" 
            className="submit-btn"
            disabled={isLoading || !file}
          >
            {isLoading ? 'Отправка...' : 'Отправить'}
          </button>

          {message && (
            <div className={`message ${message.type}`}>
              {message.text}
            </div>
          )}
        </form>
      </div>
    </div>
  );
};

export default SendCSV;
