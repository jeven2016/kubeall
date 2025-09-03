import {Button, Col, Divider, Form, Input, Radio, RadioGroup, Row, Select, Space, Typography} from "react-windy-ui";
import getOsIcon from "@/modules/icon/OsIcon.tsx";
import {useRef} from "react";
import {useMutation} from "@tanstack/react-query";
import {constants, resourceType} from "@/modules/common/Constants.tsx";
import {VmImage} from "@/crd";
import {Message} from "@/modules/common/Message.tsx";
import {ImageApi} from "@/modules/api/ImageApi.ts";
import {UploadParam} from "@/types";
import {useLocation, useNavigate} from "react-router";
import {getImageTypeText} from "@/modules/common/VmCommon.tsx";

export default function ImageEdit() {
    const navigate = useNavigate();
    const location = useLocation();
    console.log("location", location);
    const {form, watch, setValue, getValue} = Form.useForm({
        mode: 'onSubmit',
        shouldFocusError: false,
        defaultValues: {
            name: "",
            displayName: "",
            osType: "windows",
            osVersion: "",
            imageType: "iso",
            imageFrom: "upload",
            fileName: "",
        }
    });

    const goBack = () => navigate("/home/vmMgr/image");

    const fileRef = useRef<HTMLInputElement>(null);
    const saveMutation = useMutation({
        mutationFn: (data: VmImage) => {
            data.spec!!.backend = "backingimage"
            return ImageApi.create(constants.imageNamespace, resourceType.vmImage, data);
        },
        onError: (error, variables, context) => {
            Message.error("操作失败", "错误信息：" + error.message, null)
        },
        onSuccess: (data, payload, context) => {
            if (payload.spec!!.imageFrom == "upload") {
                const files = fileRef!!.current!!.files
                if (files && files.length > 0) {
                    let formData = new FormData();
                    formData.append("file", files[0]);

                    const param = {imageName: payload!!.metadata!!.name, formData: formData} as UploadParam;
                    ImageApi.upload(constants.imageNamespace, resourceType.vmImage, param.imageName, param.formData)

                    Message.ok("操作成功", "镜像创建成功，后台正在上传文件，成功前请勿刷新该页面")
                    goBack()
                }
            }
        },
        onSettled: (data, error, variables, context) => {

        },
    })

    const onSubmit = (data, e) => {
        const image = {
            apiVersion: "api.kubeall.io/v1",
            kind: "Image",
            metadata: {
                name: data.name,
                namespace: constants.imageNamespace,
                annotations: {
                    [constants.annotations.displayName]: data.displayName
                }
            },
            spec: {
                osType: data.osType,
                osVersion: data.osVersion,
                imageType: data.imageType,
                imageFrom: data.imageFrom,
                // backend: constants.imageBackend,
            }
        } as VmImage
        console.log("image=", image)
        saveMutation.mutate(image)
    };

    const uploadField = watch('imageFrom', '');
    const uploadFileName = watch('fileName', '');

    return <>
        <Form form={form} onSubmit={onSubmit} direction="horizontal" labelCol={3}>
            <div className="c-content-area">
                <div className="c-content-header">
                    <div className="c-header-panel">
                        <div className="c-book-chanel"> 创建虚拟机镜像</div>
                        <div className="c-header-desc">
                            <Typography italic>创建虚拟机专用镜像</Typography>
                        </div>
                    </div>
                </div>
                <div className={"c-content-body"}>
                    <Row>
                        <Col xl={6} lg={8} md={12}>

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
                                justifyLabel="end"
                            >
                                <Select>
                                    <Select.Option value="windows">
                                        <Space style={{height: '2rem'}}>
                                            <span>{getOsIcon('windows')}</span><span>Windows</span>
                                        </Space>
                                    </Select.Option>

                                    <Select.Option value="ubuntu">
                                        <Space style={{height: '2rem'}}>
                                            <span>{getOsIcon('ubuntu')}</span><span>Ubuntu</span>
                                        </Space>
                                    </Select.Option>

                                    <Divider/>

                                </Select>
                            </Form.Item>
                            <Form.Item
                                label="系统版本"
                                name="osVersion"
                                justifyLabel="end"
                            >
                                <Input block/>
                            </Form.Item>
                            <Form.Item
                                label="镜像类型"
                                name="imageType"
                                justifyLabel="end"
                            >
                                <Select>
                                    <Select.Option
                                        value={constants.imageType.iso}>
                                        {getImageTypeText(constants.imageType.iso)}
                                    </Select.Option>
                                    <Select.Option
                                        value={constants.imageType.disk}>
                                        {getImageTypeText(constants.imageType.disk)}
                                    </Select.Option>
                                </Select>
                            </Form.Item>

                            <Form.Item
                                label="镜像源"
                                name="imageFrom"
                                justifyLabel="end"
                                compact
                            >
                                <RadioGroup>
                                    <Radio value="upload" checkedColor="purple">
                                        本地上传
                                    </Radio>

                                    <Radio value="download" checkedColor="purple">
                                        远程下载
                                    </Radio>
                                </RadioGroup>

                            </Form.Item>
                            <Form.Item
                                label=""
                                justifyLabel="end"
                            >
                                {
                                    uploadField == 'upload' &&
                                    <Input name="localFile" placeholder="请点击选择文件" readOnly
                                           value={uploadFileName}
                                           onClick={() => {
                                               fileRef?.current?.click()
                                           }}/>
                                }
                                {
                                    uploadField == 'upload' &&
                                    <input ref={fileRef} type="file" hidden onChange={(e) => {
                                        const files = e.target.files
                                        if (files && files?.length > 0) {
                                            setValue('fileName', files[0].name, {shouldValidate: true})
                                        }

                                    }}/>
                                }
                            </Form.Item>


                            <Form.Item label="">
                                <Button color="purple" nativeType="submit" disabled={saveMutation.isPending}>
                                    保存
                                </Button>
                                <Button onClick={goBack}>
                                    返回
                                </Button>
                            </Form.Item>

                        </Col>
                    </Row>
                </div>
            </div>
        </Form>
    </>
}