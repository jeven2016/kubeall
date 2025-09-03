import {Button, Collapse, Form, Typography} from "react-windy-ui";
import {useLocation, useNavigate} from "react-router";
import {useFieldArray} from "react-hook-form";
import VmBasicForm from "@/modules/vm/vm/VmBasicForm.tsx";
import VmDiscForm from "@/modules/vm/vm/VmDiscForm.tsx";
import VmDisksForm from "@/modules/vm/vm/VmDisksForm.tsx";
import VmNetworkForm from "@/modules/vm/vm/VmNetworkForm.tsx";
import {VmDefaultValues} from "@/modules/vm/vm/VmDefaultValues.tsx";
import VmAdvancedForm from "@/modules/vm/vm/VmAdvancedForm.tsx";
import {VmFormData} from "@/types";
import {constants} from "@/modules/common/Constants.tsx";
import {addVolumes, createVm} from "@/modules/common/VmCommon.tsx";
import {V1Interface} from "@kubevirt-ui/kubevirt-api/kubevirt/models";


export default function VmEdit() {
    const navigate = useNavigate();
    const location = useLocation();

    const {form, watch, control, setValue, getValue} = Form.useForm({
        mode: 'onSubmit',
        shouldFocusError: false,
        defaultValues: VmDefaultValues,
    });

    const goBack = () => navigate("/home/vmMgr/vm");


    const onSubmit = (data: VmFormData, e) => {
        const vm = createVm(data);

        const spec = vm.spec.template.spec!!;
        const domain = spec.domain;
        //bootloader
        const advancedSettings = data.advancedSettings;
        if (advancedSettings.bootLoader == constants.uefi) {
            domain.firmware = {
                bootloader: {
                    efi: {
                        secureBoot: false,
                    },
                },
            };
        }

        //cpu model
        domain.cpu!!.model = data.advancedSettings.cpuModel;
        if (data.arch == constants.arm64) {
            domain.cpu!!.model = constants.cpuPassthrough;
        }

        //timezone
        if (data.advancedSettings.timeZone) {
            domain.clock = {
                timezone: data.advancedSettings.timeZone,
            };
        }

        // for disks
        domain.devices.disks = [];
        domain.devices.disks.push({
            bootOrder: 1,
            disk: {
                bus: "virtio",
            },
            name: "system-disk",
        });
        if (data.enableDisc) {
            domain.devices.disks.push({
                bootOrder: 2,
                cdrom: {bus: "sata",},
                name: "disc",
            });
        }

        data.dataDisks.forEach((disk, index) => {
            const name = index == 0 ? "data-disk" : `data-disk-${index}`;
            domain.devices.disks!!.push({
                disk: {
                    bus: "virtio",
                },
                name: name,
            })
        });

        // for volumes
        addVolumes(vm, "system-disk", [data.systemDisk]);
        addVolumes(vm, "data-disk", data.dataDisks);

        const netInter: V1Interface = {name: "default", model: data.primaryNetwork.deviceModel}
        if (data.primaryNetwork.networkType == "nat") {
            netInter.masquerade = {};
        }
        if (data.primaryNetwork.networkType == "bridge") {
            netInter.bridge = {};
        }
        domain.devices.interfaces!!.push(netInter);
        spec.networks!!.push({name: "default", pod: {}});

        console.log(vm);
        return
    };

    const {fields: dataDiskFields, append, remove} = useFieldArray({
        control,
        name: "dataDisks", // it corresponds to the field defined in defaultValues
    });

    return <>
        <Form form={form} onSubmit={onSubmit} direction="horizontal" labelCol={{col: 2}} controlCol={{col: 4}}>
            <div className="c-content-area">
                <div className="c-content-header">
                    <div className="c-header-panel">
                        <div className="c-book-chanel"> 创建虚拟机</div>
                        <div className="c-header-desc">
                            <Typography italic>创建可运行的虚拟机，提供虚拟化云相近的云主机体验。</Typography>
                        </div>
                    </div>
                </div>
                <div className={"c-content-body"}>
                    <Collapse accordion={false} hasBorder={false} hasBox={false}>
                        <VmBasicForm/>
                    </Collapse>
                </div>

                <div className={"c-content-body"}>
                    <Collapse accordion={false} hasBorder={false} hasBox={false}>
                        <VmDiscForm/>
                    </Collapse>
                </div>
                <div className={"c-content-body"}>
                    <Collapse accordion={false} hasBorder={false} hasBox={false}>
                        <VmDisksForm dataDiskFields={dataDiskFields} watch={watch}/>
                    </Collapse>
                </div>

                <div className={"c-content-body"}>
                    <Collapse accordion={false} hasBorder={false} hasBox={false}>
                        <VmNetworkForm/>
                    </Collapse>
                </div>

                <div className={"c-content-body"}>
                    <Collapse accordion={false} hasBorder={false} hasBox={false}>
                        <VmAdvancedForm/>
                    </Collapse>
                </div>

                <div className={"c-content-body"}>
                    <Form.Item label="">
                        <Button color="purple" nativeType="submit">
                            保存
                        </Button>
                        <Button onClick={goBack}>
                            返回
                        </Button>
                    </Form.Item>
                </div>

            </div>
        </Form>
    </>
}