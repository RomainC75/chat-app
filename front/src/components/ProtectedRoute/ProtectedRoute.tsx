import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import Loader from "../Loader/Loader";
import { verify, type AppDispatch, type RootState } from '../../store/store';


interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>()
  const { connectedUser } = useSelector((state: RootState) => ({
    connectedUser: state.user
  }));

  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token")
    console.log("-> token ", token)
    if (!token) {
      navigate("/login")
    }else {
      dispatch(verify())
    }
    
  }, []);

  // Afficher un loader pendant la vérification de l'authentification
  if (isChecking || connectedUser.isLoading || !connectedUser.user) {
    return <Loader />;
  }

  // L'utilisateur est authentifié et a le rôle requis (ou aucun rôle spécifique n'est requis)
  return <>{children}</>;
};

export default ProtectedRoute;
