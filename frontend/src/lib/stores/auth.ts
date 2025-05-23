import { writable } from 'svelte/store';
import axios from 'axios';
import { browser } from '$app/environment';

export interface User {
  id: string;
  username: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
}

function createAuth() {
  const { subscribe, set, update } = writable<AuthState>({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false
  });

  return {
    subscribe,
    login: async (username: string, password: string): Promise<boolean> => {
      try {
        const response = await axios.post('http://localhost:8080/api/auth/login', {
          username,
          password
        });
        const { access_token, refresh_token, user } = response.data;
        if (browser) {
          localStorage.setItem('access_token', access_token);
          localStorage.setItem('refresh_token', refresh_token);
        }
        set({
          user,
          accessToken: access_token,
          refreshToken: refresh_token,
          isAuthenticated: true
        });
        axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
        return true;
      } catch (err) {
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
        return false;
      }
    },
    register: async (username: string, password: string): Promise<boolean> => {
      try {
        await axios.post('http://localhost:8080/api/auth/register', {
          username,
          password
        });
        return true;
      } catch (err) {
        return false;
      }
    },
    logout: (): void => {
      if (browser) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
      }
      set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
      delete axios.defaults.headers.common['Authorization'];
    },
    refreshToken: async (): Promise<boolean> => {
      if (!browser) return false;
      const refreshToken = localStorage.getItem('refresh_token');
      if (!refreshToken) return false;
      try {
        const response = await axios.post('http://localhost:8080/api/auth/refresh', {}, {
          headers: { 'Refresh-Token': refreshToken }
        });
        const { access_token, refresh_token, user } = response.data;
        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);
        set({
          user,
          accessToken: access_token,
          refreshToken: refresh_token,
          isAuthenticated: true
        });
        axios.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
        return true;
      } catch (err) {
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
        return false;
      }
    },
    initialize: (): void => {
      if (!browser) return;
      const accessToken = localStorage.getItem('access_token');
      const refreshToken = localStorage.getItem('refresh_token');
      if (accessToken && refreshToken) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
        set({ user: null, accessToken, refreshToken, isAuthenticated: true });
      }
    }
  };
}

export const auth = createAuth(); 