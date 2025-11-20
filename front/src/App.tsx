import Signup from "./adapters/primary/react/components/Signup";
import "./App.css";
import { Navigate, Route, Routes, BrowserRouter as Router, } from "react-router-dom";
import Login from "./adapters/primary/react/components/Login";
import ProtectedRoute from "./adapters/primary/react/components/ProtectedRoute/ProtectedRoute";
import Rooms from "./adapters/primary/react/components/Rooms/Rooms";

function App() {
  return (
    <Router>
      <div className="app">
        <Routes>
          <Route
            path="/signup"
            element={
                <Signup />
            }
          />
          <Route
            path="/login"
            element={
                <Login />
            }
          /> 
          <Route
            path="/*"
            element={
              <ProtectedRoute>
                <div>
                  <h1>Application</h1>
                  <Rooms/>
                </div>
              </ProtectedRoute>
            }
          />

          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
