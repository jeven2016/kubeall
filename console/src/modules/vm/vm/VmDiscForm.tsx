import {JSX} from "react";
import {Box, Button, Checkbox, Collapse, Form, InputGroup, Select} from "react-windy-ui";
import {SiDiscogs} from "react-icons/si";
import {MdDeleteForever} from "react-icons/md";
import {useAtom} from "jotai";
import {discImagesAtom} from "@/modules/state/VmResources.ts";
import {VmImage} from "@/crd";

export default function VmDiscForm(): JSX.Element {
    const [discs, _] = useAtom<VmImage[]>(discImagesAtom);

    return <Collapse.Item header="光驱" hasBackground value={"disc"}>
        <div style={{paddingTop: "1rem"}}>


            <Form.Item
                name="enableDisc"
                label="启用光驱"
                justifyLabel="end"
            >
                <Checkbox/>
            </Form.Item>
            <Form.Item
                label="光盘镜像"
                justifyLabel="end"
            >
                <InputGroup block>
                    <InputGroup.Item autoScale={false}>
                        <div style={{padding: "0 .5rem"}}><SiDiscogs size="1.5rem"/></div>
                    </InputGroup.Item>

                    <InputGroup.Item autoScale>
                        <Form.Item
                            name="discName"
                            justifyLabel="end">
                            <Select>
                                {
                                    discs.map(() => {

                                    })
                                }
                            </Select>
                        </Form.Item>
                    </InputGroup.Item>
                </InputGroup>
                <Box/>
                <InputGroup block>
                    <InputGroup.Item autoScale={false}>
                        <div style={{padding: "0 .5rem"}}><SiDiscogs size="1.5rem"/></div>
                    </InputGroup.Item>
                    <InputGroup.Item autoScale>
                        <Form.Item
                            name="windowsDriverDisc"
                            justifyLabel="end">
                            <Select>
                                <Select.Option value="bj">Windows驱动盘</Select.Option>
                            </Select>
                        </Form.Item>
                    </InputGroup.Item>
                    <InputGroup.Item autoScale={false}>
                        <Button leftIcon={<MdDeleteForever size="1.5rem"/>}></Button>
                    </InputGroup.Item>
                </InputGroup>
            </Form.Item>


        </div>
    </Collapse.Item>;
}