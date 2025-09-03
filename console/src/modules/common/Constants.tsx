
export const resourceType = {
    clusterResources: "clusters",
    namespaceResources: "namespaces",
    allNamespace: "all",
    namespace: "namespaces",
    vmImage: "images",
    vm: "vms",
    deployment: "deployments",
    statefulSet: "statefulsets",
    DaemonSet: "daemonsets",
    pod: "pods",
    job: "jobs",
    service: "services",
    ingress: "ingresses",
    pvc: "pvcs",
    pv: "pvs",
    sc: "sc",
    configMap: "configmaps",
    secret: "secrets",
    node: "nodes",
    globalSettings: "globalsettings",
    ipPool: "pools"
}

export const constants = {
    imageNamespace: "longhorn-system",
    imageBackend: "backingimage",
    annotations: {
        displayName: "api.kubeall.io/displayName"
    },
    arm64: "arm64",
    uefi: "uefi",
    cpuPassthrough: "host-passthrough",
    pvcDiskType: "pvc",
    hostDiskType: "hostDisk",
    imageDiskType: "image",

    imageType: {
        iso: "iso",
        disk: "disk"
    }
};

export const queryKey = {
    globalSettings: "globalSettings",

}

export const vmLabelKeys = {
    vmName: "kubevirt.io/domain",
    osType: "kubeall.io/osType",
    description: "kubeall.io/desc",
    vmDisk: "kubeall.io/vmDisk",
    pvcTemplates: "kubeall.io/pvcTemplates"
}

/**
 * 当屏幕窗口变化后，组件的折叠/展开的参数
 */
export const LAYOUT_SETTING = {
    CONTENT: {MIN: '50px', MAX: '300px'},
    LEFT_SLIDER: {MIN: '-16.625rem', MAX: '0rem'}
};

export const buildResourceUri = (scop: string, type: string) => {
    return `/api/v1/${scop}/${type}`;
}

export const buildNamespacedUri = (namespace: string, type: string) => {
    return `/api/v1/namespaces/${namespace}/${type}`;
}

export const buildClusterUri = (type: string) => {
    return `/api/v1/clusters/${type}`;
}