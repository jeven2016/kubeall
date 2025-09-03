import {MenuItem} from "@/types";

export default function MenuData(): MenuItem[] {
    return [
        {
            order: 1, name: "vmMgr", displayName: "虚拟化管理", link: "", children: [
                {order: 1, name: "vm", displayName: "虚拟机", link: ""},
                {order: 2, name: "image", displayName: "镜像", link: ""},
                {order: 2, name: "ipPool", displayName: "IP地址池", link: ""},
                {order: 2, name: "disk", displayName: "磁盘", link: ""},
                {order: 2, name: "network", displayName: "网络", link: ""},
                {order: 2, name: "snapshot", displayName: "快照", link: ""},
                {order: 2, name: "backup", displayName: "备份", link: ""},
            ]
        },
        {
            order: 2, name: "containerMgr", displayName: "容器云", link: "", children: [
                {
                    order: 1, name: "workload", displayName: "工作负载", link: "", children: [
                        {order: 1, name: "Deployment", displayName: "Deployment", link: ""},
                        {order: 1, name: "StatefulSet", displayName: "StatefulSet", link: ""},
                        {order: 1, name: "DaemonSet", displayName: "DaemonSet", link: ""},
                        {order: 1, name: "Pod", displayName: "Pod", link: ""},
                        {order: 1, name: "Job", displayName: "Job", link: ""},
                    ]
                },
                {
                    order: 2, name: "ac", displayName: "访问控制", link: "", children: [
                        {order: 1, name: "Service", displayName: "Service", link: ""},
                        {order: 1, name: "Ingress", displayName: "Ingress", link: ""},
                        {order: 1, name: "NetworkPolicy", displayName: "网络策略", link: ""},
                    ]
                },
                {
                    order: 3, name: "storage", displayName: "存储", link: "", children: [
                        {order: 1, name: "pvc", displayName: "存储卷申明", link: ""},
                        {order: 1, name: "volumeSnapshot", displayName: "Volume快照", link: ""},
                        {order: 1, name: "sc", displayName: "存储类", link: ""}
                    ]
                },
                {
                    order: 4, name: "cfg", displayName: "配置", link: "", children: [
                        {order: 1, name: "cm", displayName: "Config Map", link: ""},
                        {order: 1, name: "secret", displayName: "Secret", link: ""},
                    ]
                },
                {
                    order: 5, name: "node", displayName: "节点", link: "",
                },

            ],
        },
        {order: 3, name: "tenant", displayName: "租户", link: "",},
        {order: 4, name: "system", displayName: "系统设置", link: "",}
    ] as MenuItem[]
}

