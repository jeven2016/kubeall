import axios, {AxiosRequestConfig,} from 'axios';


const createApiClient = () => {
    const instance = axios.create({
        // baseURL: "",
        // timeout: 5 * 1000
    });

    instance.interceptors.request.use(
        (config) => {
            return config;
        },
        (error) => {
            return Promise.reject(error);
        }
    );
    instance.interceptors.response.use(
        (response) => {
            const {data, status} = response
            if (status >= 400) {
                return Promise.reject(data.message)
            }
            return data;
        },
        (error: any) => {
            return Promise.reject(error);
        }
    );

    const invoke = <T = any>(config: AxiosRequestConfig): Promise<T> => {
        return instance.request<any, T>(config);
    };

    return {
        get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
            return invoke({...config, method: 'GET', url: url});
        },

        post<T = any>(url: string, data: T, config?: AxiosRequestConfig): Promise<T> {
            return invoke({...config, method: 'POST', url: url, data: data});
        },

        put<T = any>(url: string, data: T, config: AxiosRequestConfig): Promise<T> {
            return invoke({...config, method: 'PUT', url: url, data: data});
        },

        delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
            return invoke({...config, method: 'DELETE', url: url,});
        },

        upload(url: string, data: FormData, config?: AxiosRequestConfig): Promise<FormData> {
            return invoke({
                ...config, headers: {"Content-Type": "multipart/form-data"},
                method: 'POST', url: url, data: data
            });
        },
    }
}

const apiClient = createApiClient();
export default apiClient;
