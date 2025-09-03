import {Box, Button, Col, Form, Input, Row, Space, Typography} from "react-windy-ui";
import HomeIcon from "@/modules/icon/HomeIcon";
import {useNavigate} from "react-router";
import {useTranslation} from "react-i18next";

export default function LoginPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    console.log("location=", location)

    const form = Form.useForm({
        mode: "onSubmit",
        defaultValues: {
            username: "",
            password: ""
        }
    });

    const onSubmit = (data: { username: string; password: string }) => {
        navigate("/home/vmMgr");
    };

    return (
        <div className="c-container">
            <div className="c-login-container">
                <div className="c-back"/>
                <div className="c-login-row">
                    <Row>
                        <Col extraClassName="c-login-col" md={8} lg={5} xl={4}>
                            <div className="c-login-info">
                                <Space className="text color-green" block>
                                    <HomeIcon/>
                                    <span>KubeAll</span>
                                </Space>
                                <Box center={<Typography nativeType="h6">容器及虚拟化混合云平台</Typography>}></Box>
                            </div>
                            <Form extraClassName="login-form" form={form} onSubmit={onSubmit}>
                                <Form.Item
                                    name="username"
                                    rules={{
                                        required: "用户名不可为空"
                                    }}
                                    justifyLabel="end"
                                >
                                    <Input placeholder="用户名"/>
                                </Form.Item>
                                <Form.Item
                                    name="password"
                                    rules={{
                                        required: "密码不可为空"
                                    }}
                                    justifyLabel="end"
                                >
                                    <Input type="password" placeholder="密码"/>
                                </Form.Item>
                                <Form.Item direction="horizontal">
                                    <Row>
                                        <Col flexCol justify="start">
                                            <Box block
                                                 hasPadding={false}
                                                 left={<Button block color="blue" hasOutlineBackground outline nativeType="submit"
                                                               invertedOutline hasMinWidth>
                                                     登录
                                                 </Button>}
                                                 right={<Button block color="teal" hasOutlineBackground outline
                                                                invertedOutline hasMinWidth>
                                                     SSO登录
                                                 </Button>}>
                                            </Box>

                                        </Col>
                                    </Row>
                                </Form.Item>
                            </Form>
                        </Col>
                    </Row>
                </div>
            </div>
        </div>
    );
}
