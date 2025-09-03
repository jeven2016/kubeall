import LoginPage from "@/modules/login/LoginPage";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {BrowserRouter, Navigate, Outlet, Route, Routes} from "react-router";
import ModuleIndex from "@/modules/layout/ModuleIndex.tsx";
import ImageList from "@/modules/vm/image/ImageList.tsx";
import ImageEdit from "@/modules/vm/image/ImageEdit.tsx";
import VmList from "@/modules/vm/VmList.tsx";
import DeploymentList from "@/modules/container/workload/DeploymentList.tsx";
import StatefulSetList from "@/modules/container/workload/StatefulSetList.tsx";
import DaemonSetList from "@/modules/container/workload/DaemonSetList.tsx";
import PodList from "@/modules/container/workload/PodList.tsx";
import JobList from "@/modules/container/workload/JobList.tsx";
import ServiceList from "@/modules/container/access/ServiceList.tsx";
import IngressList from "@/modules/container/access/IngressList.tsx";
import PvcList from "@/modules/container/storage/PvcList.tsx";
import StorageClassList from "@/modules/container/storage/StorageClassList.tsx";
import ConfigMapList from "@/modules/container/config/ConfigMapList.tsx";
import SecretList from "@/modules/container/config/SecretList.tsx";
import VmEdit from "@/modules/vm/VmEdit.tsx";
import IpPoolList from "@/modules/vm/ippool/IpPoolList.tsx";

// 创建 QueryClient
const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: 0, // 默认重试 1 次
            // retry: false, //不重试
            //no cache, unit: mill
            staleTime: 0,
            gcTime: 0
            // staleTime: 5 * 60 * 1000 // 数据新鲜时间 5 分钟
        }
    }
});

export default function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<><Outlet/></>}>
                        <Route index element={<SecuredRoot/>}/>
                        <Route index path="login" element={<LoginPage/>}/>
                        <Route path="home">
                            <Route path="vmMgr" element={<ModuleIndex/>}>
                                <Route path="index" element={<div>index</div>}/>
                                <Route path="vm" element={<VmList/>}/>
                                <Route path="vm/:action" element={<VmEdit/>}/>
                                <Route path="image" element={<ImageList/>}/>
                                <Route path="ipPool" element={<IpPoolList/>}/>
                                <Route path="image/:action" element={<ImageEdit/>}/>
                                <Route path="*" element={<div>404</div>}/>
                            </Route>
                            <Route path="containerMgr" element={<ModuleIndex/>}>
                                <Route path="index" element={<div>index</div>}/>
                                <Route path="Deployment" element={<DeploymentList/>}/>
                                <Route path="StatefulSet" element={<StatefulSetList/>}/>
                                <Route path="DaemonSet" element={<DaemonSetList/>}/>
                                <Route path="Pod" element={<PodList/>}/>
                                <Route path="Job" element={<JobList/>}/>
                                <Route path="Service" element={<ServiceList/>}/>
                                <Route path="Ingress" element={<IngressList/>}/>
                                <Route path="pvc" element={<PvcList/>}/>
                                <Route path="sc" element={<StorageClassList/>}/>
                                <Route path="cm" element={<ConfigMapList/>}/>
                                <Route path="secret" element={<SecretList/>}/>
                                <Route path="*" element={<div>404</div>}/>
                            </Route>
                        </Route>
                        <Route path="*" element={<div>404</div>}/>
                    </Route>

                </Routes>
            </BrowserRouter>
        </QueryClientProvider>
    );
}

const SecuredRoot = () => {
    if (true) {
        return <Navigate to="/login"/>;
    }
    return <Navigate to="/home"/>;
};
