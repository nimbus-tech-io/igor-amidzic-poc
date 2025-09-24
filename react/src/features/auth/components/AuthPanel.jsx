import { useState } from 'react';
import { useAuth } from '../../../shared/AuthContext';
import Login from './Login';
import Register from './Register';

export default function AuthPanel() {
  const [isLogin, setIsLogin] = useState(true);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {isLogin ? (
          <Login />
        ) : (
          <Register />
        )}
        
        <div className="text-center mt-6">
          <button
            onClick={() => setIsLogin(!isLogin)}
            className="text-sm text-gray-600 hover:text-gray-900 transition duration-200 ease-in-out"
          >
            {isLogin 
              ? "Don't have an account? " 
              : "Already have an account? "
            }
            <span className="font-semibold text-blue-600 hover:text-blue-800">
              {isLogin ? 'Sign up' : 'Sign in'}
            </span>
          </button>
        </div>
      </div>
    </div>
  );
}