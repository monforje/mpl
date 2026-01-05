import { useState, useRef, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './UserMenu.css';

const UserMenu = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { user, logout } = useAuth();
  const menuRef = useRef<HTMLDivElement>(null);
  const iconRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        menuRef.current &&
        iconRef.current &&
        !menuRef.current.contains(event.target as Node) &&
        !iconRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, []);

  const toggleMenu = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsOpen(!isOpen);
  };

  const handleLogout = () => {
    logout();
    setIsOpen(false);
  };

  return (
    <div className="user-menu-container">
      <div 
        className="user-icon" 
        onClick={toggleMenu}
        ref={iconRef}
      />
      
      <div 
        className={`user-popup ${isOpen ? 'active' : ''}`}
        ref={menuRef}
      >
        {!user ? (
          <div className="popup-guest">
            <Link to="/login" className="popup-link" onClick={() => setIsOpen(false)}>Войти</Link>
            <Link to="/register" className="popup-link" onClick={() => setIsOpen(false)}>Зарегистрироваться</Link>
          </div>
        ) : (
          <div className="popup-user">
            <div className="popup-user-info">
              <div className="popup-user-name">{user.name}</div>
              <div className="popup-user-email">{user.email}</div>
            </div>
            <button onClick={handleLogout} className="popup-link logout-btn">Выйти</button>
          </div>
        )}
      </div>
    </div>
  );
};

export default UserMenu;
