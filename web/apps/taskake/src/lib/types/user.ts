export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface UserResponse {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
}
