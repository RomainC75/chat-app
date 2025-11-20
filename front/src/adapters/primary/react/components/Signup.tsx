import React, { useState } from 'react'
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  Container,
} from '@mui/material'

import { useDispatch, useSelector } from 'react-redux'
import { signupUser, type AppDispatch, type RootState } from '../../../../store/store';
import { Link, useNavigate } from 'react-router-dom';

const Signup = () => {
    const dispach = useDispatch<AppDispatch>();
    const navigate = useNavigate();
    const { isLoading } = useSelector(
    (state: RootState) => state.user
  );

  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: ''
  })


  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validateForm()){
        console.log("---> password don't match")
        return
    }
    await dispach(signupUser({signupUser:{email: formData.email, password: formData.password}})).unwrap()
    navigate("/login")
  }

  const validateForm = (): boolean => {
    if(formData.confirmPassword != formData.password){
        return false
    }
    return true
  }

  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Paper elevation={3} sx={{ padding: 4, width: '100%' }}>
          <Typography component="h1" variant="h4" align="center" gutterBottom>
            Sign Up
          </Typography>
          


          <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              autoFocus
              value={formData.email}
              onChange={handleChange}
            //   error={!!errors.email}
            //   helperText={errors.email}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="new-password"
              value={formData.password}
              onChange={handleChange}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="confirmPassword"
              label="Confirm Password"
              type="password"
              id="confirmPassword"
              autoComplete="new-password"
              value={formData.confirmPassword}
              onChange={handleChange}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
              disabled={isLoading}
            >
              {isLoading ? 'Signing Up...' : 'Sign Up'}
            </Button>
            <Link to="/login">Login</Link>
          </Box>
        </Paper>
      </Box>
    </Container>
  )
}

export default Signup