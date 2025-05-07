import { writable } from 'svelte/store';
import { DEFAULT_ORG_ID } from '$lib/config';
import type { UserResponse } from '$lib/types/user';

// Create writable stores
const currentUser = writable<UserResponse | null>(null);
const users = writable<UserResponse[]>([]);
const isLoading = writable<boolean>(false);
const error = writable<string | null>(null);

// API base URL
const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api.halooid.com/v1';

// User API functions
const fetchCurrentUser = async () => {
  isLoading.set(true);
  error.set(null);
  
  try {
    const response = await fetch(`${API_BASE_URL}/users/me`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch current user: ${response.statusText}`);
    }
    
    const data = await response.json();
    currentUser.set(data);
    isLoading.set(false);
    return data;
  } catch (err) {
    console.error('Error fetching current user:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

const fetchUsers = async () => {
  isLoading.set(true);
  error.set(null);
  
  try {
    const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/users`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch users: ${response.statusText}`);
    }
    
    const data = await response.json();
    users.set(data.users || []);
    isLoading.set(false);
    return data.users;
  } catch (err) {
    console.error('Error fetching users:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

// Export the store
export const userStore = {
  currentUser: {
    subscribe: currentUser.subscribe
  },
  users: {
    subscribe: users.subscribe
  },
  isLoading: {
    subscribe: isLoading.subscribe
  },
  error: {
    subscribe: error.subscribe
  },
  fetchCurrentUser,
  fetchUsers
};
