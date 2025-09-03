import {constants, vmLabelKeys} from "@/modules/common/Constants.tsx";
import {diskData, VmFormData} from "@/types";
import {PersistentVolumeClaim} from "kubernetes-models/v1";
import {V1VirtualMachine} from "@kubevirt-ui/kubevirt-api/kubevirt";

export const getImageTypeText = (imageType?: string) => {
    if (imageType == constants.imageType.iso) {
        return "光盘安装镜像"
    }
    if (imageType == constants.imageType.disk) {
        return "磁盘镜像"
    }
    return "";
}


export const createVm = (data: VmFormData) => {
    const vm: V1VirtualMachine = {
        apiVersion: 'kubevirt.io/v1',
        kind: 'VirtualMachine',
        metadata: {
            name: data.name,
            namespace: data.namespace,
            annotations: {
                [vmLabelKeys.osType]: data.osType,
                [vmLabelKeys.description]: data.description ?? "",
            },
            labels: {
                [vmLabelKeys.vmName]: data.name,
            }
        },
        spec: {
            template: {
                spec: {
                    architecture: data.arch,
                    domain: {
                        cpu: {
                            sockets: 1,
                            threads: 1,
                            cores: parseInt(data.cpu),
                        },
                        resources: {
                            requests: {
                                memory: data.memory,
                            },
                        },
                        devices: {
                            disks: [],
                            interfaces: [],
                        },
                    },
                    volumes: [],
                    networks: [],
                },
            },
        },
    };
    return vm;
}

export const addVolumes = (vm: V1VirtualMachine, volumeNamePrefix: string,
                           disks: diskData[]): void => {
    const pvcs: PersistentVolumeClaim[] = [];

    disks.forEach((d, index) => {
        const name = index == 0 ? volumeNamePrefix : `${volumeNamePrefix}-${index}`;
        const vmSpec = vm.spec.template.spec!!;
        if (d.type === constants.hostDiskType) {
            vm.spec.template.spec!!.volumes!!.push({
                name: name,
                hostDisk: {
                    capacity: d.capacity,
                    path: d.hostDiskPath!!,
                    type: "DiskOrCreate",
                }
            });
            return;
        }
        const randString = (Math.random() + 1).toString(36).substring(7);
        const pvcName = `${vm.metadata?.name}-${randString}`;
        const sc = d.type === constants.pvcDiskType ? d.storageClass : d.image;

        vmSpec.volumes!!.push({name: name, persistentVolumeClaim: {claimName: pvcName}});
        pvcs.push(new PersistentVolumeClaim({
            metadata: {
                name: pvcName,
                namespace: vm.metadata?.namespace,
                labels: {
                    [vmLabelKeys.vmName!!]: vm.metadata!!.name!!,
                    [vmLabelKeys.vmDisk]: "true"
                }
            },
            spec: {
                accessModes: ["ReadWriteOnce"],
                resources: {requests: {storage: d.capacity}},
                storageClassName: sc,
                volumeMode: d.volumeType,
            }
        }));

    })
    vm.metadata!!.annotations!![vmLabelKeys.pvcTemplates] = JSON.stringify(pvcs);
}