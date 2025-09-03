import {JSX} from "react";
import {Collapse, Form, Input, Radio, RadioGroup, Select, Space} from "react-windy-ui";
import getOsIcon from "@/modules/icon/OsIcon.tsx";
import {useAtomValue} from "jotai/index";
import {osTypesAtom} from "@/modules/state/global.ts";


export default function VmBasicForm(): JSX.Element {
    const osTypes = useAtomValue(osTypesAtom)
    console.log("data", osTypes)

    return <Collapse.Item header="基本信息" hasBackground value={"basic"}>
        <div style={{paddingTop: "1rem"}}>
            <Form.Item
                label="名称"
                name="name"
                required={true}
                justifyLabel="end"
            >
                <Input block/>
            </Form.Item>

            <Form.Item
                label="显示名称"
                name="displayName"
                justifyLabel="end"
            >
                <Input block/>
            </Form.Item>
            <Form.Item
                label="操作系统"
                name="osType"
                justifyLabel="end">
                <Select>
                    {
                        osTypes.map((osType, index) => {
                            return <Select.Option value={osType} key={`osType-${index}`}>
                                <Space style={{height: '2rem'}}>
                                    <span>{getOsIcon(osType)}</span><span>{osType}</span>
                                </Space>
                            </Select.Option>;
                        })
                    }
                </Select>
            </Form.Item>

            <Form.Item
                label="架构"
                name="arch"
                required={true}
                justifyLabel="end"
            >
                <RadioGroup
                    defaultValue="2"
                    button={true}
                    onChange={(val) => console.log(val)}
                >
                    <Radio value="amd64">X86</Radio>
                    <Radio value="arm64">ARM</Radio>
                </RadioGroup>
            </Form.Item>
            <Form.Item
                label="CPU"
                name="cpu"
                required={true}
                justifyLabel="end"
                labelCol={{col: 2}}
                controlCol={{col: 8}}
            >
                <RadioGroup
                    defaultValue="2"
                    button={true}
                    onChange={(val) => console.log(val)}
                >
                    <Radio value="1">1核</Radio>
                    <Radio value="2">2核</Radio>
                    <Radio value="4">4核</Radio>
                    <Radio value="8">8核</Radio>
                    <Radio value="16">16核</Radio>
                    <Radio value="32">32核</Radio>
                    <Radio value="64">64核</Radio>
                    <Radio value="128">96核</Radio>
                </RadioGroup>
            </Form.Item>

            <Form.Item
                label="内存"
                name="memory"
                justifyLabel="end"
                labelCol={{col: 2}}
                controlCol={{col: 8}}
            >
                <RadioGroup
                    defaultValue="2"
                    button={true}
                    onChange={(val) => console.log(val)}
                >
                    <Radio value="1">1G</Radio>
                    <Radio value="2">2G</Radio>
                    <Radio value="4">4G</Radio>
                    <Radio value="8">8G</Radio>
                    <Radio value="16">16G</Radio>
                    <Radio value="32">32G</Radio>
                    <Radio value="64">64G</Radio>
                    <Radio value="96">96G</Radio>
                    <Radio value="128">128G</Radio>
                    <Radio value="196">196G</Radio>
                    <Radio value="256">256G</Radio>
                </RadioGroup>
            </Form.Item>
            <Form.Item
                label="启动模式"
                name="advancedSettings.bootLoader"
                justifyLabel="end"
            >
                <Select block>
                    <Select.Option value="BIOS">BIOS</Select.Option>
                    <Select.Option value="UEFI">UEFI</Select.Option>
                </Select>
            </Form.Item>
            <Form.Item
                label="描述"
                name="description"
                justifyLabel="end"
            >
                <Input type="textarea" rows={5} block/>
            </Form.Item>


        </div>
    </Collapse.Item>;
}