// API configuration
export const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api.halooid.com/v1';

// Organization ID - in a real app, this would be dynamically set based on the user's context
export const DEFAULT_ORG_ID = import.meta.env.VITE_DEFAULT_ORG_ID || '00000000-0000-0000-0000-000000000000';

// Task status colors
export const STATUS_COLORS = {
  todo: '#9E9E9E',
  in_progress: '#2196F3',
  review: '#FF9800',
  done: '#4CAF50',
  cancelled: '#F44336'
};

// Task priority colors
export const PRIORITY_COLORS = {
  low: '#4CAF50',
  medium: '#FFC107',
  high: '#FF9800',
  critical: '#F44336'
};

// Project status colors
export const PROJECT_STATUS_COLORS = {
  planning: '#9C27B0',
  active: '#2196F3',
  on_hold: '#FF9800',
  completed: '#4CAF50',
  cancelled: '#F44336'
};

// Date format options
export const DATE_FORMAT_OPTIONS: Intl.DateTimeFormatOptions = {
  year: 'numeric',
  month: 'short',
  day: 'numeric'
};

// Date-time format options
export const DATE_TIME_FORMAT_OPTIONS: Intl.DateTimeFormatOptions = {
  year: 'numeric',
  month: 'short',
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit'
};
