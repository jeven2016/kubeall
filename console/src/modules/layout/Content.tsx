import {Button, Dropdown, IconAccount, Layout, Navbar} from "react-windy-ui";
import IconCollapse from "@/modules/icon/CollapseIcon";
import {useCallback} from "react";
import {useTranslation} from "react-i18next";
import {animated, useSpring} from "@react-spring/web";
import {LAYOUT_SETTING} from "@/modules/common/Constants";
import {useAtomValue} from "jotai";
import {selectedModule} from "@/modules/state/global.ts";
import {Outlet, useNavigate} from "react-router";
import {useAtom} from "jotai/index";

export default function Content(props) {
    const currentModule = useAtomValue(selectedModule);
    const {t} = useTranslation();
    const navigate = useNavigate();
    // const { mdWindow } = useContext(WindowContext);
    const mdWindow = false;

    const [module, switchModule] = useAtom(selectedModule);

    const logout = useCallback(() => {
        navigate("/login");
    }, []);

    const {marginLeft} = useSpring({
        marginLeft: mdWindow ? LAYOUT_SETTING.CONTENT.MIN : LAYOUT_SETTING.CONTENT.MAX,
        config: {duration: 200}
    });
    const AnimatedLayout = animated(Layout);


    return (
        <AnimatedLayout extraClassName="base-layout" style={{marginLeft: marginLeft}}>
            <Layout.Split>
                <Layout extraClassName="c-layout-inner" style={{overflowY: "auto"}}>
                    <div className="c-home"></div>
                    <Layout.Content extraClassName="c-layout-content">
                        <Navbar
                            hasBorder={false}
                            hasBox={false}
                            extraClassName="c-navbar-header"
                            hideOnScroll={true}>
                            <Navbar.Title>
                                <Button circle inverted hasBox={false}>
                                    <IconCollapse style={{fontSize: "1.5rem"}}/>
                                </Button>
                            </Navbar.Title>
                            <Navbar.List justify="end">
                                <Navbar.Item active={true}>
                                    <Dropdown
                                        position="bottomRight"
                                        title={
                                            <Button
                                                color="blue"
                                                inverted
                                                hasBox={false}
                                                hasRipple={false}
                                                leftIcon={<IconAccount/>}>
                                                Admin
                                            </Button>
                                        }>
                                        <Dropdown.Menu>
                                            <Dropdown.Item id="item2" onClick={logout}>
                                                {t("global.link.exit")}
                                            </Dropdown.Item>
                                            <Dropdown.Item id="item3">{t("global.link.profile_setting")}</Dropdown.Item>
                                            <Dropdown.Item id="item2">
                                                {module}
                                            </Dropdown.Item>
                                        </Dropdown.Menu>

                                    </Dropdown>
                                </Navbar.Item>
                            </Navbar.List>
                        </Navbar>

                        {/*print the content*/}
                        <Outlet/>

                    </Layout.Content>
                </Layout>
            </Layout.Split>
        </AnimatedLayout>
    );
}
