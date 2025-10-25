import axios, { type AxiosInstance } from 'axios';

const globalApi: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    }
});

globalApi.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem("token")
        if (!(token && config.headers)){
            throw new Error('no token found')

        }
        config.headers.Authorization = `Bearer ${token}`;
        return config;
    },
    (error) => Promise.reject(error)
);


export default globalApi;