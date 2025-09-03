import {JSX, useState} from "react";
import {Button, Collapse, Form, Input, Modal, Select, Table} from "react-windy-ui";
import {MdDeleteForever} from "react-icons/md";
import {constants} from "@/modules/common/Constants.tsx";


export default function VmDisksForm({dataDiskFields, watch}): JSX.Element {
    const [active, setActive] = useState(false);

    const close = () => setActive(false);
    const showDisksDialog = () => {
        setActive(true)
    }

    const systemDiskType = watch("systemDisk.type");

    return <Collapse.Item header="磁盘" hasBackground value={"disk"}>
        <div style={{paddingTop: "1rem"}}>
            <Form.Item
                label="系统盘"
                labelCol={{col: 2}}
                controlCol={{col: 8}}
                justifyLabel="end">

                <Table type="striped" hasBorder={false} hasBox={false}>
                    <thead>
                    <tr>
                        <th>磁盘类型</th>
                        <th>磁盘镜像</th>
                        <th>容量</th>
                        {systemDiskType !== constants.hostDiskType && <th>类型</th>}
                        {systemDiskType === constants.pvcDiskType && <th>存储类</th>}
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr>
                        <td>
                            <Form.Item name="systemDisk.type">
                                <Select block>
                                    <Select.Option value="pvc">空盘</Select.Option>
                                    <Select.Option value="image">磁盘镜像</Select.Option>
                                    <Select.Option value="hostDisk">主机磁盘文件</Select.Option>
                                </Select>
                            </Form.Item>
                        </td>
                        <td>
                            <Form.Item name="systemDisk.image">
                                <Input block placeholder="请点击后选择磁盘盘镜像"/>
                            </Form.Item>
                        </td>
                        <td>
                            <Form.Item name="systemDisk.capacity">
                                <Input block icon={"Gi"}/>
                            </Form.Item>
                        </td>
                        <td>
                            {
                                systemDiskType !== constants.hostDiskType &&
                                <Form.Item name="systemDisk.volumeType">
                                    <Select block>
                                        <Select.Option value="Block">块设备</Select.Option>
                                        <Select.Option value="Filesystem">文件系统</Select.Option>
                                    </Select>
                                </Form.Item>
                            }
                        </td>
                        <td>
                            {
                                systemDiskType === constants.pvcDiskType &&
                                <Form.Item name="systemDisk.storageClass">
                                    <Select block>
                                    </Select>
                                </Form.Item>
                            }

                        </td>
                    </tr>
                    </tbody>
                </Table>
            </Form.Item>

            <Form.Item
                label="数据盘"
                labelCol={{col: 2}}
                controlCol={{col: 8}}
                justifyLabel="end">
                <Table type="striped" hasBorder={false} hasBox={false}>
                    <tr>
                        <th>磁盘类型</th>
                        <th>磁盘镜像</th>
                        <th>容量</th>
                        <th>类型</th>
                        <th>存储类</th>
                        <th></th>
                    </tr>
                    {
                        dataDiskFields.map((disk, index) => {
                            return <tr key={`disk-${index}`}>
                                <td>
                                    <Form.Item name={`dataDisks[${index}].type`}>
                                        <Select block>
                                            <Select.Option value="pvc">空盘</Select.Option>
                                            <Select.Option value="image">磁盘镜像</Select.Option>
                                            <Select.Option value="hostDisk">主机磁盘文件</Select.Option>
                                        </Select>
                                    </Form.Item>
                                </td>
                                <td>
                                    <Form.Item name={`dataDisks[${index}].image`}>
                                        <Input block placeholder="请点击后选择磁盘盘镜像"
                                               onClick={() => showDisksDialog()}/>
                                    </Form.Item>
                                </td>
                                <td>
                                    <Form.Item name={`dataDisks[${index}].capacity`}>
                                        <Input block icon={"Gi"}/>
                                    </Form.Item>
                                </td>
                                <td>
                                    <Form.Item name={`dataDisks[${index}].volumeType`}>
                                        <Select block>
                                            <Select.Option value="bj">块设备</Select.Option>
                                            <Select.Option value="nj">文件系统</Select.Option>
                                        </Select>
                                    </Form.Item>
                                </td>
                                <td>
                                    <Form.Item name={`dataDisks[${index}].storageClass`}>
                                        <Select block>
                                            <Select.Option value="bj">sc1</Select.Option>
                                            <Select.Option value="nj">longhorn</Select.Option>
                                        </Select>
                                    </Form.Item>
                                </td>
                                <td><Button leftIcon={<MdDeleteForever size="1.5rem"/>} hasBorder={false}></Button></td>
                            </tr>
                        })
                    }

                </Table>
            </Form.Item>

        </div>
        <Modal
            active={active}
            type="primary"
            onCancel={close}
            size="large"
            autoOverflow={true}
            style={{width: '80%'}}
        >
            <Modal.Header>磁盘镜像</Modal.Header>
            <Modal.Body>
                {[...Array(30).keys()].map((value, index) => (
                    <span key={index}>
              Modal Content....
              <br/>
            </span>
                ))}
            </Modal.Body>
            <Modal.Footer justify="center">
                <Button hasMinWidth={true} type="primary" onClick={close}>
                    OK
                </Button>
                <Button hasMinWidth={true} onClick={close}>
                    Cancel
                </Button>
            </Modal.Footer>
        </Modal>
    </Collapse.Item>;
}