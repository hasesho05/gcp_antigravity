import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { onAuthStateChanged, signOut, type User as FirebaseUser } from 'firebase/auth';
import { auth } from '@/lib/firebase';
import { axios } from '@/lib/axios';
import Axios from 'axios';
import type { User } from '@/types/api';

interface AuthContextType {
  user: User | null;
  firebaseUser: FirebaseUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [firebaseUser, setFirebaseUser] = useState<FirebaseUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (currentUser: FirebaseUser | null) => {
      setFirebaseUser(currentUser);
      if (currentUser) {
        try {
          const token = await currentUser.getIdToken();
          axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

          // Fetch the user data from our backend
          try {
            const response = await axios.get<User>('/users/me');
            setUser(response.data);
          } catch (error) {
            if (Axios.isAxiosError(error) && error.response?.status === 404) {
              console.log('User not found in backend, creating new user...');
              // User doesn't exist in backend yet, try to create/sync
              try {
                const providerData = currentUser.providerData[0];
                const providerId = providerData?.providerId || 'password';
                
                const createResponse = await axios.post<User>('/users', {
                  email: currentUser.email,
                  provider: providerId,
                });
                setUser(createResponse.data);
              } catch (createError) {
                console.error('Failed to create/sync user:', createError);
                if (Axios.isAxiosError(createError)) {
                  console.error('Create user response:', createError.response?.data);
                  console.error('Create user status:', createError.response?.status);
                  
                  // If conflict (409), it means the user actually exists (maybe race condition or email conflict).
                  // In this case, we should try to fetch the user again or handle it gracefully.
                  if (createError.response?.status === 409) {
                     console.log('User already exists (conflict), retrying fetch...');
                     try {
                        const retryResponse = await axios.get<User>('/users/me');
                        setUser(retryResponse.data);
                     } catch (retryError) {
                        console.error('Retry fetch failed:', retryError);
                        setUser(null);
                     }
                  } else {
                    setUser(null);
                  }
                } else {
                  setUser(null);
                }
              }
            } else {
              console.error('Failed to fetch user profile:', error);
              setUser(null);
            }
          }
        } catch (error) {
          console.error('Failed to get token or auth error:', error);
          setUser(null);
        }
      } else {
        delete axios.defaults.headers.common['Authorization'];
        setUser(null);
      }
      setIsLoading(false);
    });

    return () => unsubscribe();
  }, []);

  const logout = async () => {
    try {
      await signOut(auth);
      setUser(null);
      setFirebaseUser(null);
      delete axios.defaults.headers.common['Authorization'];
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <AuthContext.Provider value={{ user, firebaseUser, isLoading, isAuthenticated: !!user, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
