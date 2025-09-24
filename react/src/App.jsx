import { AuthProvider, useAuth } from './shared/AuthContext';
import AuthPanel from './features/auth/components/AuthPanel';
import Dashboard from './features/dashboard/Dashboard';

// Komponenta za loading state
const LoadingSpinner = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50">
    <div className="text-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
      <p className="mt-4 text-gray-600">Loading...</p>
    </div>
  </div>
);

// Glavni sadržaj aplikacije
const AppContent = () => {
  const { user, initializing } = useAuth();

  // Prikaži spinner dok se app inicijalizuje
  if (initializing) {
    return <LoadingSpinner />;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {user ? <Dashboard /> : <AuthPanel />}
    </div>
  );
};

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
