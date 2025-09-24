const config = {
  development: {
    API_BASE_URL: 'http://localhost:8080/',
    TIMEOUT: 10000,
    RETRY_ATTEMPTS: 3,
  }
};

const getEnvironment = () => {
  if (import.meta.env.MODE === 'production') return 'production';
  if (import.meta.env.MODE === 'staging') return 'staging';
  return 'development';
};

const currentConfig = config[getEnvironment()];

export default {
  ...currentConfig,
  // Auth specific endpoints
  AUTH_ENDPOINTS: {
    LOGIN: 'auth/login',
    REGISTER: 'auth/register',
  },
  // Storage keys
  STORAGE_KEYS: {
    USER: 'currentUser',
  },
  // HTTP status codes
  HTTP_STATUS: {
    UNAUTHORIZED: 401,
    FORBIDDEN: 403,
    NOT_FOUND: 404,
    INTERNAL_SERVER_ERROR: 500,
  }
};