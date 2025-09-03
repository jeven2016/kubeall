import {VmFormData} from "@/types";

export const VmDefaultValues:VmFormData = {
    name: "",
    namespace: "",
    displayName: "",
    osType: "windows",
    description: "",
    arch: "amd64",
    cpu: "1",
    memory: "2",
    enableDisc: true,
    discName: "rk",
    windowsDriverDisc: "bj",
    systemDisk: {
        type: "pvc",
        image: "",
        capacity: "10",
        volumeType: "Block",
        storageClass: "",
    },
    dataDisks: [
        {
            type: "pvc",
            image: "",
            capacity: "10",
            volumeType: "Block",
            storageClass: "",
        },
        {
            type: "pvc",
            image: "",
            capacity: "10",
            volumeType: "Block",
            storageClass: "",
        }
    ],
    primaryNetwork: {
        networkType: "pod",
        deviceModel: "bridge",
        interfaceType: "virtio"
    },
    advancedSettings:{
        bootLoader: "BIOS",
        cpuModel: "host-model",
        timeZone: "Asia/Shanghai"
    }
}