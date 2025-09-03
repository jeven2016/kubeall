import {JSX} from "react";
import {Collapse, Form, Input, Select} from "react-windy-ui";


export default function VmAdvancedForm(): JSX.Element {
    return <Collapse.Item header="高级设置" hasBackground value={"basic"}>
        <div style={{paddingTop: "1rem"}}>
            <Form.Item
                label="CPU模式"
                name="advancedSettings.cpuModel"
                justifyLabel="end"
            >
                <Select block>
                    <Select.Option value="host-model">host-model</Select.Option>
                    <Select.Option value="host-passthrough">CPU直通</Select.Option>
                </Select>
            </Form.Item>

            <Form.Item
                label="虚拟机时区"
                name="advancedSettings.timeZone"
                justifyLabel="end"
            >
                <Input block/>
            </Form.Item>


        </div>
    </Collapse.Item>;
}