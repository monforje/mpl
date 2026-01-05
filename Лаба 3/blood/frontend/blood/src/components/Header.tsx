import { Link } from 'react-router-dom';
import UserMenu from './UserMenu';
import './Header.css';

interface HeaderProps {
  currentPath?: string;
}

const Header = ({ currentPath }: HeaderProps) => {
  return (
    <header className="header">
      <nav className="nav">
        <div className="nav-left">
          <Link 
            to="/" 
            className={`nav-link ${currentPath === '/' ? 'active' : ''}`}
          >
            Главная
          </Link>
          <Link 
            to="/send-csv" 
            className={`nav-link ${currentPath === '/send-csv' ? 'active' : ''}`}
          >
            Отправить анализ
          </Link>
          <Link 
            to="/scans" 
            className={`nav-link ${currentPath === '/scans' ? 'active' : ''}`}
          >
            Мои анализы
          </Link>
        </div>
        <div className="nav-right">
          <UserMenu />
        </div>
      </nav>
    </header>
  );
};

export default Header;
