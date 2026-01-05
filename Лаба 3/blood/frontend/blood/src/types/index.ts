export interface User {
  id: string;
  name: string;
  email: string;
}

export interface UserContextType {
  user: User | null;
  isLoggedIn: boolean;
  login: (user: User) => void;
  logout: () => void;
}

export interface Scan {
  id: string;
  user_id: string;
  created_at: string;
  full_name: string;
  birth_date: string;
  sex: string;
  hemoglobin: number;
  erythrocytes: number;
  hematocrit: number;
  mcv: number;
  leukocytes: number;
  neutrophils: number;
  lymphocytes: number;
  monocytes: number;
  eosinophils: number;
  basophils: number;
  platelets: number;
  mpv: number;
}
