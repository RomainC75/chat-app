import axios, { type AxiosInstance } from 'axios';

const globalApi: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    }
});


export default globalApi;