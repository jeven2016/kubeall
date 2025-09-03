import {Button, Collapse, Form, Select, Table} from "react-windy-ui";

export default function VmNetworkForm() {
    return <Collapse.Item header="网络" hasBackground value={"disk"}>
        <div style={{paddingTop: "1rem"}}>
            <Form.Item
                label="网络"
                justifyLabel="end">
                <Table type="simple" hover hasBorder={false} hasBox={false}>
                    <thead>
                    <tr>
                        <th>接入网络</th>
                        <th>设备模式</th>
                        <th>网卡类型</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr>
                        <td>
                            <Form.Item name="primaryNetwork.networkType" justifyLabel="end" compact>
                                <Select block>
                                    <Select.Option value="pod">容器网络</Select.Option>
                                    <Select.Option value="multus">multus</Select.Option>
                                </Select>
                            </Form.Item>
                        </td>
                        <td>
                            <Form.Item name="primaryNetwork.deviceModel" justifyLabel="end" compact>
                                <Select block>
                                    <Select.Option value="bridge">桥接</Select.Option>
                                    <Select.Option value="masquerade">NAT</Select.Option>
                                </Select>
                            </Form.Item>
                        </td>
                        <td>
                            <Form.Item name="primaryNetwork.interfaceType" justifyLabel="end" compact>
                                <Select block>
                                    <Select.Option value="virtio">virtio</Select.Option>
                                    <Select.Option value="e1000">e1000</Select.Option>
                                </Select>
                            </Form.Item>
                        </td>
                        <td>
                            <Button>+</Button>
                        </td>
                    </tr>
                    </tbody>
                </Table>
            </Form.Item>
        </div>
    </Collapse.Item>;
};